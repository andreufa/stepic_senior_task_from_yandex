package main

import (
	"context"
	"sync"
	"time"
)

type Request interface{}
type Response interface{}

type Backend interface {
	Invoke(ctx context.Context, req Request) (Response, error)
}
type BackendImpl struct {
	addr      string
	jobs      int32
	available bool
}

func (b *BackendImpl) Invoke(ctx context.Context, req Request) (Response, error) {
	return nil, nil
}
func NewBackend(addr string) *BackendImpl {
	return &BackendImpl{
		addr:      addr,
		available: true,
	}
}

type Balancer struct {
	services []*BackendImpl
	mu       sync.Mutex
	stopped  chan struct{}
}

func NewBalancer(addrs []string) *Balancer {
	b := &Balancer{
		services: make([]*BackendImpl, len(addrs)),
		stopped:  make(chan struct{}),
	}

	for id, adr := range addrs {
		b.services[id] = &BackendImpl{addr: adr}
	}
	go b.healthCheckLoop()

	return b
}

func (b *Balancer) healthCheckLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
Loop:
	for {
		select {
		case <-b.stopped:
			break Loop
		case <-ticker.C:
			b.checkHealth()
		}
	}
}

func (b *Balancer) checkHealth() {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_servises := make([]*BackendImpl, len(b.services))
	b.mu.Lock()
	copy(_servises, b.services)
	b.mu.Unlock()

	var wg sync.WaitGroup

	for i, s := range _servises {
		if s.available {
			continue
		}
		wg.Add(1)
		go func(idx int, svc *BackendImpl) {
			defer wg.Done()
			_, err := svc.Invoke(ctx, nil)

			if err == nil {
				indx := idx
				b.mu.Lock()
				b.services[indx].available = true
				b.mu.Unlock()
			}
		}(i, s)
	}

	wg.Wait()
}

func (b *Balancer) Invoke(ctx context.Context, req Request) (Response, error) {
	return nil, nil
}
