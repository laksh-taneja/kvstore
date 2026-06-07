package cache

import (
	"sync"

	"github.com/laksh-taneja/kvstore/internal/lru"
)

type Cache struct {
	Store *lru.Cache
	mu    sync.Mutex // can't use rmutex get also do a write on linkedlist
}

func New(capacity int) (*Cache, error) {
	core, err := lru.LRUCache(capacity)
	if err != nil {
		return nil, err
	}
	return &Cache{
		Store: core,
	}, nil
}

func (c *Cache) Access(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Store.Get(key)
}

func (c *Cache) Write(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Store.Put(key, value)
}
