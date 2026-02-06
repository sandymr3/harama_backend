package unit_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"harama/internal/worker"

	"github.com/stretchr/testify/assert"
)

type MockJob struct {
	id        string
	shouldFail bool
	executed  bool
	mu        sync.Mutex
	wg        *sync.WaitGroup
}

func (m *MockJob) ID() string {
	return m.id
}

func (m *MockJob) Execute(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.executed = true
	m.wg.Done()
	
	if m.shouldFail {
		return errors.New("job failed")
	}
	return nil
}

func TestWorkerPool(t *testing.T) {
	// Setup
	var wg sync.WaitGroup
	pool := worker.NewWorkerPool(3, 10) // 3 workers, buffer 10
	pool.Start()
	defer pool.Stop()

	// 1. Submit successful jobs
	job1 := &MockJob{id: "job-1", wg: &wg}
	job2 := &MockJob{id: "job-2", wg: &wg}
	
	wg.Add(2)
	pool.Submit(job1)
	pool.Submit(job2)

	// Wait for execution
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for jobs")
	}

	assert.True(t, job1.executed, "Job 1 should be executed")
	assert.True(t, job2.executed, "Job 2 should be executed")

	// 2. Submit failing job (should not crash worker)
	wg.Add(1)
	failJob := &MockJob{id: "fail-job", shouldFail: true, wg: &wg}
	pool.Submit(failJob)
	
	// Wait for it to "finish"
	wg.Wait() // Re-using waitgroup logic locally would require reset, but here strictly adding works if sequential or careful
	assert.True(t, failJob.executed)
}
