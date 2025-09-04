package storage

import (
	"sync"
	"sync/atomic"

	xxhash "github.com/cespare/xxhash/v2"
)

const (
	DataFile  = "elysiandb.json"
	numShards = 128
	shardMask = numShards - 1
)

type shard struct {
	mu sync.RWMutex
	m  map[string][]byte
}

type Store struct {
	shards [numShards]*shard
	saved  atomic.Bool
}

func newStore() *Store {
	s := &Store{}
	for i := 0; i < numShards; i++ {
		s.shards[i] = &shard{m: make(map[string][]byte)}
	}
	s.saved.Store(true)
	return s
}

func (s *Store) reset() {
	for i := 0; i < numShards; i++ {
		sh := s.shards[i]
		sh.mu.Lock()
		sh.m = make(map[string][]byte)
		sh.mu.Unlock()
	}
	s.saved.Store(false)
}

func shardIndex(key string) int {
	h := xxhash.Sum64String(key)
	return int(h & shardMask)
}

func (s *Store) get(key string) ([]byte, bool) {
	sh := s.shards[shardIndex(key)]
	sh.mu.RLock()
	v, ok := sh.m[key]
	sh.mu.RUnlock()
	if !ok {
		return nil, false
	}

	return v, true
}

func (s *Store) put(key string, value []byte) {
	buf := make([]byte, len(value))
	copy(buf, value)

	sh := s.shards[shardIndex(key)]
	sh.mu.Lock()
	sh.m[key] = buf
	sh.mu.Unlock()
	s.saved.Store(false)
}

func (s *Store) del(key string) {
	sh := s.shards[shardIndex(key)]
	sh.mu.Lock()
	delete(sh.m, key)
	sh.mu.Unlock()
	s.saved.Store(false)
}

func (s *Store) Iterate(fn func(k string, v []byte)) {
	for i := 0; i < numShards; i++ {
		sh := s.shards[i]
		sh.mu.RLock()
		for k, v := range sh.m {
			c := make([]byte, len(v))
			copy(c, v)
			fn(k, c)
		}

		sh.mu.RUnlock()
	}
}

func (s *Store) FromMap(src map[string][]byte) {
	for k, v := range src {
		sh := s.shards[shardIndex(k)]
		sh.mu.Lock()
		buf := make([]byte, len(v))
		copy(buf, v)
		sh.m[k] = buf
		sh.mu.Unlock()
	}
	s.saved.Store(true)
}

func (s *Store) ToMap() map[string][]byte {
	result := make(map[string][]byte)
	s.Iterate(func(k string, v []byte) {
		result[k] = v
	})

	return result
}
