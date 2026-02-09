#include <chrono>
#include <cstddef>
#include <string>
#include <unordered_map>
#include <shared_mutex>

class CacheShard {
private:
    class Entry{
    public:
        std::string value;
        std::chrono::steady_clock::time_point timestamp;
    };
    std::unordered_map<std::string, Entry> data_map;
    mutable std::shared_mutex cache_mutex;

public:
    CacheShard();
    ~CacheShard();

    void set(const std::string& key, const std::string& value, int ttl_second = 10);
    std::string get(const std::string& key);
    void del(const std::string& key);

    // interface for chard manager to remove expired entries
    void cleanup_expired_entries();
};