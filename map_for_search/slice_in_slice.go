package main

import "fmt"

func GetCommonElements(slice1, slice2 []int) []int {
	result := []int{}

	_slice1 := make(map[int]bool, len(slice1))

	for _, elem := range slice1 {
		_slice1[elem] = true
	}

	for _, elem := range slice2 {
		if _slice1[elem] {
			result = append(result, elem)
		}
	}
	return result
}

func main() {
	s1 := []int{1, 3, 5, 6, 7, 8, 911, 2, 23, 0}
	s2 := []int{3, 5, 23, 556, 9, 10, 0}

	result := GetCommonElements(s1, s2)
	fmt.Println(result)
}
