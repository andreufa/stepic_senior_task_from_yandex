package main

import (
	"fmt"
	"sync"
	"time"
)

func workers(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(1 * time.Millisecond)
}

// func main() { // Сколько времени займет?
// 	runtime.GOMAXPROCS(20)
// 	wg := &sync.WaitGroup{}
// 	const MAX_TASKS = 10_000
// 	wg.Add(MAX_TASKS)
// 	start := time.Now()
// 	for range MAX_TASKS {
// 		go workers(wg)
// 	}
// 	wg.Wait()
// 	fmt.Println(time.Since(start))
// }

func main() { // Что не так?
	counter := 0
	for i := 0; i < 100; i++ {
		go func() {
			counter++
		}()
	}
	fmt.Println(counter)
}
