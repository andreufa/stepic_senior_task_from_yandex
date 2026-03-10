package main

import "fmt"

func f(s []int) {
	fmt.Println("in f before", s)
	s = append(s, 10, 11, 12, 13, 14, 15, 16)
	fmt.Println("in f after", s)
}

func main() {
	var a = make([]int, 0, 10)
	a = append(a, 1, 2, 3, 4, 5)
	f(a[1:3])
	fmt.Println("main before append", a)
	a = append(a, 6, 7, 8)
	fmt.Println(a)
}
