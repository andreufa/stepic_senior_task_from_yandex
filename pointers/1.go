package main

import "fmt"

type Person struct {
	name string
	age  int
}

func (p *Person) ChangeName(name string) {
	p = &Person{
		name: name,
	}
}

func main() {
	p := Person{
		name: "Bob",
		age:  25,
	}

	p.ChangeName("Alice")
	fmt.Println(p)
}