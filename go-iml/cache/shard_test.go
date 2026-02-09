package cache

import (
	"testing"
	"time"
)

func TestShardSetGet(t *testing.T) {
	s := newShard()
	s.set("k1", "v1", time.Second)

	if got, ok := s.get("k1"); !ok || got != "v1" {
		t.Fatalf("expected v1, got %q, ok=%v", got, ok)
	}
}

func TestShardTTLExpiry(t *testing.T) {
	s := newShard()
	s.set("k1", "v1", 30*time.Millisecond)
	time.Sleep(60 * time.Millisecond)

	if _, ok := s.get("k1"); ok {
		t.Fatalf("expected key to expire")
	}
}

func TestShardDelete(t *testing.T) {
	s := newShard()
	s.set("k1", "v1", time.Second)
	s.delete("k1")

	if _, ok := s.get("k1"); ok {
		t.Fatalf("expected key to be deleted")
	}
}

func TestShardCleanupRemovesExpired(t *testing.T) {
	s := newShard()
	s.set("k1", "v1", 10*time.Millisecond)
	s.set("k2", "v2", time.Second)
	time.Sleep(30 * time.Millisecond)

	s.cleanup()

	if _, ok := s.entries["k1"]; ok {
		t.Fatalf("expected k1 to be removed after cleanup")
	}
	if _, ok := s.entries["k2"]; !ok {
		t.Fatalf("expected k2 to remain after cleanup")
	}
}
