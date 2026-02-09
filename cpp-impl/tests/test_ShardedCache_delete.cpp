#include <cassert>
#include <string>

#include "ShardedCache.hpp"

static void testGetNonexistent() {
    ShardedCache cache;
    assert(cache.get("missing") == "");
}

static void testDelete() {
    ShardedCache cache;
    cache.set("key1", "value1");
    cache.del("key1");
    assert(cache.get("key1") == "");
}

int main() {
    testGetNonexistent();
    testDelete();
    return 0;
}
