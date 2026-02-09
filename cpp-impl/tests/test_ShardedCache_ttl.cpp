#include <cassert>
#include <chrono>
#include <string>
#include <thread>

#include "ShardedCache.hpp"

static void testTtlExpiration() {
    ShardedCache cache;
    cache.set("temp", "value", 1);
    std::this_thread::sleep_for(std::chrono::seconds(2));
    assert(cache.get("temp") == "");
}

static void testTtlIndependentKeys() {
    ShardedCache cache;
    cache.set("short", "s", 1);
    cache.set("long", "l", 5);
    std::this_thread::sleep_for(std::chrono::seconds(2));
    assert(cache.get("short") == "");
    assert(cache.get("long") == "l");
}

int main() {
    testTtlExpiration();
    testTtlIndependentKeys();
    return 0;
}
