package main

import (
	"context"
	"errors"
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
	inflight  int32 // используем atomic для безопасного подсчета
}

func (b *BackendImpl) Invoke(_ context.Context, req Request) (Response, error) {
	// Реальная реализация здесь
	return nil, nil
}

func NewBackend(addr string) *BackendImpl {
	return &BackendImpl{
		addr:      addr,
		available: true,
	}
}

// Balancer реализует балансировку с наименьшей нагрузкой
type Balancer struct {
	services []*BackendImpl
	mu       sync.RWMutex
	stopped  chan struct{}
}

func NewBalancer(addrs []string) *Balancer {
	b := &Balancer{
		services: make([]*BackendImpl, len(addrs)),
		stopped:  make(chan struct{}),
	}

	for i, addr := range addrs {
		b.services[i] = NewBackend(addr)
	}

	// Периодическая проверка доступности
	go b.healthCheckLoop()

	return b
}

func (b *Balancer) healthCheckLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.checkHealth()
		case <-b.stopped:
			return
		}
	}
}

func (b *Balancer) checkHealth() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	b.mu.RLock()
	services := make([]*BackendImpl, len(b.services))
	copy(services, b.services)
	b.mu.RUnlock()

	var wg sync.WaitGroup
	for i, s := range services {
		if s.available {
			continue
		}

		wg.Add(1)
		go func(idx int, svc *BackendImpl) {
			defer wg.Done()

			// Пробуем вызвать с таймаутом
			_, err := svc.Invoke(ctx, nil)

			b.mu.Lock()
			if err == nil {
				b.services[idx].available = true
			}
			b.mu.Unlock()
		}(i, s)
	}

	// Ждем завершения всех проверок, но не дольше контекста
	wg.Wait()
}

// getLeastLoaded находит наименее загруженный доступный сервис
func (b *Balancer) getLeastLoaded() (*BackendImpl, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var selected *BackendImpl
	minLoad := int32(^uint32(0) >> 1) // MaxInt32

	for _, s := range b.services {
		if !s.available {
			continue
		}

		load := atomic.LoadInt32(&s.inflight)
		if load < minLoad {
			minLoad = load
			selected = s
		}
	}

	if selected == nil {
		return nil, errors.New("no available backends")
	}

	return selected, nil
}

func (b *Balancer) Invoke(ctx context.Context, req Request) (Response, error) {
	// Получаем наименее загруженный сервис
	backend, err := b.getLeastLoaded()
	if err != nil {
		return nil, err
	}

	// Увеличиваем счетчик активных запросов
	atomic.AddInt32(&backend.inflight, 1)
	defer atomic.AddInt32(&backend.inflight, -1)

	// Выполняем запрос
	resp, err := backend.Invoke(ctx, req)

	// Если ошибка - помечаем сервис как недоступный
	if err != nil {
		b.mu.Lock()
		backend.available = false
		b.mu.Unlock()

		// Пробуем ретрай с другим сервисом
		return b.retryInvoke(ctx, req)
	}

	return resp, nil
}

func (b *Balancer) retryInvoke(ctx context.Context, req Request) (Response, error) {
	// Проверяем, не истек ли контекст
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Пробуем найти другой доступный сервис
	backend, err := b.getLeastLoaded()
	if err != nil {
		return nil, errors.New("all backends failed")
	}

	atomic.AddInt32(&backend.inflight, 1)
	defer atomic.AddInt32(&backend.inflight, -1)

	resp, err := backend.Invoke(ctx, req)
	if err != nil {
		b.mu.Lock()
		backend.available = false
		b.mu.Unlock()
		return nil, err // Не ретраим повторно, чтобы избежать каскада
	}

	return resp, nil
}

func (b *Balancer) Stop() {
	close(b.stopped)
}
