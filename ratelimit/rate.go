package main

import (
	"sync"
	"time"
)

// Задача устанавливать не более MAX_CONNECT соединений в секунду. Лишние соединения сбрасываем
const MAX_CONNECT = 1000 // rps

type Connect struct {
	ip string
}

func CreateConnect(ip string) *Connect {
	return &Connect{ip: ip}
}

func ConnectionWithLimit(conn []string) {
	ticker := time.NewTicker(time.Second)
	tokens := make(chan struct{}, MAX_CONNECT)
	for i := 0; i < MAX_CONNECT; i++ {
		tokens <- struct{}{}
	}
	go func() {
		for range ticker.C { // Каждую секунду
			for i := 0; i < MAX_CONNECT; i++ {
				select {
				case tokens <- struct{}{}:
				default:
					break
				}
			}
		}
	}()

	var wg sync.WaitGroup

	for _, cn := range conn {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-tokens:
				_ = CreateConnect(cn)
			default: // Сбрасываем соединение
			}

		}()
	}
	wg.Wait()
}
