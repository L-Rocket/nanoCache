package cache

import (
	"sync"
	"time"
)

type entry struct {
	key    string
	value  string
	expiry time.Time
}

type shard struct {
	entries map[string]entry
	mutex   sync.RWMutex
}

func newShard() *shard {
	return &shard{
		entries: make(map[string]entry),
	}
}

func (s *shard) set(key string, value string, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.entries[key] = entry{
		key:    key,
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

func (s *shard) get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entry, ok := s.entries[key]
	if !ok {
		return "", false
	}

	if time.Now().After(entry.expiry) {
		return "", false
	}

	return entry.value, true
}

func (s *shard) delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.entries, key)
}

func (s *shard) cleanup() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for key, entry := range s.entries {
		if now.After(entry.expiry) {
			delete(s.entries, key)
		}
	}
}
