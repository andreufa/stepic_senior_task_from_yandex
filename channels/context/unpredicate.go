package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func unpredicateFunc() int {
	n := rand.Intn(7)
	time.Sleep(time.Duration(n) * time.Second)
	return n
}

func predicateFunc1() int {
	done := make(chan struct{})
	var res int
	go func() {
		res = unpredicateFunc()
		close(done)
	}()
	select {
	case <-done:
		fmt.Println("res", res)
		return res
	case <-time.After(5 * time.Second):
		fmt.Println("trunked by timeout")
	}
	return 0
}

func predicateFunc2(ctx context.Context) (int, error) {
	var res int
	done := make(chan struct{})

	var cancel context.CancelFunc

	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second) // Переопределили контекст, если нет Deadline
		defer cancel()
	}

	go func() {
		res = unpredicateFunc()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return res, ctx.Err()
	case <-done:
		return res, nil
	}
}

func main() {
	for range 20 {
		ctx := context.Background()
		r, err := predicateFunc2(ctx)
		fmt.Println(r, err)
	}
}
