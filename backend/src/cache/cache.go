package cache

import (
	"context"
	"fmt"
	"luna-backend/constants"
	"luna-backend/errors"
	"luna-backend/types"
	"sync"
	"time"
)

// This module lets us reduce the amounts of requests to upstreams by caching
// source and calendar objects in-between calls to Luna API.
//
// Without this, we would for example have to refetch the source for every
// calendar that we try to fetch, since it is done with a separate API call for
// each calendar.
//
// This cache is very short-lived because it is only meant to handle short
// "bursts" of API calls, such as clicking the refresh button in the frontend.
//
// To be absolute certain that sensitive data does not make its way from one
// user to another, we include the calling user's key in the key and the value
// of each cache element.
//
// Other than that, we rely on UUID's entropy to guarantee we don't run into
// type casting errors. To reduce the chance of collisions further, we could
// include the type of the cache entry like "source" or "calendar" in the key,
// but is the hash function's collision resillience really stronger than UUID's?
// To be absolutely certain, we could also add the cache entry type name to the
// actual cache entry value and perform a check as we do with the user ID, but I
// really doubt it is worth the additional overhead, when the chances of such a
// collision are astronomically low. Why do I do it for user IDs then? I think
// that the potential of some user's resources being exposed to an adversary is
// enough to justify a little bit of paranoia, so I prefer to bet on determinism
// rather than probability in this case, even if it means a few more comparisons
// need to be made by the CPU. In any case, this is still faster than querying
// the upstream.

type Cacheable interface {
	GetId() types.ID
	SupplyContext(ctx context.Context)
}

type cacheEntry struct {
	item      Cacheable
	userId    types.ID
	timestamp time.Time
}

func (entry *cacheEntry) expired() bool {
	return time.Since(entry.timestamp) > constants.LifetimeRamCache
}

type Cache struct {
	dictionary map[string]cacheEntry
	lock       *sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		dictionary: make(map[string]cacheEntry),
	}
}

func getCacheKey(userId types.ID, objectId types.ID) string {
	return fmt.Sprintf("%v;%v", userId.String(), objectId.String())
}

// This will be possible in Go 1.27 (https://github.com/golang/go/issues/77273)
//func (cache *Cache) GetCached[T any](userId types.ID, objectId types.ID, fallback func () T) T {
//
//}

func GetCached[T Cacheable](cache *Cache, userId types.ID, objectId types.ID, ctx context.Context, fallback func() (T, *errors.ErrorTrace)) (T, *errors.ErrorTrace) {
	cache.lock.RLock()
	obj, exists := cache.dictionary[getCacheKey(userId, objectId)]
	cache.lock.RUnlock()

	if exists && !obj.expired() && obj.userId == userId {
		obj.item.SupplyContext(ctx)
		return obj.item.(T), nil
	}

	cache.lock.Lock()
	defer cache.lock.Unlock()
	return fallback()
}

func (cache *Cache) Cache(userId types.ID, object Cacheable) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	cache.dictionary[getCacheKey(userId, object.GetId())] = cacheEntry{
		item:      object,
		timestamp: time.Now(),
	}
}

func (cache *Cache) DeleteStaleEntries() {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	for key, entry := range cache.dictionary {
		if entry.expired() {
			delete(cache.dictionary, key)
		}
	}
}
