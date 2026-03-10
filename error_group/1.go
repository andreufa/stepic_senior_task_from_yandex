package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

func fetchId(ctx context.Context, id int) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		// Имитируем задержку для демонстрации работы контекста
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		value := rand.Intn(100)
		if value > 50 {
			return value, nil
		}
		return value, fmt.Errorf("value %d < 50", value)
	}
}

func multyFetch(ctx context.Context, ids []int) (int, error) {
	group, ectx := errgroup.WithContext(ctx)
	group.SetLimit(5) // Ограничиваем количество одновременных горутин

	var (
		sum int
		mu  sync.Mutex // Мьютекс для защиты общей переменной sum
	)

	for _, id := range ids {
		id := id // Захватываем значение id в локальной переменной
		group.Go(func() error {
			res, err := fetchId(ectx, id)
			if err != nil {
				return fmt.Errorf("fetchId(%d) failed: %w", id, err)
			}

			// Защищаем доступ к общей переменной sum
			mu.Lock()
			defer mu.Unlock()
			sum += res
			return nil
		})
	}

	// Ждём завершения всех горутин и получаем ошибку, если она есть
	err := group.Wait()
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func main() {
	ids := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Создаём контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	sum, err := multyFetch(ctx, ids)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Total sum: %d\n", sum)
}
