# NanoCache: High-Performance In-Memory Cache

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

*Creating a baseline comparison between the two implementations.*

| Metric | C++ Implementation | Go Implementation |
| :--- | :--- | :--- |
| **Throughput (OPS)** | *TBD* | *TBD* |
| **Memory Footprint** | *TBD* | *TBD* |
| **Lock Contention** | *High/Low* | *High/Low* |

> *Benchmarks will be run on [Your CPU Specs] with 100 concurrent workers and 1M keys.*

## 📝 Key Learnings (The Migration Journey)

*This section documents the transition from C++ to Go.*

1.  **Memory Model:** Moving from `std::shared_ptr` semantics to Go's Garbage Collector.
2.  **Concurrency:** Comparing `std::thread` overhead vs. Goroutine context switching cost.
3.  **Code Complexity:** Lines of code required to implement the sharded map logic.

---

## License

MIT

## 🔁 CI/CD

GitHub Actions workflow is provided at `.github/workflows/ci.yml` with:
- Go unit tests + benchmark smoke test
- C++ CMake build + CTest + benchmark smoke test

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
