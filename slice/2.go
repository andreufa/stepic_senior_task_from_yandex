package main

import "fmt"

func change(s []int) {
	s = append(s, 4)
	s[0] = 6
}

func main() {
	sl := make([]int, 2)
	for i := range 2 {
		sl[i] = i + 1
	}
	sl = append(sl, 98)
	fmt.Println(sl)
	change(sl)
	fmt.Println("-------")
	fmt.Println(sl)
	mp := make(map[int]string)
	mp[1] = "Bob"
	changeMap(mp)
	println(mp[1])

}

func changeMap(m map[int]string) {
	m[1] = "Alice"
}
