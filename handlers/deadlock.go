package handlers

import (
	"net/http"
	"sync"
)

// CacheManager coordinates cache updates and invalidations using two mutexes.
// The bug: UpdateCache locks mu1 then mu2; InvalidateCache locks mu2 then mu1.
// When both run concurrently, they can deadlock. This produces a fatal error
// that cannot be recovered by the recovery middleware.
type CacheManager struct {
	mu1   sync.Mutex
	mu2   sync.Mutex
	cache map[string]string
}

// NewCacheManager returns a cache manager with empty cache.
func NewCacheManager() *CacheManager {
	return &CacheManager{cache: make(map[string]string)}
}

// UpdateCache acquires mu1 then mu2 and updates the cache.
func (c *CacheManager) UpdateCache(key, value string) {
	c.mu1.Lock()
	defer c.mu1.Unlock()
	c.mu2.Lock()
	defer c.mu2.Unlock()
	c.cache[key] = value
}

// InvalidateCache acquires mu2 then mu1 (reversed order) and clears an entry.
// With UpdateCache running in another goroutine, this order causes deadlock.
func (c *CacheManager) InvalidateCache(key string) {
	c.mu2.Lock()
	defer c.mu2.Unlock()
	c.mu1.Lock()
	defer c.mu1.Unlock()
	delete(c.cache, key)
}

// Deadlock handles GET /error/deadlock.
// It runs UpdateCache and InvalidateCache concurrently to trigger deadlock.
// Note: Fatal error "all goroutines are asleep - deadlock" is not recoverable;
// the process will exit and recovery middleware will not run.
func Deadlock(mgr *CacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go mgr.UpdateCache("foo", "bar")
		mgr.InvalidateCache("foo")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
