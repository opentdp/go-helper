package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	createdAt time.Time
	value     any
}

type TimedCache struct {
	mutex    sync.RWMutex
	interval time.Duration
	stopChan chan struct{}
	items    map[string]CacheItem
}

func NewTimedCache(interval time.Duration) *TimedCache {
	cache := &TimedCache{
		interval: interval,
		stopChan: make(chan struct{}),
		items:    make(map[string]CacheItem),
	}

	go func() {
		ticker := time.NewTicker(cache.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cache.DeleteExpired()
			case <-cache.stopChan:
				return
			}
		}
	}()

	return cache
}

func (c *TimedCache) Set(key string, value any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = CacheItem{
		createdAt: time.Now(),
		value:     value,
	}
}

func (c *TimedCache) Get(key string) (any, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if item, ok := c.items[key]; ok {
		if time.Since(item.createdAt) < c.interval {
			return item.value, true
		}
		delete(c.items, key)
	}
	return nil, false
}

func (c *TimedCache) DeleteExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, item := range c.items {
		if time.Since(item.createdAt) > c.interval {
			delete(c.items, key)
		}
	}
}

func (c *TimedCache) StopCleanup() {
	close(c.stopChan)
}
