package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id int
}

func NewWorker(id int) *Worker {
	return &Worker{id: id}
}

type WorkerPool struct {
	queue   chan int
	workers chan *Worker
	handler func(int)
	stopCh  chan struct{}
	wg      sync.WaitGroup
	once    sync.Once
}

func NewWorkerPool(workerCount, queueSize int, fn func(int)) *WorkerPool {
	wp := &WorkerPool{
		queue:   make(chan int, queueSize),
		workers: make(chan *Worker, workerCount),
		handler: fn,
		stopCh:  make(chan struct{}),
	}
	for i := 0; i < workerCount; i++ {
		worker := NewWorker(i)
		wp.workers <- worker
	}
	return wp
}

func (w *WorkerPool) Submit(value int) {
	select {
	case w.queue <- value:
		// Задача добавлена в очередь
	case <-w.stopCh:
		// Пул остановлен, задача не будет обработана
		return
	}
}

func (w *WorkerPool) Start() {
	go func() {
		for {
			select {
			case <-w.stopCh:
				close(w.queue) // Закрываем очередь для graceful shutdown
				return
			case v, ok := <-w.queue:
				if !ok {
					// Очередь закрыта, выходим
					return
				}
				w.wg.Add(1) // Увеличиваем счётчик только при получении задачи
				worker := <-w.workers
				go func(worker *Worker, v int) {
					defer func() {
						w.workers <- worker // Возвращаем рабочего в пул
						w.wg.Done()         // Уменьшаем счётчик задач
					}()
					fmt.Printf("worker %d processing task %d\n", worker.id, v)
					w.handler(v)
				}(worker, v)
			}
		}
	}()
}

func (w *WorkerPool) Wait() {
	w.wg.Wait()
}

func (w *WorkerPool) Stop() {
	w.once.Do(func() {
		close(w.stopCh) // Закрываем канал остановки
	})
	w.Wait() // Ждём завершения всех задач
}

func main() {
	fn := func(v int) {
		time.Sleep(1 * time.Second)
		fmt.Printf("task %d worked\n", v)
	}

	pool := NewWorkerPool(3, 10, fn)
	pool.Start()

	for i := range 20 {
		pool.Submit(i)
	}

	time.Sleep(3 * time.Second)
	pool.Submit(32)
	time.Sleep(10 * time.Second)
	pool.Stop()
}
