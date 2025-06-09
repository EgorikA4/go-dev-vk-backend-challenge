package worker_test

import (
	"sync"
	"testing"
	"time"

	"github.com/EgorikA4/go-dev-vk-backend-challenge/internal/workerpool/worker"
	"github.com/stretchr/testify/assert"
)

func TestWorkerSuccess(t *testing.T) {
	output := make(chan string, 1)
	cb := func(task string) {
		output <- task
	}

	w := worker.NewWorker(1, cb)
	assert.Equal(t, 1, w.GetID())

	tasks := make(chan string, 1)
	done := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)
	w.Start(tasks, done, &wg)

	task := "Test Task"
	tasks <- task

	select {
	case got := <-output:
		assert.Equal(t, task, got)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("worker did not process task in time")
	}

	close(tasks)
	close(done)
	wg.Wait()
}

func TestWorkerStopsOnDone(t *testing.T) {
    var mu sync.Mutex
    output := make([]string, 0)
    cb := func(task string) {
        mu.Lock()
        output = append(output, task)
        mu.Unlock()
    }

    w := worker.NewWorker(1, cb)
    assert.Equal(t, 1, w.GetID())

    tasks := make(chan string, 1)
    done := make(chan struct{})

    var wg sync.WaitGroup
    wg.Add(1)
    w.Start(tasks, done, &wg)

    close(done)

    doneCh := make(chan struct{})
    go func() {
        wg.Wait()
        close(doneCh)
    }()

    select {
    case <-doneCh:
    case <-time.After(100 * time.Millisecond):
        t.Fatal("worker did not stop after done was closed")
    }

    tasks <- "should-not-be-processed"
    time.Sleep(20 * time.Millisecond)

    mu.Lock()
    defer mu.Unlock()
    assert.Empty(t, output, "callback should not be called after done signal")
}