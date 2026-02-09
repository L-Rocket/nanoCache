#include <cstddef>
#include <memory>
#include <string>
#include <vector>
#include <thread>
#include "CacheShard.hpp"

namespace {
    constexpr size_t SHARD_COUNT = 256;
}

class ShardedCache {
private:
    std::vector<std::unique_ptr<CacheShard>> shards;
    std::thread janitor_thread_;
    std::atomic<bool> running_;

    inline uint64_t fnv1a_hash(const std::string& key) const;
    inline size_t get_shard_index(const std::string& key) const;
    void cleanup_expired_entries();

public:
    ShardedCache();
    ShardedCache(size_t num_shards);
    ~ShardedCache();
    
    void set(const std::string& key, const std::string& value, int ttl_second = 10);
    std::string get(const std::string& key);
    void del(const std::string& key);

};