# 资源类型系统内核

基于操作系统内核对象类型系统概念设计的电力系统资源类型系统，提供POSIX风格的管理接口。

## 设计理念

参考操作系统内核中的对象类型系统（Object Type System），设计了一个元层模型（DSL）来定义类型系统，将CIM模型转换为类型系统，并定义类似POSIX的管理接口。

### 核心概念映射

| 内核概念 | 电力系统映射 | 说明 |
|---------|------------|------|
| Object Type | ResourceType | 定义资源类型的结构、属性和操作 |
| Object Instance | ResourceInstance | 资源类型的实例 |
| Object Manager | ResourceManager | 统一管理所有资源实例 |
| Handle | ResourceHandle | 用户态访问资源的抽象句柄 |
| Type Descriptor | TypeDescriptor | 类型描述符，定义类型元数据 |
| Operations | Capabilities | 资源支持的操作（Control、Measurement等） |

## 架构

```
┌─────────────────────────────────────┐
│   POSIX-like API Layer              │  ← 用户态接口
│   (Open, Read, Write, Ioctl)        │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│   Resource Kernel                   │  ← 内核层
│   - Type Registry                   │
│   - Resource Manager                │
│   - Handle Management               │
│   - Capability Router               │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│   Type System (Meta Layer)          │  ← 元层
│   - Type Descriptors                │
│   - CIM to TypeSystem Converter     │
│   - Operation Validator             │
└─────────────────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│   Actor System (Implementation)     │  ← 实现层
│   - PowerSystemResourceActor        │
│   - Capacity Implementations        │
│   - Message Routing                 │
└─────────────────────────────────────┘
```

## 文件结构

```
kernel/
├── typesystem.yaml          # 类型系统DSL定义
├── types.go                 # 类型描述符数据结构
├── registry.go              # 类型注册表
├── loader.go                # YAML加载器
├── cim_converter.go         # CIM到类型系统转换器
├── kernel.go                # 资源内核实现
├── go.mod                    # Go模块定义
└── cmd/
    └── example/
        └── main.go          # 使用示例
```

## 使用方法

### 1. 定义类型系统

在 `typesystem.yaml` 中定义资源类型：

```yaml
resource_types:
  - name: Breaker
    base_type: PowerSystemResource
    attributes:
      - name: ratedCurrent
        type: float
        required: false
    capabilities:
      - name: SwitchControl
        operations: [open, close, queryState]
```

### 2. 加载类型系统

```go
k := kernel.NewResourceKernel(system)
err := k.LoadTypeSystem("typesystem.yaml")
```

### 3. 使用POSIX风格接口

```go
// 打开资源
fd, err := k.Open("Breaker", "BREAKER_001", 0)
defer k.Close(fd)

// 查询资源信息
stat, err := k.Stat(ctx, fd)

// 读取资源状态
state, err := k.Read(ctx, fd)

// 执行控制操作
result, err := k.Ioctl(ctx, fd, 0x1004, map[string]interface{}{"value": 150.0})
```

## 核心接口

### ResourceKernel

资源内核，提供POSIX风格的系统调用接口：

- `Open(resourceType, resourceID, flags)` - 打开资源
- `Close(fd)` - 关闭资源
- `Read(ctx, fd)` - 读取资源状态
- `Write(ctx, fd, req)` - 写入资源状态
- `Stat(ctx, fd)` - 查询资源信息
- `Ioctl(ctx, fd, request, argp)` - 控制操作
- `Find(resourceType, filter)` - 查找资源
- `Watch(fd, events)` - 监听资源变化

### TypeRegistry

类型注册表，管理所有类型描述符：

- `Register(descriptor)` - 注册类型
- `Get(typeName)` - 获取类型描述符
- `List()` - 列出所有类型
- `Exists(typeName)` - 检查类型是否存在

### TypeDescriptor

类型描述符，定义资源类型的元数据：

- `GetAttribute(name)` - 获取属性描述符
- `GetAllAttributes()` - 获取所有属性（包括继承）
- `HasCapability(name)` - 检查是否具有能力
- `GetCapability(name)` - 获取能力描述符
- `ValidateOperation(capability, operation)` - 验证操作

## CIM模型转换

`CIMToTypeSystemConverter` 可以将CIM模型转换为类型系统：

```go
converter := kernel.NewCIMToTypeSystemConverter(registry)

entityInfo := kernel.EntityInfo{
    Name: "Breaker",
    Inherits: "PowerSystemResource",
    Attributes: [...],
    Relationships: [...],
}

descriptor, err := converter.ConvertEntityToType("Breaker", entityInfo)
```

转换器会：
1. 从CIM实体提取属性
2. 从CIM关系推断能力（Control、Measurement等）
3. 生成类型描述符

## 操作映射

在 `typesystem.yaml` 中定义ioctl命令到能力操作的映射：

```yaml
operation_mapping:
  ioctl_commands:
    - ioctl_cmd: 0x1004
      capability: Control
      operation: execute
      message_type: SetPointMessage
```

## 示例

参见 `cmd/example/main.go` 了解完整的使用示例。

## 与Actor系统集成

类型系统内核基于现有的Actor系统实现：

1. **ResourceManager** - 管理Actor实例的句柄
2. **类型验证** - 在执行操作前验证资源类型和能力
3. **操作路由** - 将ioctl命令路由到对应的Capacity

## 未来扩展

- [ ] 实现Watch接口（资源变化监听）
- [ ] 完善Find接口（按类型和属性查找）
- [ ] 添加访问控制（基于角色的权限管理）
- [ ] 支持资源生命周期管理
- [ ] 实现资源快照和恢复

