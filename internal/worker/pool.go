package worker

import (
	"concurrent-job-system/internal/db"
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/logger"
	"concurrent-job-system/internal/metrics"
	"context"
	"sync"
)

type Pool struct {
	numWorkers int
	queues     *JobQueueSet
	dispatcher *Dispatcher
	wg         sync.WaitGroup
	logger     logger.ILogger
}

func NewPool(numWorkers int, storage db.IJobRepository, log logger.ILogger) *Pool {
	stats := &JobStats{}
	queues := NewJobQueueSet(100)
	executor := NewJobExecutor(stats, storage, log)
	dispatcher := NewDispatcher(queues, executor, log)

	return &Pool{
		numWorkers: numWorkers,
		queues:     queues,
		dispatcher: dispatcher,
		logger:     log,
	}
}

func (p *Pool) Start(ctx context.Context) {
	p.logger.Info("Starting %d workers", p.numWorkers)
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			p.dispatcher.Run(ctx, workerID)
		}(i)
	}
}

func (p *Pool) Submit(j job.IProcessable) {
	metrics.QueuedJobs.WithLabelValues(j.Type()).Inc()
	metrics.QueuedGauge.Inc()
	id, err := p.dispatcher.Save(j)
	if err != nil {
		p.logger.Error("Failed to save job %d: %v", j.GetId(), err)
		return
	}
	err = j.SetId(id)
	if err != nil {
		p.logger.Error("Failed to set job id %d: %v", j.GetId(), err)
		return
	}
	p.queues.Enqueue(j)
	p.logger.Debug("Enqueued job %d with priority %v", j.GetId(), j.GetPriority())
}

func (p *Pool) Wait() {
	p.wg.Wait()
	p.logger.Info("All workers shut down.")
}

func (p *Pool) Shutdown() {
	p.logger.Info("Shutting down pool...")
	p.Wait()
	p.logger.Info("Pool shutdown complete.")
}

func (p *Pool) Stats() *JobStats {
	return p.dispatcher.stats
}
