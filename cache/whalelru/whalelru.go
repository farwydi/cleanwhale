// Package whalelru proxy Cache to github.com/hashicorp/golang-lru
package whalelru

import (
	"github.com/farwydi/cleanwhale/cache"
	lru "github.com/hashicorp/golang-lru"
)

// NewCacheLRU make cache.Cache object.
// Initial lur cache with maxSize.
func NewCacheLRU(maxSize int) (cache.Cache, error) {
	c, err := lru.New(maxSize)
	if err != nil {
		return nil, err
	}
	return &lurWithTTL{
		cache: c,
	}, nil
}

type lurWithTTL struct {
	cache *lru.Cache
}

func (c *lurWithTTL) Add(key string, value interface{}) {
	c.cache.Add(key, value)
}

func (c *lurWithTTL) Has(key string) bool {
	return c.cache.Contains(key)
}

func (c *lurWithTTL) Get(key string) (value interface{}, found bool) {
	return c.cache.Get(key)
}
