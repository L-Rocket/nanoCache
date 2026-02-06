#include "SimpleCache.hpp"
#include <chrono>
#include <cstddef>
#include <mutex>
#include <shared_mutex>

SimpleCache::SimpleCache() : clean_loop_cycle(1000), running_(true) {
    janitor_thread_ = std::thread(&SimpleCache::cleanup_loop, this);
}


SimpleCache::SimpleCache(size_t clean_loop_cycle) : clean_loop_cycle(clean_loop_cycle), running_(true) {
    janitor_thread_ = std::thread(&SimpleCache::cleanup_loop, this);
}

SimpleCache::~SimpleCache() {
    running_ = false;
    if (janitor_thread_.joinable()) {
        janitor_thread_.join();
    }
}


void SimpleCache::cleanup_loop() {
    size_t cycle_count = 0;
    while (running_) {
        std::this_thread::sleep_for(std::chrono::seconds(1));

        if (!running_) {
            break;
        }

        cycle_count++;
        if (cycle_count >= clean_loop_cycle) {
            cycle_count = 0;

            // Cleanup expired entries
            {
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

        } else {
            continue;
        }
    }
}


void SimpleCache::set(const std::string& key, const std::string& value, int ttl_second) {
    std::unique_lock<std::shared_mutex> lock(cache_mutex);
    data_map[key] = Entry{value, std::chrono::steady_clock::now() + std::chrono::seconds(ttl_second)};
}

std::string SimpleCache::get(const std::string& key) {
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

void SimpleCache::del(const std::string& key) {
    std::unique_lock<std::shared_mutex> lock(cache_mutex);
    data_map.erase(key);
}