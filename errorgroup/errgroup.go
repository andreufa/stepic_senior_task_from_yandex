package main

import (
	"context"
	"fmt"
	"math/rand"

	"golang.org/x/sync/errgroup"
)

// Используйте errgroup для ограничения количества горутин

func fetch(ctx context.Context, val int) error {
	select {
	case <-ctx.Done(): // Если контекст отменен, то сразу выходим
		return ctx.Err()
	default:
	}
	if val < 50 {
		return fmt.Errorf("Value < 50 -> %d", val)
	}
	return nil
}

func main() {

	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(5)

	for i := 0; i < 100; i++ {
		rValue := rand.Intn(100)
		group.Go(func() error {
			return fetch(ctx, rValue)

		})
	}
	err := group.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
