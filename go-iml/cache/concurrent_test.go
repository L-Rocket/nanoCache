package cache

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCacheConcurrentSetGet(t *testing.T) {
	c := newCache(8)
	defer c.Close()

	var wg sync.WaitGroup
	numGoroutines := 100
	opsPerGoroutine := 1000

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := fmt.Sprintf("value-%d-%d", id, j)
				c.Set(key, value, time.Second)
				if got, ok := c.Get(key); !ok || got != value {
					t.Errorf("concurrent set/get failed: key=%s, expected=%s, got=%s, ok=%v", key, value, got, ok)
				}
			}
		}(i)
	}
	wg.Wait()
	t.Logf("Successfully completed %d goroutines with %d ops each", numGoroutines, opsPerGoroutine)
}

func TestCacheConcurrentDeleteGet(t *testing.T) {
	c := newCache(4)
	defer c.Close()

	var wg sync.WaitGroup
	numGoroutines := 50
	opsPerGoroutine := 500

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d", j%100) // Intentionally create conflicts
				value := fmt.Sprintf("value-%d", j)
				c.Set(key, value, time.Second)

				if id%2 == 0 {
					c.Get(key)
				} else {
					c.Delete(key)
				}
			}
		}(i)
	}
	wg.Wait()
	t.Logf("Successfully handled concurrent deletes and gets")
}

func TestShardConcurrentSetGet(t *testing.T) {
	s := newShard()

	var wg sync.WaitGroup
	numGoroutines := 50
	opsPerGoroutine := 2000

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("shard-key-%d-%d", id, j)
				value := fmt.Sprintf("shard-value-%d-%d", id, j)
				s.set(key, value, time.Minute)
				if got, ok := s.get(key); !ok || got != value {
					t.Errorf("shard concurrent set/get failed: key=%s", key)
				}
			}
		}(i)
	}
	wg.Wait()
	t.Logf("Shard: successfully handled %d goroutines", numGoroutines)
}

func TestCacheConcurrentMixedOperations(t *testing.T) {
	c := newCache(16)
	defer c.Close()

	var wg sync.WaitGroup
	var successCount int64
	numGoroutines := 200
	opsPerGoroutine := 500

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("mixed-%d", j%500)
				op := (id + j) % 4

				switch op {
				case 0: // Set
					c.Set(key, fmt.Sprintf("val-%d-%d", id, j), time.Second*5)
					atomic.AddInt64(&successCount, 1)
				case 1: // Get
					if _, ok := c.Get(key); ok {
						atomic.AddInt64(&successCount, 1)
					}
				case 2: // Delete
					c.Delete(key)
					atomic.AddInt64(&successCount, 1)
				case 3: // Get again
					c.Get(key)
					atomic.AddInt64(&successCount, 1)
				}
			}
		}(i)
	}
	wg.Wait()
	t.Logf("Completed %d mixed operations without panic", atomic.LoadInt64(&successCount))
}

func TestCacheConcurrentWithJanitor(t *testing.T) {
	c := newCache(8)
	defer c.Close()

	var wg sync.WaitGroup
	numGoroutines := 30
	opsPerGoroutine := 200

	// Let janitor run and potentially clean while we're adding/reading
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("ttl-key-%d-%d", id, j)
				ttl := time.Duration((j%3)+1) * 100 * time.Millisecond
				c.Set(key, "value", ttl)
				time.Sleep(time.Duration((id*j)%50) * time.Millisecond)
				c.Get(key) // May or may not exist due to expiry
			}
		}(i)
	}
	wg.Wait()
	t.Logf("Successfully handled concurrent operations with janitor cleanup")
}

func BenchmarkCacheSet(b *testing.B) {
	c := newCache(16)
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench-key-%d", i)
		c.Set(key, "value", time.Second)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	c := newCache(16)
	defer c.Close()

	c.Set("bench-key", "value", time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("bench-key")
	}
}

func BenchmarkCacheConcurrentSetGet(b *testing.B) {
	c := newCache(16)
	defer c.Close()

	b.ResetTimer()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < b.N/16; j++ {
				key := fmt.Sprintf("bench-%d-%d", id, j)
				c.Set(key, "value", time.Second)
				c.Get(key)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkShardConcurrentSetGet(b *testing.B) {
	s := newShard()

	b.ResetTimer()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < b.N/16; j++ {
				key := fmt.Sprintf("shard-bench-%d-%d", id, j)
				s.set(key, "value", time.Minute)
				s.get(key)
			}
		}(i)
	}
	wg.Wait()
}
