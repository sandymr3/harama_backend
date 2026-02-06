package worker

import (
	"context"
	"log"
	"sync"
)

type Job interface {
	Execute(ctx context.Context) error
	ID() string
}

type WorkerPool struct {
	numWorkers int
	jobQueue   chan Job
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewWorkerPool(numWorkers int, bufferSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		numWorkers: numWorkers,
		jobQueue:   make(chan Job, bufferSize),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	log.Printf("Started %d workers", p.numWorkers)
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()
	log.Printf("Worker %d ready", id)

	for {
		select {
		case <-p.ctx.Done():
			log.Printf("Worker %d shutting down", id)
			return
		case job, ok := <-p.jobQueue:
			if !ok {
				return
			}
			log.Printf("Worker %d starting job: %s", id, job.ID())
			if err := job.Execute(p.ctx); err != nil {
				log.Printf("Worker %d job %s failed: %v", id, job.ID(), err)
			} else {
				log.Printf("Worker %d job %s completed successfully", id, job.ID())
			}
		}
	}
}

func (p *WorkerPool) Submit(job Job) {
	p.jobQueue <- job
}

func (p *WorkerPool) Stop() {
	p.cancel()
	close(p.jobQueue)
	p.wg.Wait()
	log.Println("All workers stopped")
}
