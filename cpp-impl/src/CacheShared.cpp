#include "CacheShard.hpp"
#include <chrono>
#include <cstddef>
#include <mutex>
#include <shared_mutex>

CacheShard::CacheShard() = default;


CacheShard::~CacheShard() = default;


void CacheShard::set(const std::string& key, const std::string& value, int ttl_second) {
    std::unique_lock<std::shared_mutex> lock(cache_mutex);
    data_map[key] = Entry{value, std::chrono::steady_clock::now() + std::chrono::seconds(ttl_second)};
}

std::string CacheShard::get(const std::string& key) {
    std::shared_lock<std::shared_mutex> lock(cache_mutex);
    auto it = data_map.find(key);
    if (it != data_map.end()) {
        if (std::chrono::steady_clock::now() > it->second.timestamp) {
            return "";
        }
        return it->second.value;
    }
    return ""; // Return empty string if key not found
}

void CacheShard::del(const std::string& key) {
    std::unique_lock<std::shared_mutex> lock(cache_mutex);
    data_map.erase(key);
}

void CacheShard::cleanup_expired_entries() {
        std::unique_lock<std::shared_mutex> lock(cache_mutex);
        auto now = std::chrono::steady_clock::now();
        for (auto it = data_map.begin(); it != data_map.end(); ) {
            if (now > it->second.timestamp) {
                it = data_map.erase(it);
            } else {
                ++it;
            }
        }            
        
    }