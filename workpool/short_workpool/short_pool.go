package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Worker struct {
	id        int
	handle    func(string)
	jobsCount int32
}

type WorkerPool struct {
	maxWorkers int
	workers    chan *Worker
	queue      chan string
	stoped     chan struct{}
	wg         sync.WaitGroup
	once       sync.Once
}

func NewWorkerPool(queueSize, maxWorkers int, handler func(string)) *WorkerPool {
	wp := &WorkerPool{
		maxWorkers: maxWorkers,
		workers:    make(chan *Worker, maxWorkers),
		queue:      make(chan string, queueSize),
		stoped:     make(chan struct{}),
	}
	for i := 0; i < maxWorkers; i++ {
		wp.workers <- &Worker{id: i, handle: handler}
	}
	return wp
}

func (p *WorkerPool) Submit(value string) error {
	select {
	case <-p.stoped:
		return fmt.Errorf("workerpool stoped")
	default:
		p.queue <- value
		return nil
	}
}

func (p *WorkerPool) Start() {
	for {
		select {
		case <-p.stoped:
			return
		case value := <-p.queue:
			worker := <-p.workers
			p.wg.Add(1)
			go func() {
				defer func() {
					p.wg.Done()
					select {
					case <-p.stoped:
						// Не возвращаем при остановке
					default:
						p.workers <- worker
					}
				}()
				worker.handle(value)
				atomic.AddInt32(&worker.jobsCount, 1)
			}()
		}
	}
}

func (p *WorkerPool) Shutdown() {
	p.once.Do(func() {
		close(p.stoped)
	})
	p.wg.Wait()
	close(p.queue)
	for i := 0; i < p.maxWorkers; i++ {
		w := <-p.workers
		fmt.Printf("worker %d had %d tasks\n", w.id, atomic.LoadInt32(&w.jobsCount))
	}
	close(p.workers)
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
		pool.Submit(phrase)
	}

	var input string
	fmt.Scan(&input)
	pool.Shutdown()
	pool.Submit("more")
}
