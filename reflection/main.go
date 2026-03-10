package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string
	Age  int
}

func main() {
	anna := User{
		Name: "Anna",
		Age:  35,
	}

	val := reflect.ValueOf(anna)
	if val.Kind() == reflect.Struct {
		fmt.Println(val.NumField())
		for i := 0 ; i <val.NumField(); i++{
			fmt.Println(val.Field(i))
		}
		for i := 0 ; i <val.NumField(); i++{
			fmt.Println(val.Type().Field(i))
		}
	}
}
