package main

import (
	"hash/fnv"
	"sync"
)

// Реализуйте InMemoryCashe
type Cache interface {
	Set(k, v string)
	Get(k string) (string, bool)
}

type Shard struct {
	data map[string]string
	mu   sync.RWMutex
}
type InMemoryCache struct {
	shards []Shard
}

func (c *Shard) Set(k, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[k] = v
}

func (c *Shard) Get(k string) (string, bool) {
	c.mu.RLock()
	value, ok := c.data[k]
	c.mu.RUnlock()
	return value, ok
}

func NewInMemoryCache(count int) *InMemoryCache {
	shards := make([]Shard, count)
	for i := range count {
		shards[i] = Shard{
			data: make(map[string]string),
			mu:   sync.RWMutex{},
		}
	}

	return &InMemoryCache{
		shards: shards,
	}
}

func (c *InMemoryCache) Set(k, v string) {
	h := hasher(k)
	shardId := h % len(c.shards)
	c.shards[shardId].Set(k, v)
}

func (c *InMemoryCache) Get(k string) (string, bool) {
	h := hasher(k)
	shardId := h % len(c.shards)
	value, ok := c.shards[shardId].data[k]
	return value, ok
}

func hasher(s string) int {
	//TO DO
	h := fnv.New32()
	h.Write([]byte(s))
	return int(h.Sum32())
}
