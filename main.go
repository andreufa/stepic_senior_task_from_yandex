package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Request interface{}
type Response interface{}

type Backend interface {
	Invoke(ctx context.Context, req Request) (Response, error)
}

type BackendImpl struct {
	addr      string
	available bool
	jobs      int32
}

func (b *BackendImpl) Invoke(ctx context.Context, req Request) (Response, error) {
	return nil, nil
}

type Balancer struct {
	services []*BackendImpl
	mu       sync.Mutex
	wg       sync.WaitGroup
	stopped  chan struct{}
}

func NewBalancer(addrs []string) *Balancer {
	b := &Balancer{
		services: make([]*BackendImpl, len(addrs)),
		stopped:  make(chan struct{}),
	}

	for i, adr := range addrs {
		b.services[i] = &BackendImpl{
			addr:      adr,
			available: true,
		}
	}

	go func() {
		for {
			select {
			case <-b.stopped:
				return
			case <-time.After(2 * time.Second):
				b.checkHealth()
			}
		}

	}()

	return b
}
func (b *Balancer) Destroy() {
	close(b.stopped)
}

func (b *Balancer) checkHealth() {
	_services := make([]*BackendImpl, len(b.services))
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	b.mu.Lock()
	copy(_services, b.services)
	b.mu.Unlock()
	for i := 0; i < len(b.services); i++ {
		if _services[i].available {
			continue
		}
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			_, err := _services[i].Invoke(ctx, nil)
			if err == nil {
				b.mu.Lock()
				b.services[i].available = true
				b.mu.Unlock()
			}
		}()
	}
	b.wg.Wait()
}

func (b *Balancer) getIdFreeService() (int, error) {
	minJobs := int32(math.MaxInt32)
	serviceId := -1
	for i, s := range b.services {
		b.mu.Lock()
		available := s.available
		b.mu.Unlock()
		if !available {
			continue
		}

		if jobs := atomic.LoadInt32(&s.jobs); jobs < minJobs {
			minJobs = jobs
			serviceId = i
		}
	}
	if serviceId == -1 {
		return serviceId, fmt.Errorf("not available service")
	}
	return serviceId, nil
}

func (b *Balancer) Invoke(ctx context.Context, req Request) (Response, error) {
	serviceId, err := b.getIdFreeService()
	if err != nil {
		return nil, err
	}
	atomic.AddInt32(&b.services[serviceId].jobs, 1)
	resp, err := b.services[serviceId].Invoke(ctx, req)
	if err != nil {
		b.mu.Lock()
		b.services[serviceId].available = false
		b.mu.Unlock()
		resp, err = b.RetryInvoke(ctx, req)
		return resp, err
	}
	return resp, nil

}

func (b *Balancer) RetryInvoke(ctx context.Context, req Request) (Response, error) {
	serviceId, err := b.getIdFreeService()
	if err != nil {
		return nil, err
	}
	resp, err := b.services[serviceId].Invoke(ctx, req)
	if err != nil {
		b.mu.Lock()
		b.services[serviceId].available = false
		b.mu.Unlock()
		return nil, err
	}
	return resp, nil

}
