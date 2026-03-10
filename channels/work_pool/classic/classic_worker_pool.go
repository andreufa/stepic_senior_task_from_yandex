package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id        int
	jobsCount int32
	handler   func(int)
}

func NewWorker(id int, fn func(int)) *Worker {
	return &Worker{id: id, handler: fn}
}

type WorkerPool struct {
	queune  chan int
	workers []*Worker
	stopCh  chan struct{}
	wg      sync.WaitGroup
	once    sync.Once
}

func NewWorkerPool(maxWorkers, queuneSize int, fn func(int)) *WorkerPool {
	wp := &WorkerPool{
		queune:  make(chan int, queuneSize),
		workers: make([]*Worker, maxWorkers),
		stopCh:  make(chan struct{}),
	}
	for i := range maxWorkers {
		wp.workers[i] = &Worker{id: i, handler: fn}
	}
	return wp
}

func (p *WorkerPool) Submit(task int) {
	select {
	case <-p.stopCh:
		return
	case p.queune <- task:
	}
}

func (p *WorkerPool) Start() {
	for _, w := range p.workers {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-p.stopCh:
					return
				case val, ok := <-p.queune:
					if !ok {
						return
					}
					w.handler(val)
				}
			}

		}()
	}
}

func (p *WorkerPool) Stop() {
	p.once.Do(func() {
		close(p.stopCh)
	})
	p.wg.Wait()
}

func (p *WorkerPool) Statistic() {
	for _, w := range p.workers {
		fmt.Printf("worker %d had % d jobs\n", w.id, w.jobsCount)
	}
}

func main() {
	tasks := make([]int, 10)
	for i := range 10 {
		tasks[i] = i
	}

	fn := func(task int) {
		time.Sleep(time.Second)
		fmt.Printf("task %d proccesing\n", task)
	}

	pool := NewWorkerPool(3, 6, fn)
	pool.Start()

	for t := range tasks {
		pool.Submit(t)
	}
	time.Sleep(10 * time.Second)
	pool.Stop()
	pool.Statistic()
}
