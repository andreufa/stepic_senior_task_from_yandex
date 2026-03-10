package main

import "fmt"

func changeNo() int { // Вернет 0
	var p int
	defer func() {
		p = 100
	}()
	return p
}

func change() (p int) { // Вернет 100
	defer func() {
		p = 100
	}()
	return
}

func main() {
	p := change()
	fmt.Println(p)
}
