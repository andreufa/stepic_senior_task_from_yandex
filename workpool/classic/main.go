package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Worker struct {
	id        int
	handler   func(string)
	jobCounts int32
}

func (w *Worker) Loop(queue <-chan string) {
	for value := range queue {
		w.handler(value)
		atomic.AddInt32(&w.jobCounts, 1)
	}
}

type WorkerPool struct {
	queue   chan string
	workers []*Worker
	stopped chan struct{}
	once    sync.Once
	wg      sync.WaitGroup
}

func NewWorkerPool(maxWorkers, queueSize int, fn func(string)) *WorkerPool {
	wp := &WorkerPool{
		queue:   make(chan string, queueSize),
		workers: make([]*Worker, maxWorkers),
		stopped: make(chan struct{}),
	}

	for i := 0; i < maxWorkers; i++ {
		wp.workers[i] = &Worker{id: i, handler: fn}
	}
	return wp
}

func (p *WorkerPool) Submit(value string) error {
	select {
	case <-p.stopped:
		return fmt.Errorf("pool stopped")
	case p.queue <- value:
		return nil
	}
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

func (p *WorkerPool) Shutdown() {
	p.once.Do(func() {
		close(p.stopped)
		close(p.queue)
	})
	p.wg.Wait()
}

func (p *WorkerPool) Statistic() {
	for _, w := range p.workers {
		fmt.Printf("worker %d had %d tasks", w.id, atomic.LoadInt32(&w.jobCounts))
	}
}

func main() {
	say := func(phrase string) {
		time.Sleep(time.Second)
		fmt.Printf("say %s", phrase)

	}
	pool := NewWorkerPool(5, 2, say)
	go pool.Start()

	for i := 0; i < 11; i++ {
		phrase := fmt.Sprintf("phrase %d\n", i)
		err := pool.Submit(phrase)
		if err != nil {
			fmt.Println(err)
		}
	}

	var input string
	fmt.Scan(&input)
	pool.Shutdown()
}
