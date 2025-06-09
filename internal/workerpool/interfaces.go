package interfaces

import "sync"

type TaskCallback func(task string)

type Worker interface {
	Start(tasks <-chan string, done <-chan struct{}, wg *sync.WaitGroup)
	GetID() int
	Stop()
}

type Pool interface {
	AddWorker(cb TaskCallback) int
	RemoveWorkerByID(id int) error
	GetWorkersID() []int
	Submit(string)
	Shutdown()
}
