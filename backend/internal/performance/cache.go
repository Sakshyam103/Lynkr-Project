/**
 * Cache Layer
 * In-memory caching for improved API response times
 */

package performance

import (
	// "encoding/json"
	"sync"
	"time"
)

type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

type Cache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
}

func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]CacheItem),
	}

	// Start cleanup routine
	go cache.cleanup()

	return cache
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]CacheItem)
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

func (c *Cache) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	expired := 0
	now := time.Now()
	for _, item := range c.items {
		if now.After(item.ExpiresAt) {
			expired++
		}
	}

	return map[string]interface{}{
		"total_items":   len(c.items),
		"expired_items": expired,
		"active_items":  len(c.items) - expired,
	}
}

// Specialized cache methods for common use cases
func (c *Cache) CacheEventData(eventID string, data interface{}) {
	c.Set("event:"+eventID, data, 10*time.Minute)
}

func (c *Cache) GetEventData(eventID string) (interface{}, bool) {
	return c.Get("event:" + eventID)
}

func (c *Cache) CacheUserData(userID string, data interface{}) {
	c.Set("user:"+userID, data, 5*time.Minute)
}

func (c *Cache) GetUserData(userID string) (interface{}, bool) {
	return c.Get("user:" + userID)
}

func (c *Cache) CacheAnalytics(key string, data interface{}) {
	c.Set("analytics:"+key, data, 15*time.Minute)
}

func (c *Cache) GetAnalytics(key string) (interface{}, bool) {
	return c.Get("analytics:" + key)
}

func (c *Cache) InvalidatePattern(pattern string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key := range c.items {
		if contains(key, pattern) {
			delete(c.items, key)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
