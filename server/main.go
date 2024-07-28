package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var cache *Cache

func main() {
	cache = NewCache(10, 5*time.Second)

	r := gin.Default()

	r.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, found := cache.Get(key)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"message": "Key not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

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

	r.DELETE("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		cache.Delete(key)
		c.JSON(http.StatusOK, gin.H{"message": "Key deleted"})
	})

	r.Run(":8080")
}
