#include <string>
#include <unordered_map>

class SimpleCache {
private:
    std::unordered_map<std::string, std::string> data_map;

public:
    SimpleCache() = default;

    void set(const std::string& key, const std::string& value) {
        data_map[key] = value;
    }

    std::string get(const std::string& key) {
        auto it = data_map.find(key);
        if (it != data_map.end()) {
            return it->second;
        }
        return ""; // or throw an exception
    }

    void del(const std::string& key) {
        data_map.erase(key);
    }

    ~SimpleCache() = default;
};