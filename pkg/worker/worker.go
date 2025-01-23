package worker

import (
	"sync"
)

type Task func()

type Pool struct {
	tasks     chan Task
	workers   int
	waitGroup sync.WaitGroup
}

var instance *Pool

func GetPool() *Pool {
	return instance
}

func InitWorkerPool(workers int) {
	workerPool := &Pool{
		tasks:     make(chan Task, workers),
		workers:   workers,
		waitGroup: sync.WaitGroup{},
	}

	workerPool.startWorkers()

	instance = workerPool
}

func (wp *Pool) AddTask(task Task) {
	wp.tasks <- task
}

func (wp *Pool) startWorkers() {
	for i := 0; i < wp.workers; i++ {
		wp.waitGroup.Add(1)
		go func(workerID int) {
			defer wp.waitGroup.Done()
			for task := range wp.tasks {
				task()
			}
		}(i)
	}
}
