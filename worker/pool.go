package worker

import (
	"concurrent-job-system/internal"
	"concurrent-job-system/models"
	"context"
	"sync"
	"time"
)

type Pool struct {
	numWorkers  int
	lowQueue    chan models.Processable
	normalQueue chan models.Processable
	highQueue   chan models.Processable
	wg          sync.WaitGroup
	Stats       *JobStats
	Storage     config.IStorage
}

func NewPool(numWorkers int, s config.IStorage) *Pool {
	return &Pool{
		numWorkers:  numWorkers,
		highQueue:   make(chan models.Processable, 100),
		normalQueue: make(chan models.Processable, 100),
		lowQueue:    make(chan models.Processable, 100),
		Stats:       &JobStats{},
		Storage:     s,
	}
}

func (p *Pool) Start(ctx context.Context) {
	jobLog.info.Printf("Starting %d workers", p.numWorkers)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			p.Stats.Print()
		}
	}()

	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					jobLog.info.Printf("Worker %d exiting", workerID)
					return
				default:
					var job models.Processable
					var ok bool

					select {
					case job, ok = <-p.highQueue:
					default:
						select {
						case job, ok = <-p.normalQueue:
						default:
							select {
							case job, ok = <-p.lowQueue:
							default:
								time.Sleep(100 * time.Millisecond) // Nothing to do
								continue
							}
						}
					}

					if !ok || job == nil {
						continue
					}
					p.HandleJob(job, ctx)
				}
			}

		}(i)
	}
}

func (p *Pool) Submit(job models.Processable, ctx context.Context) {
	err := p.Storage.Save(job)
	if err != nil {
		return
	}
	select {
	case <-ctx.Done():
		return
	default:
		switch job.GetPriority() {
		case models.Low:
			p.lowQueue <- job
		case models.Normal:
			p.normalQueue <- job
		case models.High:
			p.highQueue <- job
		}
	}
}

func (p *Pool) Wait() {
	close(p.lowQueue)
	close(p.normalQueue)
	close(p.highQueue)
	p.wg.Wait()
	jobLog.info.Println("All workers shut down gracefully.")
}

func (p *Pool) HandleJob(job models.Processable, ctx context.Context) {
	p.Stats.IncTotal()

	for job.GetRetries() < job.GetMaxRetryCount() {
		job.IncRetry()
		LogJobStart(job)

		if err := job.Process(); err == nil {
			p.Stats.IncSuccess()
			LogJobSuccess(job)
			p.Stats.RecordStatus(job.GetId(), "success")
			p.Storage.UpdateStatus(job.GetId(), "success")
			return
		}

		backoff := time.Duration(1<<job.GetRetries()) * time.Second
		LogJobRetry(job, backoff)

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			LogJobCanceled(job)
			p.Stats.RecordStatus(job.GetId(), "canceled")
			p.Storage.UpdateStatus(job.GetId(), "canceled")
			return
		}
	}

	p.Stats.IncFailed()
	LogJobFail(job)
	p.Stats.RecordStatus(job.GetId(), "failed")
	p.Storage.UpdateStatus(job.GetId(), "failed")
}
