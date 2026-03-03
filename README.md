# NanoCache: High-Performance In-Memory Cache

English | [简体中文](./README.zh-CN.md)

![Language](https://img.shields.io/badge/language-C%2B%2B17%20%7C%20Go-blue)
![License](https://img.shields.io/badge/license-MIT-green)

NanoCache is a sharded, thread-safe in-memory key-value store.

This project serves as a **comparative study between Modern C++ (C++17) and Go**, implementing the same architectural design to explore the differences in:
- Concurrency models (OS Threads vs. Goroutines)
- Memory management (RAII/Smart Pointers vs. Garbage Collection)
- Locking strategies and contention profiles

## 🏗 Architecture

To minimize lock contention in high-concurrency scenarios, NanoCache uses a **Sharding Strategy** (similar to BigCache or FreeCache).

- **Sharding:** The keyspace is divided into 256 shards based on FNV-1a hashing.
- **Locking:** Each shard has its own independent `RWMutex` (Go) or `std::shared_mutex` (C++).
- **Eviction:** Support for TTL (Time-To-Live) with both lazy deletion and background cleanup (Janitor).

## 📂 Structure

- **[`cpp-impl/`](./cpp-impl)**: The baseline implementation using C++17. Focuses on manual memory management using `std::shared_ptr` and `std::shared_mutex`.
- **[`go-iml/`](./go-iml)**: The target implementation using Go. Focuses on Goroutines, Channels, and the Go runtime scheduler.


## 🚀 Getting Started

### C++ Implementation

Requirements: CMake >= 3.10, C++17 compliant compiler (GCC/Clang/MSVC).

```bash
cd cpp-impl
mkdir build && cd build
cmake ..
make
./nano_cache_cpp
```

### Go Implementation

Requirements: Go >= 1.20.

```bash
go mod tidy
go run cmd/server/main.go
```

## 📊 Performance Benchmark

Baseline result from the latest local run on **2026-03-03** (Apple M4, darwin/arm64; your machine may differ).

Method:
- Go: `go test ./go-iml/cache -run '^$' -bench 'BenchmarkCache(Set|Get|ConcurrentSetGet)$' -benchmem -count=3`
- C++: `./cpp-impl/build/benchmark_sharded_cache --ops 300000 --threads 16` (3 runs, averaged)

| Scenario | C++ (ops/s) | Go (ops/s) | Faster |
| :--- | :--- | :--- | :--- |
| Set | `7,876,393` | `3,028,890` | `C++ ~2.60x` |
| Get | `10,965,167` | `30,401,036` | `Go ~2.77x` |
| Concurrent Set+Get | `29,095,933` | `4,518,869` | `C++ ~6.44x` |

### Why performance may differ (quick analysis)

1. **Runtime model differences**
   - C++ runs directly on native threads and can keep hot paths close to lower runtime overhead.
   - Go includes goroutine scheduling and GC behavior, which improves concurrency ergonomics but adds runtime cost per operation.

2. **Allocation and memory management differences**
   - C++ object lifetime is more explicitly controlled.
   - Go `Set`/concurrent paths may introduce allocation and GC pressure (also visible from `allocs/op` in benchmark output).

3. **Locking and contention behavior**
   - Both implementations use sharding + RW locks.
   - Different lock implementations and scheduler behavior can change throughput/latency under contention.

4. **Benchmark methodology differences still exist**
   - Even after unit conversion, benchmark harnesses are different (`go test -bench` vs custom C++ timer).
   - For stricter comparison, align workload profile + metric collection (e.g., throughput + p99 latency) across both languages.

## 🧪 Local Performance Comparison (C++ vs Go)

Run one command from repo root:

```bash
./scripts/compare_cpp_go_perf.sh
```

It will:
1. Run Go benchmarks in `go-iml/cache`
2. Build C++ benchmark target
3. Run C++ benchmark and print ops/s
4. Save raw outputs to `perf_go.txt` and `perf_cpp.txt`

## 📝 Key Learnings (The Migration Journey)

*This section documents the transition from C++ to Go.*

1. **Memory Model:** Moving from `std::shared_ptr` semantics to Go's Garbage Collector.
2. **Concurrency:** Comparing `std::thread` overhead vs. Goroutine context switching cost.
3. **Code Complexity:** Lines of code required to implement the sharded map logic.

---

## License

MIT
