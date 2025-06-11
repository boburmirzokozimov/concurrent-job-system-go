package worker

import (
	"context"
	"log"
	"sync"
	"time"
)

type Pool struct {
	numWorkers int
	jobQueue   chan Job
	wg         sync.WaitGroup

	Stats *JobStats
}

func NewPool(numWorkers int) *Pool {
	return &Pool{
		numWorkers: numWorkers,
		jobQueue:   make(chan Job, 100),
		Stats:      &JobStats{},
	}
}

func (p *Pool) Start(ctx context.Context) {
	log.Printf("Starting %d workers", p.numWorkers)
	go func() {
		for {
			time.Sleep(1 * time.Second)
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
					log.Printf("Worker %d exiting", workerID)
					return
				case job, ok := <-p.jobQueue:
					if !ok {
						return
					}
					p.HandleJob(job, ctx)
				}
			}
		}(i)
	}
}

func (p *Pool) Submit(job Job, ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case p.jobQueue <- job:
	}
}

func (p *Pool) Wait() {
	close(p.jobQueue)
	p.wg.Wait()
}

func (p *Pool) HandleJob(job Job, ctx context.Context) {
	p.Stats.IncTotal()

	for job.Retries = 0; job.Retries < job.MaxRetryCount; job.Retries++ {
		if err := job.process(); err == nil {
			p.Stats.IncSuccess()
			p.Stats.RecordStatus(job.ID, "success")
			return
		}
		backoff := time.Duration(1<<job.Retries) * time.Second
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			p.Stats.RecordStatus(job.ID, "canceled")
			return
		}
	}

	p.Stats.IncFailed()
	p.Stats.RecordStatus(job.ID, "failed")
}
