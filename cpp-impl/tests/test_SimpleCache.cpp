#include <atomic>
#include <cassert>
#include <string>
#include <thread>
#include <vector>

#include "SimpleCache.hpp"

static void testSetAndGet() {
    SimpleCache cache;
    cache.set("key1", "value1");
    assert(cache.get("key1") == "value1");
}

static void testGetNonexistent() {
    SimpleCache cache;
    assert(cache.get("nonexistent") == "");
}

static void testDelete() {
    SimpleCache cache;
    cache.set("key1", "value1");
    cache.del("key1");
    assert(cache.get("key1") == "");
}

static void testOverwrite() {
    SimpleCache cache;
    cache.set("key1", "value1");
    cache.set("key1", "value2");
    assert(cache.get("key1") == "value2");
}

static void testConcurrentAccess() {
    SimpleCache cache;
    const int thread_count = 8;
    const int keys_per_thread = 200;
    std::vector<std::thread> threads;
    threads.reserve(thread_count);

    for (int t = 0; t < thread_count; ++t) {
        threads.emplace_back([t, keys_per_thread, &cache]() {
            for (int i = 0; i < keys_per_thread; ++i) {
                const std::string key = "k" + std::to_string(t) + "_" + std::to_string(i);
                const std::string value = "v" + std::to_string(t) + "_" + std::to_string(i);
                cache.set(key, value);
            }
        });
    }

    for (auto& th : threads) {
        th.join();
    }

    for (int t = 0; t < thread_count; ++t) {
        for (int i = 0; i < keys_per_thread; ++i) {
            const std::string key = "k" + std::to_string(t) + "_" + std::to_string(i);
            const std::string value = "v" + std::to_string(t) + "_" + std::to_string(i);
            assert(cache.get(key) == value);
        }
    }
}

static void testConcurrentMixedOperations() {
    SimpleCache cache;
    std::atomic<bool> saw_empty{false};
    const int iterations = 1000;

    std::thread writer([&cache, iterations]() {
        for (int i = 0; i < iterations; ++i) {
            cache.set("shared", "value");
        }
    });

    std::thread reader([&cache, &saw_empty, iterations]() {
        for (int i = 0; i < iterations; ++i) {
            if (cache.get("shared").empty()) {
                saw_empty.store(true, std::memory_order_relaxed);
            }
        }
    });

    std::thread deleter([&cache, iterations]() {
        for (int i = 0; i < iterations; ++i) {
            cache.del("shared");
        }
    });

    writer.join();
    reader.join();
    deleter.join();

    cache.set("shared", "value");
    assert(cache.get("shared") == "value");
    (void)saw_empty;
}

int main() {
    testSetAndGet();
    testGetNonexistent();
    testDelete();
    testOverwrite();
    testConcurrentAccess();
    testConcurrentMixedOperations();
    return 0;
}
