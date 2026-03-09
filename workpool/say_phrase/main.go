package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func say(id int, phrase string) {
	time.Sleep(2 * time.Second)
	fmt.Printf("worker %d say %s\n", id, phrase)
}

// Создать WorkerPool обрабатывающий слайс фраз

type Worker struct {
	id        int
	handler   func(int, string)
	jobsCount int32
}

func NewWorker(id int, fn func(int, string)) *Worker {
	return &Worker{
		id:      id,
		handler: fn,
	}
}

func (w *Worker) Loop(queue <-chan string) {
	for v := range queue {
		w.handler(w.id, v)
		atomic.AddInt32(&w.jobsCount, 1)
	}
}

type WorkerPool struct {
	queue   chan string
	workers []*Worker
	stopCh  chan struct{}
	once    sync.Once
	wg      sync.WaitGroup
}

func NewWorkerPool(countWorkers, queueSize int, fn func(int, string)) *WorkerPool {
	wp := &WorkerPool{
		queue:   make(chan string, queueSize),
		workers: make([]*Worker, countWorkers),
		stopCh:  make(chan struct{}),
	}
	for i := 0; i < countWorkers; i++ {
		wp.workers[i] = &Worker{id: i, handler: fn}
	}
	return wp

}

func (p *WorkerPool) Submit(val string) {
	select {
	case <-p.stopCh:
		return
	case p.queue <- val:
	}
}

func (p *WorkerPool) Stop() {
	p.once.Do(func() {
		close(p.stopCh)
		close(p.queue)
	})
	p.wg.Wait()
}

func (p *WorkerPool) Start() {
	for _, worker := range p.workers {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			worker.Loop(p.queue)
		}()
	}
}

func (p *WorkerPool) Statistic() {
	for _, w := range p.workers {
		fmt.Printf("worker %d had tasks %d\n", w.id, atomic.LoadInt32(&w.jobsCount))
	}
}

func main() {
	prahses := []string{"one", "two", "tree", "dog", "cat", "airplane", "car", "home", "golang"}
	pool := NewWorkerPool(2, 3, say)
	pool.Start()
	var wg sync.WaitGroup
	for _, p := range prahses {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool.Submit(p)
		}()
	}
	wg.Wait()
	pool.Stop()
	pool.Statistic()
}
