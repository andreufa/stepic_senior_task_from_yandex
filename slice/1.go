package main

import (
	"fmt"
)

func main() {
	var arr []*int
	number := 0

	for i := 0; i < 3; i++ {
		arr = append(arr, &number)
		number++
	}

	fmt.Println(arr) // 1) что выведет программа?

	// 2) окей, что выведется в таком случае?
	for _, num := range arr {
		fmt.Println(*num)
	}

	trippleSlice(arr)
	fmt.Println("after triple")
	// 4) что выведется потом?
	for _, num := range arr {
		fmt.Println(*num)
	}
}

// 3) задание: реализовать функцию, которая увеличит каждый элемент слайса в 3 раза
func trippleSlice(slice []*int) {
	length := len(slice)
	old := &slice
	fmt.Println("old addr:",old)
	slice = append(slice, &length)
	fmt.Println("slice addr:",slice)
	mp := make(map[*int]struct{}, len(slice))

	for _, num := range slice {
		_, ok := mp[num]
		if ok {
			continue
		}
		mp[num] = struct{}{}
		*num = *num * 3

	}
}
