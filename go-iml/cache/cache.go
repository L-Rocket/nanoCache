package cache

import (
	"time"
)

type cache struct {
	shards []*shard
	stop   chan struct{}
}

func newCache(numShards int) *cache {

	cnt := 2
	for cnt < numShards {
		cnt <<= 1
	}

	c := &cache{
		shards: make([]*shard, cnt),
		stop:   make(chan struct{}),
	}

	for i := 0; i < cnt; i++ {
		c.shards[i] = newShard()
	}

	go c.janitor_clean()

	return c
}

func (c *cache) janitor_clean() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			{
				for _, shard := range c.shards {
					shard.cleanup()
				}
			}

		case <-c.stop:
			{
				return
			}
		}
	}
}

func (c *cache) Close() {
	close(c.stop)
}

func (c *cache) getIndex(key string) uint64 {
	return fnv1a_hash(key) & (uint64(len(c.shards)) - 1)
}

func (c *cache) Set(key string, value string, ttl time.Duration) {
	index := c.getIndex(key)
	c.shards[index].set(key, value, ttl)
}

func (c *cache) Get(key string) (string, bool) {
	index := c.getIndex(key)
	return c.shards[index].get(key)
}

func (c *cache) Delete(key string) {
	index := c.getIndex(key)
	c.shards[index].delete(key)
}

func fnv1a_hash(key string) uint64 {
	var fnv_prime uint64 = 0x100000001b3
	var hash uint64 = 0xcbf29ce484222325
	for _, c := range key {
		hash ^= uint64(c)
		hash *= fnv_prime
	}
	return hash
}
