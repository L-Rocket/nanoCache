package cache

import (
	"sync"
	"time"
)

type entry struct {
	key    string
	expiry time.Time
}

type shard struct {
	entries map[string]entry
	mutex   sync.Mutex
}


func newShard() *shard {
	return &shard{
		entries: make(map[string]entry),
	}
}

func (s )