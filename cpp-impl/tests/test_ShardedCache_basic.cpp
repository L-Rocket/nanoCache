#include <cassert>
#include <string>

#include "ShardedCache.hpp"

static void testSetAndGet() {
    ShardedCache cache;
    cache.set("key1", "value1");
    assert(cache.get("key1") == "value1");
}

static void testOverwrite() {
    ShardedCache cache;
    cache.set("key1", "value1");
    cache.set("key1", "value2");
    assert(cache.get("key1") == "value2");
}

static void testCustomShardCountConstructor() {
    ShardedCache cache(5);
    cache.set("alpha", "one");
    cache.set("beta", "two");
    assert(cache.get("alpha") == "one");
    assert(cache.get("beta") == "two");
}

int main() {
    testSetAndGet();
    testOverwrite();
    testCustomShardCountConstructor();
    return 0;
}
