# NanoCache：高性能内存缓存

[English](./README.md) | 简体中文

![Language](https://img.shields.io/badge/language-C%2B%2B17%20%7C%20Go-blue)
![License](https://img.shields.io/badge/license-MIT-green)

NanoCache 是一个分片（sharded）、线程安全的内存键值缓存。

本项目是 **Modern C++（C++17）与 Go 的对比实践**，使用相同的架构设计来观察差异：
- 并发模型（OS 线程 vs. Goroutine）
- 内存管理（RAII/智能指针 vs. GC）
- 锁策略与竞争特征

## 🏗 架构

为了降低高并发场景下的锁竞争，NanoCache 使用了 **分片策略**（类似 BigCache/FreeCache）。

- **分片：** 使用 FNV-1a hash 将 key 空间划分为 256 个 shard。
- **加锁：** 每个 shard 独立使用 `RWMutex`（Go）或 `std::shared_mutex`（C++）。
- **过期：** 支持 TTL，包含惰性删除与后台清理（Janitor）。

## 📂 目录结构

- **[`cpp-impl/`](./cpp-impl)**：C++17 实现，偏向手动内存管理与原生并发控制。
- **[`go-iml/`](./go-iml)**：Go 实现，偏向 Goroutine 与运行时调度模型。


## 🚀 快速开始

### C++ 实现

要求：CMake >= 3.10，支持 C++17 的编译器（GCC/Clang/MSVC）。

```bash
cd cpp-impl
mkdir build && cd build
cmake ..
make
./nano_cache_cpp
```

### Go 实现

要求：Go >= 1.20。

```bash
go mod tidy
go run cmd/server/main.go
```

## 📊 性能基线（示例）

以下为一次本机运行样例（不同机器会有差异）。

| 场景 | C++ (ops/s) | Go (ops/s) |
| :--- | :--- | :--- |
| Set | `~1.22e6` | `~0.71e6` |
| Get | `~2.78e6` | `~12.7e6` |
| 并发 Set+Get | `~1.06e7` | `~1.34e6` |

### 性能差异的简单分析

1. **运行时模型不同**
   - C++ 更接近原生线程模型，热点路径额外运行时开销较小。
   - Go 有 goroutine 调度与 GC 机制，开发并发更友好，但单操作存在运行时成本。

2. **内存分配与回收机制不同**
   - C++ 生命周期控制更显式。
   - Go 在 `Set`/并发路径可能出现更多分配与 GC 压力（可从 `allocs/op` 观察）。

3. **锁实现与竞争行为不同**
   - 两端都采用分片 + 读写锁。
   - 但锁实现细节与调度策略不同，会影响冲突时延迟和吞吐。

4. **即使统一单位，基准方法仍有差异**
   - Go 使用 `go test -bench`。
   - C++ 当前是自定义计时程序。
   - 若需更严格横向比较，建议统一负载模型与指标（如吞吐 + p99）。

## 🧪 本机一键性能对比（C++ vs Go）

在仓库根目录执行：

```bash
./scripts/compare_cpp_go_perf.sh
```

脚本会：
1. 跑 Go benchmark（`go-iml/cache`）
2. 构建 C++ benchmark 目标
3. 运行 C++ benchmark 并输出 ops/s
4. 将原始结果保存到 `perf_go.txt` 和 `perf_cpp.txt`

## 📝 迁移中的关键观察

1. **内存模型：** 从 `std::shared_ptr` 语义迁移到 Go GC 思维。
2. **并发模型：** `std::thread` 开销与 goroutine 调度成本对比。
3. **代码复杂度：** 分片 map 逻辑在两种语言中的实现复杂度差异。

## License

MIT
