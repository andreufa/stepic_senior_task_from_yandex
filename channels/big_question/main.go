package main

import (
	"fmt"
	"sync"
)

// Задача 1.1: Базовый WorkerPool
// Концепции: горутины, каналы, sync.WaitGroup

// Создайте структуру WorkerPool, которая:

// Принимает количество воркеров при создании

// Имеет метод Submit(func()) для добавления задач

// Имеет метод Wait() который ждет завершения всех задач

// Воркеры должны быть горутинами, которые запускаются при создании пула

type WorkerPool struct {
	tasks chan func()    // Канал задач (очередь)
	wg    sync.WaitGroup // Для ожидания завершения
	queneSize int
}

func NewWorkerPool(workerCount, queneSize int) *WorkerPool {
	wp := &WorkerPool{
		queneSize: queneSize,
		tasks: make(chan func(), queneSize), // Буферизированный канал
	}

	// Запускаем воркеров
	for i := 0; i < workerCount; i++ {
		go wp.worker(i)
	}

	return wp
}

// worker - горутина, которая выполняет задачи
func (wp *WorkerPool) worker(id int) {
	for task := range wp.tasks {
		fmt.Printf("Воркер %d выполняет задачу\n", id)
		task()
		wp.wg.Done() // Уменьшаем счетчик после выполнения
	}
}

func (wp *WorkerPool) Submit(task func()) {
	wp.wg.Add(1) // Увеличиваем счетчик ДО отправки
	wp.tasks <- task
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()    // Ждем выполнения всех задач
	close(wp.tasks) // Закрываем канал (воркеры завершатся)
}

func main() {
	pool := NewWorkerPool(3)

	for i := 0; i < 10; i++ {
		i := i
		pool.Submit(func() {
			fmt.Printf("Task %d\n", i)
		})
	}

	pool.Wait()
}
