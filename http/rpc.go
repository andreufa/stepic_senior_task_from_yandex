package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"https://google.com",
		"https://yandex.ru",
		"https://amazon.com",
		"https://youtube.com",
	}
	mp := procces(urls)
	fmt.Println(mp)
}

//Реализовать паралельные запросы по адресам
// Подсчитать количество для каждого StatusCode ответа
// Предусмотреть возможность отмены запроса по таймауту

const (
	TIME_OUT = 2 * time.Second
)

var client = &http.Client{
	Timeout: TIME_OUT,
}

type Limiter struct {
	ticker *time.Ticker
	tokens chan struct{}
	once   sync.Once
	stopCh chan struct{}
}

func NewLimiter(rps int) *Limiter {
	limiter := &Limiter{
		ticker: time.NewTicker(time.Second / time.Duration(rps)),
		tokens: make(chan struct{}, rps),
		stopCh: make(chan struct{}),
	}
	for i := 0; i < rps; i++ {
		limiter.tokens <- struct{}{}
	}
	go func() {
		defer limiter.ticker.Stop()
	Loop:
		for {
			select {
			case <-limiter.ticker.C:
				limiter.tokens <- struct{}{}
			case <-limiter.stopCh:
				break Loop
			}
		}
	}()
	return limiter
}
func (l *Limiter) Allow() bool {
	select {
	case <-l.tokens:
		return true
	default:
		return false
	}
}

func (l *Limiter) Stop() {
	l.once.Do(func() {
		close(l.stopCh)
	})
}

func procces(urls []string) map[int]int {
	statusCodes := map[int]int{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	limiter := NewLimiter(10)

	for _, url := range urls {
		if !limiter.Allow() {
			continue
			// Тут может быть разное поведение. continue - сбрасываем текущий реквест но продолжаем идти по циклу
			// Можно <- limiter.tokens - просто блокировать функцию до поступления новых лимитов
			// break - прекратить обработку при достижении лимита
		}
		wg.Add(1)
		go func(ursl string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				fmt.Printf("wrong url | %s\n", err)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("error by request %s\n", err)
				return
			}
			defer resp.Body.Close() // Не забываем закрывать тело ответа
			mu.Lock()
			statusCodes[resp.StatusCode]++
			mu.Unlock()
		}(url)
	}
	wg.Wait()
	limiter.Stop()
	return statusCodes
}
