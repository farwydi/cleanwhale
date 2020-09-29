// Package cache defines the Cache interface.
package cache

// Cache sets interface for various cache types.
type Cache interface {
	Add(key string, value interface{})
	Has(key string) bool
	Get(key string) (value interface{}, found bool)
}
