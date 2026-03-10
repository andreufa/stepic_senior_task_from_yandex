package main

import "fmt"

// Что произойдет ?

func spawnMessage(n int) chan string {
	ch := make(chan string, 1)
	go func() {
		for i := 0; i < n; i++ {
			ch <- fmt.Sprintf("msg %d", i+1)
		}
	}()

	return ch
}

func main() {
	n := 10
	for msg := range spawnMessage(n) {
		fmt.Println("recived: ", msg)
	}
}
