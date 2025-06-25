package utils

import (
	"time"

	"github.com/etkecc/go-kit"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

// vCache is a cache for validation results
var vCache = *expirable.NewLRU[string, bool](1000, nil, 1*time.Hour)

// GetCachedValidation retrieves a cached validation result for the given entry.
func GetCachedValidation(entry string) (valid, cached bool) {
	if entry == "" {
		return false, false
	}

	key := kit.Hash(entry)
	return vCache.Get(key)
}

// SetCachedValidation sets a validation result in the cache for the given entry.
func SetCachedValidation(entry string, valid bool) {
	if entry == "" {
		return
	}

	key := kit.Hash(entry)
	vCache.Add(key, valid)
}
