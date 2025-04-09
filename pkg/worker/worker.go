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

var (
	pools     = make(map[string]*Pool)
	poolsLock sync.RWMutex
)

func GetPool(key string) *Pool {
	poolsLock.RLock()
	defer poolsLock.RUnlock()
	return pools[key]
}

func InitWorkerPool(key string, workers int) *Pool {
	poolsLock.Lock()
	defer poolsLock.Unlock()

	if existing, ok := pools[key]; ok {
		return existing
	}

	workerPool := &Pool{
		tasks:     make(chan Task, workers),
		workers:   workers,
		waitGroup: sync.WaitGroup{},
	}

	workerPool.startWorkers()
	pools[key] = workerPool

	return workerPool
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
