package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Нужно ограничить количество запросов в секунду

const (
	RPS = 10
)

var client http.Client

func Request(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return -1, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

type Limiter struct {
	ticker *time.Ticker
	tokens chan struct{}
	wg     sync.WaitGroup
}

func NewLimiter(ctx context.Context, rps int) *Limiter {
	l := &Limiter{
		ticker: time.NewTicker(time.Second / time.Duration(rps)),
		tokens: make(chan struct{}, rps),
	}
	for i := 0; i < rps; i++ {
		l.tokens <- struct{}{}
	}
	go func() {
		defer l.ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-l.ticker.C:
				l.tokens <- struct{}{}

			}
		}
	}()

	return l
}

func (l *Limiter) Allow() bool {
	select {
	case <-l.tokens:
		return true
	default:
		return false
	}
}

func main() {
	urls := []string{
		"https://google.com",
		"https://yandex.ru",
		"https://amazon.com",
		"https://youtube.com",
	}
	ctx := context.Background()
	limiter := NewLimiter(ctx, RPS)
	for _, url := range urls {
		if limiter.Allow() {
			limiter.wg.Add(1)
			go func() {
				defer limiter.wg.Done()
				code, err := Request(ctx, url)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(code)
			}()
		}
	}
	limiter.wg.Wait()

}
