package main

import "fmt"

func modify(start, end int) []byte {
	arr := [999999]byte{}

	slice := make([]byte, end-start)

	copy(slice, arr[start:end])
	return slice

}

func main() {
	

	s := modify(10, 20)
	fmt.Println(cap(s))

}
