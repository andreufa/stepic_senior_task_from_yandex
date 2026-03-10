package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	tokens chan struct{}
}

func NewSemaphore(maxTokens int) *Semaphore {
	return &Semaphore{
		tokens: make(chan struct{}, maxTokens),
	}
}

func (s *Semaphore) Aquiare() {
	s.tokens <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.tokens
}

const MAX_WORKS = 10

func main() {
	tasks := make([]func(), MAX_WORKS)
	for i := range MAX_WORKS {
		tasks[i] = func() {
			fmt.Printf("working func %d\n", i)
			time.Sleep(time.Second)
		}
	}
	sem := NewSemaphore(3)
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func() {
					sem.Aquiare()
			defer wg.Done()
			task()
			sem.Release()
		}()
	}
	wg.Wait()

}
