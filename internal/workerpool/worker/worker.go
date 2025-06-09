package worker

import (
	"log"
	"sync"

	interfaces "github.com/EgorikA4/go-dev-vk-backend-challenge/internal/workerpool"
)

type Worker struct {
	id int
	cb interfaces.TaskCallback
}

func NewWorker(id int, cb interfaces.TaskCallback) interfaces.Worker {
	return &Worker{
		id: id,
		cb: cb,
	}
}

func (w *Worker) GetID() int {
	return w.id
}

func (w *Worker) Start(tasks <-chan string, done <-chan struct{}, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			select {
			case task, ok := <-tasks:
				if !ok {
					log.Printf("Worker %d: tasks closed, exiting\n", w.id)
					return
				}
				log.Printf("Worker %d processing: %s\n", w.id, task)
				w.cb(task)
			case <-done:
				log.Printf("Worker %d: stop received\n", w.id)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {}
