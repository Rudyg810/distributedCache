package cache

import "sync"

type SafeCache struct {
	mu         sync.Mutex
	lru        *Cache
	cacheBytes int64
}

func NewSafeCache(cacheBytes int64) *SafeCache {
	return &SafeCache{
		cacheBytes: cacheBytes,
	}
}

func (c *SafeCache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *SafeCache) Get(key string) (value ByteView, ok bool) {
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
