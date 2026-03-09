package main

import "sync"

// Задача выполнить все  задачи строго определенным количеством горутин

const MAX_GOROUTINES = 5

func Procedure (tasks <-chan func()){
var wg sync.WaitGroup

	for range MAX_GOROUTINES{
		wg.Add(1)
		go func(){
			defer wg.Done()
			for v := range tasks{
				v()
			}
		}()
	}
	wg.Wait()
}
