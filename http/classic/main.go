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
// Предусмотреть возможность отмены запроса по таймаутуbefore
const (
	TIME_OUT = 2 * time.Second
)

var client = &http.Client{
	Timeout: TIME_OUT,
}

func procces(urls []string) map[int]int {
	statusCodes := map[int]int{}
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, url := range urls {
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
	return statusCodes
}
