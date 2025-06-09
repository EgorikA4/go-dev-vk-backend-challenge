package pool

import (
	"fmt"
	"log"
	"sync"

	interfaces "github.com/EgorikA4/go-dev-vk-backend-challenge/internal/workerpool"
	"github.com/EgorikA4/go-dev-vk-backend-challenge/internal/workerpool/worker"
)

type WorkerPool struct {
	tasks   chan string
	workers map[int]interfaces.Worker
	doneMap map[int]chan struct{}
	mu      sync.Mutex
	nextID  int
	wg      sync.WaitGroup
}

func NewWorkerPool(buffer int) interfaces.Pool {
	return &WorkerPool{
		tasks:   make(chan string, buffer),
		workers: make(map[int]interfaces.Worker),
		doneMap: make(map[int]chan struct{}),
		nextID:  1,
	}
}

func (wp *WorkerPool) AddWorker(cb interfaces.TaskCallback) int {
	wp.mu.Lock()
	id := wp.nextID
	wp.nextID++
	done := make(chan struct{})

	worker := worker.NewWorker(id, cb)
	wp.workers[id] = worker
	wp.doneMap[id] = done
	wp.mu.Unlock()

	wp.wg.Add(1)
	worker.Start(wp.tasks, done, &wp.wg)

	log.Printf("Started worker %d\n", id)
	return id
}

func (wp *WorkerPool) RemoveWorkerByID(id int) error {
	wp.mu.Lock()
	done, ok := wp.doneMap[id]
	if !ok {
		wp.mu.Unlock()
		return fmt.Errorf("worker with id: %d not found", id)
	}
	delete(wp.workers, id)
	delete(wp.doneMap, id)
	wp.mu.Unlock()

	close(done)
	log.Printf("Signaled worker %d to stop\n", id)
	return nil
}

func (wp *WorkerPool) Submit(task string) {
	wp.tasks <- task
}

func (wp *WorkerPool) GetWorkersID() []int {
	wp.mu.Lock()
	workersID := make([]int, 0, len(wp.workers))
	for workerID := range wp.workers {
		workersID = append(workersID, workerID)
	}
	wp.mu.Unlock()
	return workersID
}

func (wp *WorkerPool) Shutdown() {
	close(wp.tasks)
	wp.mu.Lock()
	for id, done := range wp.doneMap {
		close(done)
		log.Printf("Signaled worker %d to stop\n", id)
	}
	wp.mu.Unlock()
	wp.wg.Wait()
	log.Println("Pool shut down")
}
