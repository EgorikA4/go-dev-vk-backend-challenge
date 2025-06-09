package pool_test

import (
	"testing"
	"time"

	"github.com/EgorikA4/go-dev-vk-backend-challenge/internal/workerpool/pool"
	"github.com/stretchr/testify/assert"
)

func TestPoolWithOneWorker(t *testing.T) {
	const bufferSize = 4
	wp := pool.NewWorkerPool(bufferSize)
	assert.Empty(t, wp.GetWorkersID(), "expected no workers at start")

	output := make(chan string, 1)
	cb := func(task string) {
		output <- task
	}

	id := wp.AddWorker(cb)
	assert.Equal(t, 1, id, "first worker should have ID=1")
	assert.Len(t, wp.GetWorkersID(), 1, "should have one worker after AddWorker")

	task1 := "Test Task"
	wp.Submit(task1)
	select {
	case got := <-output:
		assert.Equal(t, task1, got)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("worker did not process first task in time")
	}

	err := wp.RemoveWorkerByID(2)
	assert.Error(t, err, "removing unknown worker should return error")

	err = wp.RemoveWorkerByID(id)
	assert.NoError(t, err, "removing existing worker should succeed")
	assert.Empty(t, wp.GetWorkersID(), "no workers should remain after removal")

	task2 := "Test Task2"
	wp.Submit(task2)
	select {
	case got := <-output:
		t.Fatalf("unexpectedly processed task after removal: %s", got)
	case <-time.After(100 * time.Millisecond):
	}

	wp.Shutdown()
}

func TestShutdownSignalsAllWorkers(t *testing.T) {
	const bufferSize = 1

	wp := pool.NewWorkerPool(bufferSize)
	wp.AddWorker(func(string) {})
	wp.AddWorker(func(string) {})

	done := make(chan struct{})
	go func() {
		wp.Shutdown()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Shutdown did not complete in time; done channels may not have been closed")
	}
}

func TestSubmitPanicsAfterShutdown(t *testing.T) {
	const bufferSize = 1

	wp := pool.NewWorkerPool(bufferSize)
	wp.AddWorker(func(string) {})

	wp.Shutdown()
	assert.Panics(t, func() {
		wp.Submit("should panic")
	}, "expected Submit to panic when called after Shutdown")
}