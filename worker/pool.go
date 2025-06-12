package worker

import (
	"context"
	"sync"
	"time"
)

type Pool struct {
	numWorkers  int
	lowQueue    chan Processable
	normalQueue chan Processable
	highQueue   chan Processable
	wg          sync.WaitGroup
	Stats       *JobStats
	Storage     *FileJobStorage
}

func NewPool(numWorkers int, s *FileJobStorage) *Pool {
	return &Pool{
		numWorkers:  numWorkers,
		highQueue:   make(chan Processable, 100),
		normalQueue: make(chan Processable, 100),
		lowQueue:    make(chan Processable, 100),
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
					var job Processable
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

func (p *Pool) Submit(job Processable, ctx context.Context) {
	p.Storage.Save(job)
	select {
	case <-ctx.Done():
		return
	default:
		switch job.GetPriority() {
		case Low:
			p.lowQueue <- job
		case Normal:
			p.normalQueue <- job
		case High:
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

func (p *Pool) HandleJob(job Processable, ctx context.Context) {
	p.Stats.IncTotal()

	for job.GetRetries() < job.GetMaxRetryCount() {
		job.IncRetry()
		LogJobStart(job)

		if err := job.Process(); err == nil {
			p.Stats.IncSuccess()
			LogJobSuccess(job)
			p.Stats.RecordStatus(job.GetId(), "success")
			p.Storage.MarkCompleted(job.GetId(), "success")
			return
		}

		backoff := time.Duration(1<<job.GetRetries()) * time.Second
		LogJobRetry(job, backoff)

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			LogJobCanceled(job)
			p.Stats.RecordStatus(job.GetId(), "canceled")
			p.Storage.MarkCompleted(job.GetId(), "canceled")
			return
		}
	}

	p.Stats.IncFailed()
	LogJobFail(job)
	p.Stats.RecordStatus(job.GetId(), "failed")
	p.Storage.MarkCompleted(job.GetId(), "failed")
}
