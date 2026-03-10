package main

import (
	"fmt"
	"math/rand"
)

// Требуется реализовать функцию, которая генирирует слайс требуемой длины и с уникальными значениями

func getUniqSlice(n int) []int {
	m := make(map[int]struct{})
	s := make([]int, 0, n)
	for len(s) < n {
		value := rand.Intn(100)
		if _, ok := m[value]; ok {
			continue
		}
		s = append(s, value)
		m[value] = struct{}{}
	}
	return s
}

func main() {
	s := getUniqSlice(5)
	fmt.Println(s)
}
