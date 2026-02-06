
# 项目代号：NanoCache (C++ Prototype)

## 1. 项目目标

构建一个高性能、线程安全、支持过期时间的内存键值存储库（In-Memory Key-Value Store）。
**核心指标：** 在高并发（100+ 线程）读写混合场景下，吞吐量（OPS）显著优于单一大锁（`std::mutex`）的实现。

## 2. 核心接口规范 (API)

你的 C++ 类 `ShardedCache` 需要对外暴露以下公共接口。为了简化，我们暂时统一使用 `std::string` 作为 Key 和 Value。

```cpp
class ShardedCache {
public:
    // 构造函数：初始化分片数（默认 256）
    ShardedCache(size_t shard_count = 256);

    // 设置键值对
    // ttl_ms: 过期时间（毫秒）。如果为 0，则永不过期（可选，或视为立即过期，建议设为永不过期）
    void Set(const std::string& key, const std::string& value, int64_t ttl_ms);

    // 获取值
    // 返回：如果存在且未过期，返回 value；否则返回空（可以用 std::optional, 指针, 或空字符串+bool）
    // 建议：使用 std::optional<std::string> (C++17) 或 bool Get(key, &value)
    std::optional<std::string> Get(const std::string& key);

    // (可选) 手动删除
    void Del(const std::string& key);
};

```

## 3. 详细技术约束

### 3.1 分片架构 (Sharding Strategy)

* **分片数量：** 固定为 256 (或者 `2^n`)。
* **路由算法：** 必须使用哈希算法将 Key 映射到特定的分片索引。
* 公式：`shard_index = FNV64(key) % shard_count`
* *注：FNV-1a 是一个简单且分布均匀的哈希算法，比 std::hash 更适合这种场景，或者直接用 `std::hash` 也可以，重点是取模逻辑。*


* **锁隔离：** 每个分片必须拥有**独立**的读写锁（`std::shared_mutex`）。严禁使用全局锁。

### 3.2 存储单元 (Storage Unit)

每个分片内部维护一个 `std::unordered_map`。
Map 的 Value 必须包含：

* 实际数据 (`value`)
* 过期时间戳 (`expiry_timestamp`)

### 3.3 过期策略 (Eviction Policy)

必须实现两种过期机制：

1. **惰性删除 (Lazy Deletion)：**
* 当调用 `Get(key)` 时，先检查当前时间是否超过 `expiry_timestamp`。
* 如果已过期，视为 Key 不存在（返回空），并**顺手**从 map 中删除该记录（释放内存）。


2. **主动清理 (Active Expiration / Janitor)：**
* （C++ 版可选，但强烈建议实现）启动一个后台线程。
* 每隔一段时间（例如 1 秒），遍历所有分片的 map。
* 删除已过期的 Key。
* *挑战点：* 遍历时如何加锁？全锁会导致性能抖动，试着每个分片锁一下，清理完解锁，再处理下一个分片。



## 4. 验证与测试 (Benchmark)

写完类之后，你需要一个 `main.cpp` 来证明你的代码是“高并发安全”的。

**测试场景：**

* **数据量：** 100,000 个不同的 Key。
* **并发度：** 10 个线程同时运行。
* **混合操作：** 每个线程循环执行 10,000 次操作。
* 90% 概率调用 `Get`。
* 10% 概率调用 `Set`。


* **结果观察：** 程序不崩溃（Crash），且最终内存占用稳定（得益于过期清理）。

---

### 建议的开发流

1. **Level 1:** 实现 `Set` 和 `Get`，跑通单线程逻辑。
2. **Level 2:** 引入 `std::shared_mutex` 和分片逻辑，跑通多线程逻辑。
3. **Level 3:** 实现 `main` 函数里的 Benchmark，测试是否有 Race Condition (可以用 ThreadSanitizer: `g++ -fsanitize=thread ...`)。