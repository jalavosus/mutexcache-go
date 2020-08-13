# mutexcache-go

A small utility library for dynamically creating mutexes based on cache keys. 

## Use case

Say you're implementing a graphql server, with object fields which asynchronously resolve independently of each other. 
Multiple fields perform the same operation, and thus should use a mutex and some basic caching to ensure that the 
database query happens at most once. If you have an array of these objects, suddenly they're all using the same mutex,
which can degrade performance. 

Instead of using one mutex to rule them all, dynamically create multiple short-lived mutexes which each object can use 
independently of other resolving objects. By using the same cache key for the mutexes as you would for your cache check,
you can almost transparently use dynamically created mutexes without worrying about performance or allocation/deallocation
of mutexes. 

With mutexcache, if a mutex associated with a cache key is already stored, then it will be returned. Otherwise, a new 
mutex will silently be created, stored for future use, and returned. 

## Installation

`go get -u github.com/jalavosus/mutexcache-go`

## Usage

```go
package main

import (
    "time"
    "github.com/jalavosus/mutexcache-go"
)

func main() {
    mutexCache := mutexcache.New(30*time.Second) // or whatever default ttl you'd like to use
    
    cacheKeyA := "key_a"
    cacheKeyB := "key_b"
    
    mutA := mutexCache.Get(cacheKeyA)
    // use GetWithExpiration() to specify a ttl other than the default ttl.
    mutB := mutexCache.GetWithExpiration(cacheKeyB, 1*time.Minute)
    
    [...do whatever with your mutexes...]
}
```

## Testing

`go test`

Fair warning, `TestExpiration` runs `time.Sleep()` for 30 seconds to ensure that
cache expiry works. 