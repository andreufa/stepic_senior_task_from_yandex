package main

import (
	"fmt"
	"sync"
)

func fanin(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	go func() {
		for _, ch := range channels {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range ch {
					out <- v
				}
			}()
		}
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	in1 := make(chan int)
	in2 := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range 4 {
			in1 <- i + 300
		}
	}()
	go func() {
		defer wg.Done()
		for i := range 4 {
			in2 <- i + 700
		}
	}()

	go func() {
		wg.Wait()
		close(in1)
		close(in2)
	}()

	out := fanin(in1, in2)

	for v := range out {
		fmt.Println(v)
	}
}
