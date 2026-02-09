#include <atomic>
#include <cassert>
#include <string>
#include <thread>
#include <vector>

#include "ShardedCache.hpp"

static void testConcurrentWritesAndReads() {
    ShardedCache cache;
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
    ShardedCache cache;
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
    testConcurrentWritesAndReads();
    testConcurrentMixedOperations();
    return 0;
}
