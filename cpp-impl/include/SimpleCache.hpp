#include <atomic>
#include <chrono>
#include <cstddef>
#include <string>
#include <unordered_map>
#include <shared_mutex>
#include <thread>

class SimpleCache {
private:
    class Entry{
    public:
        std::string value;
        std::chrono::steady_clock::time_point timestamp;
    };
    std::unordered_map<std::string, Entry> data_map;
    mutable std::shared_mutex cache_mutex;
    size_t clean_loop_cycle;
    std::atomic<bool> running_;
    std::thread janitor_thread_;

    void cleanup_loop(); // Placeholder for potential cleanup logic
public:
    SimpleCache();
    SimpleCache(size_t clean_loop_cycle);
    ~SimpleCache();

    void set(const std::string& key, const std::string& value, int ttl_second = 10);
    std::string get(const std::string& key);
    void del(const std::string& key);
};