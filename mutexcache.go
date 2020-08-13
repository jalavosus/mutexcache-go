package mutexcache

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var cleanupDuration = 30 * time.Second

type MutexCache struct {
	mutCache          *cache.Cache
	defaultExpiration time.Duration
}

// Returns a new instance of MutexCache.
// By default, a MutexCache instance has its
// cache checked and pruned every 30 seconds.
func New(defaultExpiration time.Duration) *MutexCache {
	return &MutexCache{
		mutCache:          cache.New(defaultExpiration, cleanupDuration),
		defaultExpiration: defaultExpiration,
	}
}

// Checks for an existing mutex.
// If found, returns a *sync.Mutex (to prevent lock-copying issues),
// otherwise, creates a new *sync.Mutex,
// stores it using the instances default expiration, and returns that.
func (m *MutexCache) Get(cacheKey string) *sync.Mutex {
	var mut *sync.Mutex

	mutInterface, ok := m.mutCache.Get(cacheKey)
	if !ok {
		mut = m.newMutex(cacheKey, m.defaultExpiration)
	} else {
		mut = mutInterface.(*sync.Mutex)
	}

	return mut
}

// Checks for an existing mutex.
// If found, returns a *sync.Mutex (to prevent lock-copying issues),
// otherwise, creates a new *sync.Mutex,
// stores it using the provided expiration param, and returns that.
func (m *MutexCache) GetWithExpiration(cacheKey string, expiration time.Duration) *sync.Mutex {
	var mut *sync.Mutex

	mutInterface, ok := m.mutCache.Get(cacheKey)
	if !ok {
		mut = m.newMutex(cacheKey, expiration)
	} else {
		mut = mutInterface.(*sync.Mutex)
	}

	return mut
}

func (m *MutexCache) newMutex(cacheKey string, expiration time.Duration) *sync.Mutex {
	mut := &sync.Mutex{}
	m.mutCache.Set(cacheKey, mut, expiration)

	return mut
}
