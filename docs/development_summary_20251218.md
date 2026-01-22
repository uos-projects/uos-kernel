# Development Summary - 2025-12-18

## 核心架构演通过程：从 POSIX 到 分布式持久化对象

本次开发迭代主要集中在完善 `ResourceKernel` 的语义完整性，使其不仅具备 POSIX 的外壳，更拥有符合电力业务需求的强类型和持久化内核。

### 1. 资源创建与类型安全 (Open-or-Create)

#### 1.1 `O_CREAT` 语义实现
实现了 POSIX 标准的 `O_CREAT` 标志支持，确立了“获取即存在”的原子语义。

*   **自动创建机制**：
    *   当 `Open` 调用指定 `O_CREAT` 且资源不存在时，内核会自动根据 `resourceType` 创建 `CIMResourceActor`。
    *   **智能推导**：内核从类型系统 (`TypeRegistry`) 推导 `OWLClassURI`，并尝试自动注入默认的能力（如检测到 `Control` 能力组时自动添加 `CommandCapacity`）。
*   **示例**：
    ```go
    // 如果 BREAKER_001 不存在则自动创建，存在则直接打开
    fd, err := k.Open("Breaker", "BREAKER_001", resource.O_CREAT)
    ```

#### 1.2 运行时类型安全增强
在 POSIX 原有的文件描述符基础上，增加了面向对象的强类型检查。
*   **类型校验**：`Open` 一个已存在的资源时，内核会强制检查请求的 `resourceType` 与实际运行时的 `Actor.ResourceType` 是否匹配。
*   **异常处理**：类型不匹配时返回 `resource type mismatch` 错误，防止了错误的类型转换和操作。

### 2. 对象属性的写操作 (Write Semantics)

将 POSIX 的 `write(fd, buf)` 语义在对象系统中重定义为 **“属性更新 (Property Update)”**。

*   **WriteRequest**：
    *   不传输原始字节流，而是传输 `Updates map[string]interface{}`。
*   **Actor 集成**：
    *   利用 `actors.PropertyHolder` 接口，使 `ResourceManager` 能够通用的修改 `CIMResourceActor` 的属性。
    *   **全链路验证**：实现了从 `k.Write` 修改属性 -> `k.Read` 验证修改结果的完整闭环。

### 3. 持久化架构：内存+快照模式 (Iceberg Integration)

确立了 **"Memory for Runtime, Iceberg for History"** 的分层存储架构，以解决高频状态变更与低频大规模分析之间的矛盾。

#### 3.1 架构设计
*   **热数据 (Hot)**：Actor 内存状态，毫秒级响应 `Read`/`Write`/`Ioctl`。
*   **冷数据/温数据 (Cold/Warm)**：通过快照 (`Snapshot`) 形式持久化到后端存储（如 Apache Iceberg），提供时间旅行 (`Time Travel`) 和审计能力。

#### 3.2 核心抽象：`SnapshotStore`
定义了存储无关的快照接口 (`actor/state/snapshot_store.go`)：
```go
type SnapshotStore interface {
    Save(ctx context.Context, resourceID string, state interface{}) error
    LoadAt(ctx context.Context, resourceID string, timestamp int64) (interface{}, error)
}
```
*   **MemorySnapshotStore**：实现了基于内存的 Mock 版本，用于单元测试和原型验证。

#### 3.3 系统控制命令：`CMD_SYNC`
为了支持显式的持久化控制，引入了系统级 IOCTL 命令。
*   **命令定义**：`CMD_SYNC (0x2001)`。
*   **通道穿透**：`ResourceKernel.Ioctl` 特殊处理该命令，绕过常规的业务 Capability 检查，直接调用 `CIMResourceActor.SaveSnapshotAsync`。
*   **使用方式**：
    ```go
    k.Ioctl(ctx, fd, resource.CMD_SYNC, nil) // 强制触发快照保存
    ```

## 总结
通过本次迭代，`uos-kernel` 已经具备了一个微内核操作系统的雏形：
1.  **文件系统层**：完善的 `Open` (Create/TypeCheck)。
2.  **数据层**：对象化的 `Read`/`Write`。
3.  **驱动层**：基于 Capability 的 `Ioctl`。
4.  **存储层**：基于快照的 `SnapshotStore`。
