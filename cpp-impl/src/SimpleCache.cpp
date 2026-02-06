#include "SimpleCache.hpp"
#include <mutex>
#include <shared_mutex>

SimpleCache::SimpleCache() = default;
SimpleCache::~SimpleCache() = default;



void SimpleCache::set(const std::string& key, const std::string& value) {
    std::unique_lock<std::shared_mutex> lock(cache_mutex);
    data_map[key] = value;
}

std::string SimpleCache::get(const std::string& key) {
    std::shared_lock<std::shared_mutex> lock(cache_mutex);
    auto it = data_map.find(key);
    if (it != data_map.end()) {
        return it->second;
    }
    return ""; // Return empty string if key not found
}

void SimpleCache::del(const std::string& key) {
    std::unique_lock<std::shared_mutex> lock(cache_mutex);
    data_map.erase(key);
}