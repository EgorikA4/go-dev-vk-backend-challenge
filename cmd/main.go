package main

import (
	"fmt"
	"time"

	"github.com/EgorikA4/go-dev-vk-backend-challenge/internal/workerpool/pool"
)

const bufferSize = 10

func main() {
	wp := pool.NewWorkerPool(bufferSize)

	cb := func(task string) {
		fmt.Println("completed task: ", task)
	}
	for range 4 {
		wp.AddWorker(cb)
	}

	for taskID := range 50 {
		wp.Submit(fmt.Sprintf("task%d", taskID+1))
	}

	wp.RemoveWorkerByID(2)
	wp.RemoveWorkerByID(4)

	for taskID := range 50 {
		wp.Submit(fmt.Sprintf("task%d", taskID+51))
	}

	time.Sleep(time.Second)
	wp.Shutdown()
}
