package worker

import (
	"github.com/labstack/echo/v4"
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

func RunWithWorker[T any](handler func(c echo.Context) error, poolName string) echo.HandlerFunc {
	return func(c echo.Context) error {
		var body T
		if err := c.Bind(&body); err != nil {
			return c.JSON(400, map[string]string{"error": "Invalid request format"})
		}

		if err := c.Validate(body); err != nil {
			return c.JSON(422, map[string]string{"error": "Validation failed"})
		}

		done := make(chan error)
		wp := GetPool(poolName)

		wp.AddTask(func() {
			done <- handler(c)
		})

		return <-done
	}
}
