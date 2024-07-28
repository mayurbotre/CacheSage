package main

import (
	"container/list"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Cache struct {
	mu       sync.RWMutex
	capacity int
	cacheMap map[string]*list.Element
	ll       *list.List
}

type Entry struct {
	Key        string        `json:"key"`
	Value      string        `json:"value"`
	Timestamp  time.Time     `json:"timestamp"`
	Expiration time.Duration `json:"expiration"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		cacheMap: make(map[string]*list.Element),
		ll:       list.New(),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, found := c.cacheMap[key]; found {
		kv := elem.Value.(*Entry)
		if time.Since(kv.Timestamp) > kv.Expiration {
			c.ll.Remove(elem)
			delete(c.cacheMap, key)
			return "", false
		}
		c.ll.MoveToFront(elem)
		return kv.Value, true
	}
	return "", false
}

func (c *Cache) Set(key, value string, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, found := c.cacheMap[key]; found {
		c.ll.MoveToFront(elem)
		kv := elem.Value.(*Entry)
		kv.Value = value
		kv.Timestamp = time.Now()
		kv.Expiration = expiration
		return
	}

	if c.ll.Len() == c.capacity {
		oldest := c.ll.Back()
		if oldest != nil {
			c.ll.Remove(oldest)
			kv := oldest.Value.(*Entry)
			delete(c.cacheMap, kv.Key)
		}
	}

	elem := c.ll.PushFront(&Entry{Key: key, Value: value, Timestamp: time.Now(), Expiration: expiration})
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

func (c *Cache) GetCacheData() map[string]map[string]interface{} {
	data := make(map[string]map[string]interface{})
	c.mu.RLock()
	defer c.mu.RUnlock()

	for key, elem := range c.cacheMap {
		kv := elem.Value.(*Entry)
		if time.Since(kv.Timestamp) <= kv.Expiration {
			data[key] = map[string]interface{}{
				"value":      kv.Value,
				"expiration": int(kv.Expiration.Seconds()) - int(time.Since(kv.Timestamp).Seconds()),
			}
		}
	}
	return data
}

func main() {
	cache := NewCache(10)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, found := cache.Get(key)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"message": "Key not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	r.GET("/cache", func(c *gin.Context) {
		cacheData := cache.GetCacheData()
		c.JSON(http.StatusOK, cacheData)
	})

	r.POST("/cache", func(c *gin.Context) {
		var req struct {
			Key       string `json:"key"`
			Value     string `json:"value"`
			ExpiresIn int    `json:"expires_in"` // Expecting seconds
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}
		expiration := time.Duration(req.ExpiresIn) * time.Second
		cache.Set(req.Key, req.Value, expiration)
		c.JSON(http.StatusOK, gin.H{"message": "Key added/updated"})
	})

	r.DELETE("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		cache.Delete(key)
		c.JSON(http.StatusOK, gin.H{"message": "Key deleted"})
	})

	r.GET("/cache-updates", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upgrade connection"})
			return
		}
		defer conn.Close()

		for {
			time.Sleep(5 * time.Second)
			cacheData := cache.GetCacheData()
			err := conn.WriteJSON(cacheData)
			if err != nil {
				return
			}
		}
	})

	r.Run(":8080")
}
