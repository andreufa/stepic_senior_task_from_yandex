package main

import "fmt"

type ABC interface {
	A()
	B()
	C()
}

type BC interface {
	B()
	C()
}

type abc struct{}

func (a abc) A() {}
func (a abc) B() {}
func (a abc) C() {}

func main() {
	var a interface{}
	var a1 interface{}
	a = abc{}

	a1 = a.(BC)
	fmt.Printf("%T\n", a1) // main.abc (но тип переменной interface{})
	a1.(BC).B()            // a1.B() так не работает

	a1 = a1.(ABC)
	fmt.Printf("%T\n", a1) // main.abc (но тип переменной interface{})
	a1.(ABC).A()           //a1.A() не работает
}
