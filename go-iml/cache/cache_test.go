package cache

import (
	"testing"
	"time"
)

func TestCacheSetGetDelete(t *testing.T) {
	c := NewCache(4)
	defer c.Close()

	c.Set("k1", "v1", time.Second)
	if got, ok := c.Get("k1"); !ok || got != "v1" {
		t.Fatalf("expected v1, got %q, ok=%v", got, ok)
	}

	c.Delete("k1")
	if _, ok := c.Get("k1"); ok {
		t.Fatalf("expected key to be deleted")
	}
}

func TestCacheTTLExpiry(t *testing.T) {
	c := NewCache(2)
	defer c.Close()

	c.Set("k1", "v1", 50*time.Millisecond)
	time.Sleep(80 * time.Millisecond)

	if _, ok := c.Get("k1"); ok {
		t.Fatalf("expected key to expire")
	}
}

func TestCacheShardCountPowerOfTwo(t *testing.T) {
	c1 := NewCache(3)
	defer c1.Close()
	if got := len(c1.shards); got != 4 {
		t.Fatalf("expected 4 shards, got %d", got)
	}

	c2 := NewCache(8)
	defer c2.Close()
	if got := len(c2.shards); got != 8 {
		t.Fatalf("expected 8 shards, got %d", got)
	}

	c3 := NewCache(1)
	defer c3.Close()
	if got := len(c3.shards); got != 2 {
		t.Fatalf("expected 2 shards, got %d", got)
	}
}
