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

以下为最新一次本机实测（**2026-10-20**，Apple M4，macOS/darwin arm64；不同机器会有差异）。

测试方法：
- 两端均使用 256 个 shard、10 秒 TTL 和 1 秒后台清理间隔
- 10 个 worker、300,000 个精确工作项、相同的预生成 key/value 形式，并使用优化构建
- 运行 5 个相互隔离的进程；下表为算术平均值及样本变异系数（CV）
- Set 和 Get 使用每秒操作数；并发 Set+Get 使用每秒 pair 数，一个 pair 表示一次 Set 后紧接一次 Get

| 场景 | C++ 平均吞吐/s (CV) | Go 平均吞吐/s (CV) | C++ / Go |
| :--- | :--- | :--- | :--- |
| Set（操作） | `8,155,290` (`12.35%`) | `6,025,967` (`3.71%`) | `1.35x` |
| Get（操作） | `14,808,860` (`1.25%`) | `8,873,410` (`2.08%`) | `1.67x` |
| 并发 Set+Get（pair） | `15,396,280` (`4.42%`) | `11,752,032` (`1.82%`) | `1.31x` |

先前的并发 `5.35x` 将 C++ 的单次操作数/秒与 Go 的 Set+Get pair 数/秒直接比较。统一负载与单位后，本次结果为 `1.31x`。吞吐表本身不能说明剩余差距的原因。

Go CPU/内存分配剖析、交错 GC 对照实验、实测影响因素及归因限制见[基准分析](docs/benchmark-analysis.md)。

## 🧪 本机一键性能对比（C++ vs Go）

在仓库根目录执行：

```bash
./scripts/compare_cpp_go_perf.sh
```

脚本会：
1. 在隔离进程中运行 Go benchmark（`go-iml/cache`）
2. 构建 C++ benchmark 目标
3. 在隔离进程中运行 C++ benchmark
4. 输出统一口径的平均值、变异和 C++/Go 比率
5. 将原始结果保存到 `perf_go.txt` 和 `perf_cpp.txt`

## 📝 迁移中的关键观察

1. **内存模型：** 从 `std::shared_ptr` 语义迁移到 Go GC 思维。
2. **并发模型：** `std::thread` 开销与 goroutine 调度成本对比。
3. **代码复杂度：** 分片 map 逻辑在两种语言中的实现复杂度差异。

## License

MIT
