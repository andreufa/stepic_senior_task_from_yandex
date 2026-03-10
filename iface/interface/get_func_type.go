package main

import (
	"fmt"
	"reflect"
	"strings"
)

//Задание: Реализуй функцию process так, чтобы она:

// Для int: выводила число в квадрате

// Для string: выводила строку задом наперед

// Для []int: выводила сумму элементов

// Для map[string]int: выводила количество ключей

// Для nil: выводила "got nil"

// Для функций: выводила тип функции и её адрес

// Для всех остальных типов: выводила "unknown type: %T"
func process(v interface{}) {
	// Твой код здесь
	switch val := v.(type) {
	case int:
		fmt.Println(val * val)
	case []int:
		sum := 0
		for _, vl := range val {
			sum += vl
		}
		fmt.Println(sum)
	case string:
		runes := []rune(val)
		var builder strings.Builder
		for i := len(runes) - 1; i >= 0; i-- {
			builder.WriteRune(runes[i])
		}
		fmt.Println(builder.String())
	case map[string]int:
		fmt.Println(len(val))
	case nil:
		fmt.Println("got nil")
	default:
		vl := reflect.ValueOf(val)
		if vl.Kind() == reflect.Func {
			fmt.Printf("function with %d args, returns %d values\n",
				vl.Type().NumIn(), vl.Type().NumOut())
		} else {
			fmt.Printf("unknown type: %T\n", v)
		}

	}

}

func main() {
	process(42)
	process("hello")
	process([]int{1, 2, 3})
	process(map[string]int{"a": 1})
	process(nil)
	process(process) // да, передаем саму функцию
}
