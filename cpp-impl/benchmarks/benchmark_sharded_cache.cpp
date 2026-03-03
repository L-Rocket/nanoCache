#include <chrono>
#include <iostream>
#include <string>
#include <thread>
#include <vector>

#include "ShardedCache.hpp"

using Clock = std::chrono::steady_clock;

struct Result {
    std::string name;
    double seconds;
    size_t operations;
};

std::vector<std::string> build_keys(size_t operations) {
    std::vector<std::string> keys;
    keys.reserve(operations);
    for (size_t i = 0; i < operations; ++i) {
        keys.push_back("k" + std::to_string(i));
    }
    return keys;
}

Result run_set_benchmark(size_t operations) {
    ShardedCache cache;
    const auto keys = build_keys(operations);

    auto start = Clock::now();
    for (size_t i = 0; i < operations; ++i) {
        cache.set(keys[i], "v");
    }
    auto end = Clock::now();
    std::chrono::duration<double> diff = end - start;
    return {"cpp_set", diff.count(), operations};
}

Result run_get_benchmark(size_t operations) {
    ShardedCache cache;
    const auto keys = build_keys(operations);

    for (size_t i = 0; i < operations; ++i) {
        cache.set(keys[i], "v");
    }

    auto start = Clock::now();
    for (size_t i = 0; i < operations; ++i) {
        (void)cache.get(keys[i]);
    }
    auto end = Clock::now();
    std::chrono::duration<double> diff = end - start;
    return {"cpp_get", diff.count(), operations};
}

Result run_concurrent_set_get_benchmark(size_t operations, size_t threads) {
    ShardedCache cache;
    std::vector<std::thread> workers;
    workers.reserve(threads);

    const size_t per_thread = operations / threads;
    std::vector<std::vector<std::string>> keys(threads);
    for (size_t t = 0; t < threads; ++t) {
        auto& thread_keys = keys[t];
        thread_keys.reserve(per_thread);
        for (size_t i = 0; i < per_thread; ++i) {
            thread_keys.push_back("k" + std::to_string(t) + "_" + std::to_string(i));
        }
    }

    auto start = Clock::now();
    for (size_t t = 0; t < threads; ++t) {
        workers.emplace_back([t, per_thread, &cache, &keys]() {
            for (size_t i = 0; i < per_thread; ++i) {
                const std::string& key = keys[t][i];
                cache.set(key, "v");
                (void)cache.get(key);
            }
        });
    }

    for (auto& w : workers) {
        w.join();
    }

    auto end = Clock::now();
    std::chrono::duration<double> diff = end - start;
    return {"cpp_concurrent_set_get", diff.count(), per_thread * threads * 2};
}

void print_result(const Result& result) {
    const double ops_per_second = result.operations / result.seconds;
    std::cout << result.name
              << " seconds=" << result.seconds
              << " operations=" << result.operations
              << " ops_per_second=" << ops_per_second
              << '\n';
}

int main(int argc, char** argv) {
    size_t operations = 200000;
    size_t threads = 16;

    for (int i = 1; i < argc; ++i) {
        std::string arg = argv[i];
        if (arg == "--ops" && i + 1 < argc) {
            operations = static_cast<size_t>(std::stoull(argv[++i]));
        } else if (arg == "--threads" && i + 1 < argc) {
            threads = static_cast<size_t>(std::stoull(argv[++i]));
        }
    }

    print_result(run_set_benchmark(operations));
    print_result(run_get_benchmark(operations));
    print_result(run_concurrent_set_get_benchmark(operations, threads));
    return 0;
}
