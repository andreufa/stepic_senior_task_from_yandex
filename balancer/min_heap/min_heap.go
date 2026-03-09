package main

import (
	"container/heap"
	"fmt"
)

// Пациент с номерком
type Patient struct {
	number int
	name   string
}

// MinHeap реализует интерфейс heap.Interface
type MinHeap []Patient

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].number < h[j].number }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Patient))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func main() {
	// Создаём кучу
	h := &MinHeap{}
	heap.Init(h)

	// Добавляем пациентов
	heap.Push(h, Patient{5, "Анна"})
	heap.Push(h, Patient{2, "Борис"})
	heap.Push(h, Patient{7, "Виктор"})
	heap.Push(h, Patient{1, "Галина"})

	fmt.Println("Очередь к врачу:")
	for h.Len() > 0 {
		patient := heap.Pop(h).(Patient)
		fmt.Printf("Приглашается %s с номерком %d\n", patient.name, patient.number)
	}
}
