package main

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	mu         sync.RWMutex
	capacity   int
	expiration time.Duration
	cacheMap   map[string]*list.Element
	ll         *list.List
}

type entry struct {
	key       string
	value     string
	timestamp time.Time
}

func NewCache(capacity int, expiration time.Duration) *Cache {
	return &Cache{
		capacity:   capacity,
		expiration: expiration,
		cacheMap:   make(map[string]*list.Element),
		ll:         list.New(),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, found := c.cacheMap[key]; found {
		kv := elem.Value.(*entry)
		if time.Since(kv.timestamp) > c.expiration {
			c.ll.Remove(elem)
			delete(c.cacheMap, key)
			return "", false
		}
		c.ll.MoveToFront(elem)
		return kv.value, true
	}
	return "", false
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cacheMap[key]; found {
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*entry)
		kv.value = value
		kv.timestamp = time.Now()
		return
	}

	if c.ll.Len() == c.capacity {
		oldest := c.ll.Back()
		if oldest != nil {
			c.ll.Remove(oldest)
			kv := oldest.Value.(*entry)
			delete(c.cacheMap, kv.key)
		}
	}

	elem := c.ll.PushFront(&entry{key: key, value: value, timestamp: time.Now()})
	c.cacheMap[key] = elem
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cacheMap[key]; found {
		c.ll.Remove(elem)
		delete(c.cacheMap, key)
	}
}
