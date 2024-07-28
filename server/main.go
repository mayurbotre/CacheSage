// main.go
package main

import (
	"container/list"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Cache is a simple LRU cache
type Cache struct {
	mu         sync.RWMutex
	capacity   int
	expiration time.Duration
	cacheMap   map[string]*list.Element
	ll         *list.List
}

// Entry holds a key-value pair and its timestamp
type Entry struct {
	key       string
	value     string
	timestamp time.Time
}

// NewCache initializes a new Cache
func NewCache(capacity int, expiration time.Duration) *Cache {
	return &Cache{
		capacity:   capacity,
		expiration: expiration,
		cacheMap:   make(map[string]*list.Element),
		ll:         list.New(),
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, found := c.cacheMap[key]; found {
		kv := elem.Value.(*Entry)
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

// Set adds or updates a key-value pair in the cache
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cacheMap[key]; found {
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*Entry)
		kv.value = value
		kv.timestamp = time.Now()
		return
	}

	if c.ll.Len() == c.capacity {
		oldest := c.ll.Back()
		if oldest != nil {
			c.ll.Remove(oldest)
			kv := oldest.Value.(*Entry)
			delete(c.cacheMap, kv.key)
		}
	}

	elem := c.ll.PushFront(&Entry{key: key, value: value, timestamp: time.Now()})
	c.cacheMap[key] = elem
}

// Delete removes a key-value pair from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cacheMap[key]; found {
		c.ll.Remove(elem)
		delete(c.cacheMap, key)
	}
}

func main() {
	cache := NewCache(10, 5*time.Second)

	r := gin.Default()

	// Endpoint to get a value from the cache
	r.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, found := cache.Get(key)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"message": "Key not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	// Endpoint to add or update a key-value pair
	r.POST("/cache", func(c *gin.Context) {
		var req struct {
			Key       string        `json:"key"`
			Value     string        `json:"value"`
			ExpiresIn time.Duration `json:"expires_in"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}
		cache.Set(req.Key, req.Value)
		c.JSON(http.StatusOK, gin.H{"message": "Key added/updated"})
	})

	// Endpoint to delete a key-value pair
	r.DELETE("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		cache.Delete(key)
		c.JSON(http.StatusOK, gin.H{"message": "Key deleted"})
	})

	r.Run(":8080") // Start the server on port 8080
}
