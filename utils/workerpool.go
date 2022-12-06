package utils

import "sync"

type T = interface{}

type IWorkerPool interface {
	Run()
	AddTask(task func())
	Wait()
}

type WorkerPool struct {
	maxWorker   int
	queuedTaskC chan func()
	wg          *sync.WaitGroup
}

func NewWorkerPool(maxWorker int) *WorkerPool {
	return &WorkerPool{
		maxWorker:   maxWorker,
		queuedTaskC: make(chan func()),
		wg:          &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.maxWorker; i++ {
		wp.wg.Add(1)
		go func(workerID int) {
			for task := range wp.queuedTaskC {
				task()
			}
			wp.wg.Done()
		}(i + 1)
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func (wp *WorkerPool) AddTask(task func()) {
	wp.queuedTaskC <- task
}
