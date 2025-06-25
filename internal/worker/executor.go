package worker

import (
	"concurrent-job-system/internal/db"
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/logger"
	"concurrent-job-system/internal/metrics"
	"context"
	"time"
)

type JobExecutor struct {
	stats   *JobStats
	storage db.IJobRepository
	logger  logger.ILogger
}

func NewJobExecutor(stats *JobStats, storage db.IJobRepository, log logger.ILogger) *JobExecutor {
	return &JobExecutor{
		stats:   stats,
		storage: storage,
		logger:  log,
	}
}
func (e *JobExecutor) Save(j job.IProcessable) (int, error) {
	return e.storage.Save(j)
}

func (e *JobExecutor) Execute(j job.IProcessable, ctx context.Context) {
	start := time.Now()
	metrics.RunningJobs.Inc()
	defer func() {
		metrics.JobLatency.Observe(time.Since(start).Seconds())
		metrics.RunningJobs.Dec()
	}()

	e.stats.IncTotal()
	e.logger.Info("Executing job %d", j.GetId())

	var lastErr error
	for j.GetRetries() < j.GetMaxRetryCount() {
		j.IncRetry()
		e.logger.Debug("Job %d retry #%d", j.GetId(), j.GetRetries())
		lastErr = j.Process()
		if lastErr == nil {
			e.markJobAs(j, "success")
			metrics.SucceededJobs.Inc()
			return
		} else {
			e.logger.Warn("Job %d failed attempt #%d", j.GetId(), j.GetRetries())
		}

		backoff := time.Duration(1 << j.GetRetries() * int(time.Second))
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			e.markJobAs(j, "canceled")
			metrics.DeadLetterJobs.Inc()
			return
		}
	}

	metrics.FailedJobs.Inc()
	e.markJobAsFailed(j, lastErr)
}

func (e *JobExecutor) markJobAs(j job.IProcessable, status string) {
	e.logger.Info("Job %d marked as %s", j.GetId(), status)
	switch status {
	case "success":
		e.stats.IncSuccess()
	case "failed":
		e.stats.IncFailed()
	}
	e.stats.RecordStatus(j.GetId(), status)
	err := e.storage.UpdateStatus(j.GetId(), status)
	if err != nil {
		e.logger.Warn("Job %d failed to mark as %s - %s", j.GetId(), status, err)
		return
	}
}

func (e *JobExecutor) markJobAsFailed(j job.IProcessable, lastErr error) {
	status := "failed"
	e.logger.Info("Job %d marked as %s with error: %v", j.GetId(), status, lastErr)
	e.stats.IncFailed()
	e.stats.RecordStatus(j.GetId(), status)

	errMessage := "unknown error"
	if lastErr != nil {
		errMessage = lastErr.Error()
	}

	err := e.storage.UpdateFailedStatus(j.GetId(), status, errMessage)
	if err != nil {
		e.logger.Warn("Job %d failed to mark as %s - %s", j.GetId(), status, err)
		return
	}
}

func (e *JobExecutor) LoadPending() job.IProcessable {
	j, err := e.storage.LoadPending()
	if err != nil {
		e.logger.Error("Failed to load pending job")
	}
	return j
}
