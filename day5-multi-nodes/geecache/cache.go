package geecache

import (
	"github.com/ATao13/gee-chache/day5-multi-nodes/geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

func (c *cache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}
