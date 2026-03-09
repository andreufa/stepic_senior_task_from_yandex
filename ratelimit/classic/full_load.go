package main

import (
	"fmt"
	"sync"
	"time"
)

func Request(id string) {
	fmt.Printf("request %s\n", id)
}

type Limiter struct {
	ticker *time.Ticker
	tokens chan struct{}
	once   sync.Once
	stopCh chan struct{}
}

func NewLimiter(rpc int) *Limiter {
	limiter := &Limiter{
		ticker: time.NewTicker(time.Second),
		tokens: make(chan struct{}, rpc),
		stopCh: make(chan struct{}),
	}
	for i := 0; i < rpc; i++ {
		limiter.tokens <- struct{}{}
	}
	go func() {
		defer limiter.ticker.Stop()
	MainLoop:
		for {
			select {
			case <-limiter.ticker.C:
				for i := 0; i < rpc; i++ {
					select {
					case limiter.tokens <- struct{}{}:
					default:
						i = rpc
					}
				}
			case <-limiter.stopCh:
				break MainLoop
			}
		}
	}()

	return limiter
}

func (l *Limiter) Stop() {
	l.once.Do(func() {
		close(l.stopCh)
		close(l.tokens)
	})
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
	limiter := NewLimiter(10)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if limiter.Allow() {
				Request(fmt.Sprintf("id %d", id))
			}
		}(i)
	}
	wg.Wait()
	limiter.Stop()
}
