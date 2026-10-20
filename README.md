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

Baseline result from the latest local run on **2026-10-20** (Apple M4, macOS/darwin arm64; your machine may differ).

Method:
- 256 shards, 10-second TTL, and a one-second background-cleanup interval in both implementations
- 10 workers, 300,000 exact work items, identical precomputed key/value shapes, and optimized builds
- Five isolated process runs; results below are arithmetic means with sample coefficient of variation (CV)
- Set and Get use operations/second. Concurrent Set+Get uses pairs/second, where one pair is one Set followed by one Get.

| Scenario | C++ mean/s (CV) | Go mean/s (CV) | C++ / Go |
| :--- | :--- | :--- | :--- |
| Set (operations) | `8,155,290` (`12.35%`) | `6,025,967` (`3.71%`) | `1.35x` |
| Get (operations) | `14,808,860` (`1.25%`) | `8,873,410` (`2.08%`) | `1.67x` |
| Concurrent Set+Get (pairs) | `15,396,280` (`4.42%`) | `11,752,032` (`1.82%`) | `1.31x` |

The earlier `5.35x` concurrent figure compared C++ individual operations/second with Go Set+Get pairs/second. With the workload and units normalized, this run measured `1.31x`. The throughput table alone does not identify the cause of the remaining gap.

See [the benchmark analysis](docs/benchmark-analysis.md) for Go CPU/allocation profiles, an interleaved GC control, measured contributors, and the limits of the attribution.

## 🧪 Local Performance Comparison (C++ vs Go)

Run one command from repo root:

```bash
./scripts/compare_cpp_go_perf.sh
```

It will:
1. Run isolated Go benchmark processes in `go-iml/cache`
2. Build C++ benchmark target
3. Run isolated C++ benchmark processes
4. Print normalized means, variation, and C++/Go ratios
5. Save raw outputs to `perf_go.txt` and `perf_cpp.txt`

## 📝 Key Learnings (The Migration Journey)

*This section documents the transition from C++ to Go.*

1. **Memory Model:** Moving from `std::shared_ptr` semantics to Go's Garbage Collector.
2. **Concurrency:** Comparing `std::thread` overhead vs. Goroutine context switching cost.
3. **Code Complexity:** Lines of code required to implement the sharded map logic.

---

## License

MIT
