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

Baseline result (sample from local run; your machine may differ):

| Scenario | C++ | Go |
| :--- | :--- | :--- |
| Set | `~1.22e6 ops/s` | `~1258–1629 ns/op` |
| Get | `~2.78e6 ops/s` | `~74.7–83.4 ns/op` |
| Concurrent Set+Get | `~1.06e7 ops/s` | `~671.6–847.3 ns/op` |

> Note: C++ benchmark currently reports **ops/s**, while Go benchmark reports **ns/op** from `testing.B`. This is suitable for trend comparison, not strict apples-to-apples absolute comparison.

### 简单原因分析（为什么会有差异）

1. **运行时模型差异**
   - C++ 直接运行在原生线程模型上，热点路径更容易接近“零额外运行时开销”。
   - Go 有 goroutine 调度与 GC 参与，吞吐在高并发下很稳，但单次操作会有调度/运行时成本。

2. **内存分配与回收机制差异**
   - C++ 主要由程序控制对象生命周期，路径可预测。
   - Go 在 `Set`/并发场景下会有分配与 GC 压力（从 benchmark 的 allocs/op 也能看出）。

3. **锁与并发访问模式**
   - 两者都用分片 + 读写锁，设计思路一致。
   - 但具体锁实现和调度策略不同，会影响在冲突场景下的延迟与吞吐表现。

4. **基准实现方式不同**
   - Go 用 `go test -bench`，统计维度是 `ns/op`、`B/op`、`allocs/op`。
   - C++ 当前是自定义计时程序，指标以 `ops/s` 为主。若要更严格对齐，建议后续统一到同一负载模型、同一指标单位（例如都换算为吞吐+P99）。

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
