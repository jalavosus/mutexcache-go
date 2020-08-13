package mutexcache_test

import (
	"github.com/jalavosus/mutexcache-go"
	"testing"
	"time"
)

var (
	defaultExpiration = 30 * time.Second
	testKeyA          = "test_key_a"
	testKeyB          = "test_key_b"
)

// Test creation and retrieval of a single *sync.Mutex
func TestSingle(t *testing.T) {
	mutexCache := mutexcache.New(defaultExpiration)

	mut := mutexCache.Get(testKeyA)
	if mut == nil {
		t.Fail()
	}

	// Ensure that the mutex pointer is actually being stored,
	// and not recreated on subsequent calls.
	sameMut := mutexCache.Get(testKeyA)
	if sameMut == nil {
		t.Fail()
	}

	if sameMut != mut {
		t.Fail()
	}
}

// Tests storage and retrieval of multiple *sync.Mutexes
func TestMulti(t *testing.T) {
	mutexCache := mutexcache.New(defaultExpiration)

	mutA := mutexCache.Get(testKeyA)
	if mutA == nil {
		t.Fail()
	}

	mutB := mutexCache.Get(testKeyB)
	if mutB == nil {
		t.Fail()
	}

	// Ensure that the cache isn't returning
	// the same mutex for some reason.
	if mutA == mutB {
		t.Fail()
	}
}

// Tests cache expiration.
func TestExpiration(t *testing.T) {
	mutexCache := mutexcache.New(defaultExpiration)

	mut := mutexCache.Get(testKeyA)

	// wait for the mutex to expire and the
	// internal cache to prune it
	time.Sleep(31 * time.Second)

	newMut := mutexCache.Get(testKeyA)

	if newMut == mut {
		t.Fail()
	}
}
