#include <string>
#include <unordered_map>
#include <shared_mutex>

class SimpleCache {
private:
    std::unordered_map<std::string, std::string> data_map;
    mutable std::shared_mutex cache_mutex;
public:
    SimpleCache();
    ~SimpleCache();

    void set(const std::string& key, const std::string& value);
    std::string get(const std::string& key);

    void del(const std::string& key);
};