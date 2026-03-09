package main

import (
	"fmt"
	"sync"
	"time"
)

// На выполнение всех задач запустить ограниченое количество горутин работающих одновременно
const (
	TASK_COUNT  = 30
	MAX_WORKERS = 4
)

func Proccessing(maxWorkers int, tasks []func()) {
	workers := make(chan struct{}, maxWorkers)
	wg := sync.WaitGroup{}

	for _, task := range tasks {
		wg.Add(1)
		workers <- struct{}{}
		go func() {
			defer wg.Done()
			task()
			<-workers

		}()
	}
	wg.Wait()
}



func main() {
	tasks := make([]func(), TASK_COUNT)
	for i := 0; i < TASK_COUNT; i++ {
		tasks[i] = func() {
			time.Sleep(1 * time.Second)
			fmt.Println("working task", i)
		}
	}
	Proccessing(MAX_WORKERS, tasks)

}
