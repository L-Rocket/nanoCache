package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"nanoCache/go-iml/cache"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() (*gin.Engine, *cache.Cache) {
	// Use gin test mode to reduce log noise
	gin.SetMode(gin.TestMode)
	c := cache.NewCache(4)
	r := SetupRouter(c)
	return r, c
}

func TestSetCacheHandler(t *testing.T) {
	r, c := setupTestRouter()
	defer c.Close()

	// Test Case 1: Successful Set
	payload := map[string]interface{}{
		"key":   "testKey",
		"value": "testValue",
		"ttl":   60,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/cache/set", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify the value was actually set in the cache
	val, ok := c.Get("testKey")
	if !ok {
		t.Error("Expected key 'testKey' to be in cache")
	}
	if val != "testValue" {
		t.Errorf("Expected value 'testValue', got '%s'", val)
	}

	// Test Case 2: Bad Request (Missing fields)
	badPayload := map[string]interface{}{
		"key": "testKey2",
		// missing value and ttl
	}
	badBody, _ := json.Marshal(badPayload)
	reqBad, _ := http.NewRequest("POST", "/cache/set", bytes.NewBuffer(badBody))
	reqBad.Header.Set("Content-Type", "application/json")
	wBad := httptest.NewRecorder()
	r.ServeHTTP(wBad, reqBad)

	if wBad.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for bad payload, got %d", wBad.Code)
	}
}

func TestGetCacheHandler(t *testing.T) {
	r, c := setupTestRouter()
	defer c.Close()

	// Pre-populate cache
	c.Set("existingKey", "existingValue", 60*time.Second)

	// Test Case 1: Get Existing Key
	req, _ := http.NewRequest("GET", "/cache/existingKey", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}
	if response["value"] != "existingValue" {
		t.Errorf("Expected value 'existingValue', got '%s'", response["value"])
	}

	// Test Case 2: Get Non-Existent Key
	reqNotFound, _ := http.NewRequest("GET", "/cache/nonExistentKey", nil)
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for missing key, got %d", wNotFound.Code)
	}
}

func TestDeleteCacheHandler(t *testing.T) {
	r, c := setupTestRouter()
	defer c.Close()

	// Pre-populate cache
	c.Set("keyToDelete", "valueToDelete", 60*time.Second)

	// Test Case 1: Delete Existing Key
	req, _ := http.NewRequest("DELETE", "/cache/keyToDelete", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify key is gone
	_, ok := c.Get("keyToDelete")
	if ok {
		t.Error("Expected key 'keyToDelete' to be deleted from cache")
	}

	// Test Case 2: Delete Non-Existent Key (Should still be OK usually, or logic dependent)
	// The current controller logic just calls delete and returns OK regardless.
	reqDeleteAgain, _ := http.NewRequest("DELETE", "/cache/nonExistentKey", nil)
	wDeleteAgain := httptest.NewRecorder()
	r.ServeHTTP(wDeleteAgain, reqDeleteAgain)

	if wDeleteAgain.Code != http.StatusOK {
		t.Errorf("Expected status 200 even for non-existent key deletion, got %d", wDeleteAgain.Code)
	}
}

func TestConcurrentServerRequests(t *testing.T) {
	r, c := setupTestRouter()
	defer c.Close()

	var wg sync.WaitGroup
	var successCount int64
	numGoroutines := 50
	opsPerGoroutine := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))

			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d", rng.Intn(100)) // Shared keys to cause contention
				op := rng.Intn(100)

				if op < 40 { // 40% Writes
					payload := map[string]interface{}{
						"key":   key,
						"value": fmt.Sprintf("val-%d", j),
						"ttl":   10,
					}
					body, _ := json.Marshal(payload)
					req, _ := http.NewRequest("POST", "/cache/set", bytes.NewBuffer(body))
					req.Header.Set("Content-Type", "application/json")
					w := httptest.NewRecorder()
					r.ServeHTTP(w, req)
					if w.Code == http.StatusOK {
						atomic.AddInt64(&successCount, 1)
					}

				} else if op < 80 { // 40% Reads
					req, _ := http.NewRequest("GET", "/cache/"+key, nil)
					w := httptest.NewRecorder()
					r.ServeHTTP(w, req)
					// We expect 200 or 404, both are valid successes in a concurrent mixed load
					if w.Code == http.StatusOK || w.Code == http.StatusNotFound {
						atomic.AddInt64(&successCount, 1)
					}

				} else { // 20% Deletes
					req, _ := http.NewRequest("DELETE", "/cache/"+key, nil)
					w := httptest.NewRecorder()
					r.ServeHTTP(w, req)
					if w.Code == http.StatusOK {
						atomic.AddInt64(&successCount, 1)
					}
				}
			}
		}(i)
	}

	wg.Wait()
	t.Logf("Completed %d concurrent HTTP requests", successCount)
	if successCount != int64(numGoroutines*opsPerGoroutine) {
		t.Errorf("Expected %d successful ops, got %d", numGoroutines*opsPerGoroutine, successCount)
	}
}