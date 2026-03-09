package main

import (
	"fmt"
	"sync"
	"time"
)

func say(id int, phrase string) {
	time.Sleep(20 * time.Millisecond)
	fmt.Printf("worker %d say %s\n", id, phrase)
}

// Создать WorkerPool обрабатывающий слайс фраз
// func makePool(poolSize int, handler func(int, string))(handle func(string), wait func())

type Worker struct {
	id     int
	handle func(int, string)
}

func makePool(poolSize int, handler func(int, string)) (handle func(string), wait func()) {
	workers := make(chan Worker, poolSize)
	var wg sync.WaitGroup

	for i := 0; i < poolSize; i++ {
		workers <- Worker{id: i, handle: handler}
	}
	handle = func(val string) {
		worker := <-workers
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker.handle(worker.id, val)
			workers <- worker
		}()
	}

	wait = func() {
		wg.Wait()
	}

	return handle, wait

}

func main() {
	phrases := []string{}

	for i := 0; i < 100; i++ {
		phrases = append(phrases, fmt.Sprintf("phrase %d", i))
	}

	handle, wait := makePool(5, say)
	for _, p := range phrases {
		handle(p)
	}
	wait()

}
