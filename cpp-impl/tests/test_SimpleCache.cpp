#include <cassert>
#include <string>

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

int main() {
    testSetAndGet();
    testGetNonexistent();
    testDelete();
    testOverwrite();
    return 0;
}
