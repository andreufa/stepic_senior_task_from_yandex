package main

import (
	"errors"
	"sync"
	"time"
)

// Создать in memory cach с методами добавления, поиска и удаления элемента
// Кеш должен быть конкуретно безопасен
// Каждый элемент должен иметь время жизни

type ICach interface {
	Set(string, string) error
	Get(string) (string, error)
	Del(string) error
}

var ErrorNotFound = errors.New("key not found")

const (
	TTL = 10 * time.Second
)

type Item struct {
	value    string
	exp_time time.Time
}
type Cach struct {
	data    map[string]Item
	mu      sync.RWMutex
	stopped chan struct{}
}

func NewCach() *Cach {
	ticker := time.NewTicker(time.Second)
	c := &Cach{
		data:    make(map[string]Item),
		stopped: make(chan struct{}),
	}

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-c.stopped:
				return
			case <-ticker.C:
				c.clear()
			}
		}
	}()
	return c
}

func (c *Cach) Set(key, value string) error {
	item := Item{value: value, exp_time: time.Now().Add(TTL)}
	c.mu.Lock()
	c.data[key] = item
	c.mu.Unlock()
	return nil
}

func (c *Cach) Get(key string) (string, error) {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		return "", ErrorNotFound
	}
	if item.exp_time.Before(time.Now()) {
		c.Del(key)
		return "", ErrorNotFound
	}
	return item.value, nil
}

func (c *Cach) Del(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *Cach) Stop() {
	c.stopped <- struct{}{}
}

func (c *Cach) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, val := range c.data {
		if val.exp_time.Before(time.Now()) {
			delete(c.data, key)
		}
	}
}
