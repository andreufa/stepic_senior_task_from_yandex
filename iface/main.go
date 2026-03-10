package main

import (
	"fmt"
)

type ABC interface {
	A()
	B()
	C()
}
type AB interface {
	A()
	B()
}

type abc struct{}

func (a abc) A() {}
func (a abc) B() {}
func (a abc) C() {}

type ab struct{}

func (a ab) A() {}
func (a ab) B() {}

func main() {
	var a interface{}
	a = abc{}

	a = a.(AB)

	switch a.(type) {
	case ab:
		println("ab")
	default:
		fmt.Printf("%T\n", a)
	}

	a1, ok := a.(ABC)
	fmt.Println(ok)
	fmt.Printf("%T\n", a1)

}
