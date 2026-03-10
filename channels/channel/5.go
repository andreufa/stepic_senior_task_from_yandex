package main

import "fmt"

func main() {
	ch := make(chan int)
	close(ch)
	select {
	case <-ch: // zero value
		fmt.Println("zero value")
	default:
		fmt.Println("default")
	}
	// select {
	// case ch <- 4: // panic
	// default:
	// 	fmt.Println("default")
	// }
}
