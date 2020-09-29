// Package whalelruttl proxy Cache to github.com/hashicorp/golang-lru
// With TTL control.
package whalelruttl

import (
	"sync"
	"time"

	"github.com/farwydi/cleanwhale/cache"
	lru "github.com/hashicorp/golang-lru"
)

type ttlController struct {
	sync.Mutex

	Expiry time.Time
	Value  interface{}
}

// NewCacheLRUWithTTL make cache.Cache object.
// Initial lur cache with maxSize and default defaultTTL.
func NewCacheLRUWithTTL(maxSize int, defaultTTL time.Duration) (cache.Cache, error) {
	c, err := lru.New(maxSize)
	if err != nil {
		return nil, err
	}
	return &lurWithTTL{
		defaultTTL: defaultTTL,
		cache:      c,
	}, nil
}

type lurWithTTL struct {
	defaultTTL time.Duration
	cache      *lru.Cache
}

func (c *lurWithTTL) checkTTLOrRemove(key string, value *ttlController) (interface{}, bool) {
	value.Lock()
	defer value.Unlock()
	isActive := value.Expiry.Before(time.Now())
	if isActive {
		return value.Value, true
	}
	c.cache.Remove(key)
	return nil, false
}

func (c *lurWithTTL) Add(key string, value interface{}) {
	c.cache.Add(key, ttlController{Expiry: time.Now().Add(c.defaultTTL), Value: value})
}

func (c *lurWithTTL) Has(key string) bool {
	value, found := c.cache.Get(key)
	if !found {
		return false
	}
	_, ok := c.checkTTLOrRemove(key, value.(*ttlController))
	return ok
}

func (c *lurWithTTL) Get(key string) (value interface{}, found bool) {
	value, found = c.cache.Get(key)
	if !found {
		return nil, false
	}
	return c.checkTTLOrRemove(key, value.(*ttlController))
}
