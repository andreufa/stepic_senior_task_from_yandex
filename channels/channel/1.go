package main

import (
	"fmt"
	"time"
)

// Сколько времени займет?

func worker() <-chan int {
	ch := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		close(ch)
	}()
	return ch
}

func main() {
	start := time.Now()
	// _, _ = worker(), worker()
	_, _ = <-worker(), <-worker()
	fmt.Println(time.Since(start))

}
