#include "ShardedCache.hpp"
#include <chrono>
#include <cstddef>


ShardedCache::ShardedCache() : shards(SHARD_COUNT), running_(true) {
    for (size_t i = 0; i < SHARD_COUNT; ++i) {
        shards[i] = std::make_unique<CacheShard>();
    }
    janitor_thread_ = std::thread(&ShardedCache::cleanup_expired_entries, this);
}


ShardedCache::ShardedCache(size_t num_shards) :  running_(true){
    size_t shard_count = 2;
    while (shard_count < num_shards) {
        shard_count <<= 1;
    }
    shards = std::vector<std::unique_ptr<CacheShard>>(shard_count);
    for (size_t i = 0; i < shard_count; ++i) {
        shards[i] = std::make_unique<CacheShard>();
    }
    janitor_thread_ = std::thread(&ShardedCache::cleanup_expired_entries, this);
}

ShardedCache::~ShardedCache() {
    running_ = false;
    if (janitor_thread_.joinable()) {
        janitor_thread_.join();
    }
}


inline uint64_t ShardedCache::fnv1a_hash(const std::string& key) const {
    const uint64_t fnv_prime = 0x100000001b3;
    uint64_t hash = 0xcbf29ce484222325;
    for (char c : key) {
        hash ^= static_cast<uint64_t>(c);
        hash *= fnv_prime;
    }
    return hash;    
}

inline size_t ShardedCache::get_shard_index(const std::string& key) const {
    return fnv1a_hash(key) & (shards.size() - 1);
}

void ShardedCache::cleanup_expired_entries() {
    while (running_) {
        std::this_thread::sleep_for(std::chrono::seconds(1));
        if (!running_) {
            break;
        }
        for (auto& shard : shards) {
            if (shard) {
                shard->cleanup_expired_entries();
            }
        }
    }
}



void ShardedCache::set(const std::string& key, const std::string& value, int ttl_second) {
    size_t shard_index = get_shard_index(key);
    shards[shard_index]->set(key, value, ttl_second);
}

std::string ShardedCache::get(const std::string& key) {
    size_t shard_index = get_shard_index(key);
    return shards[shard_index]->get(key);
}

void ShardedCache::del(const std::string& key) {
    size_t shard_index = get_shard_index(key);
    shards[shard_index]->del(key);
}