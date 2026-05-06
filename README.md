# UOS Kernel — 人机物统一资源操作系统内核

UOS Kernel 将电力系统中的**物理资源、人力资源、信息资源**统一抽象为 Resource，为每个资源内建数字孪生（状态历史），并通过拓扑图管理资源间关系。上层应用和 Agent 通过 ROSIX 接口访问世界模型，实现感知、推理和执行。

## 核心思想

```
世界模型 = 全部资源的数字孪生 + 拓扑关系图
Agent 通过 ROSIX 接口读写世界模型，完成感知-推理-执行闭环
```

### 三类资源统一管理

| 类别 | 实例 | 说明 |
|------|------|------|
| 物理资源 | 变电站、馈线、开关站、FTU/DTU | 电网设备实体 |
| 人力资源 | 运维人员、班组、服务区域 | 组织与人 |
| 信息资源 | 缺陷工单、告警规则、知识条目 | 业务对象 |

### 数字孪生 — 内建的状态历史

每个资源自带不可变的状态变化时间线：

```
Twin(终端设备) = 静态属性 + 当前状态 + [状态事件序列]

StateEvent:
  Timestamp: 2026-02-02 08:12
  Field:     "status"
  OldValue:  "online"
  NewValue:  "offline"
  Cause:     "专网信号差"
  Actor:     "system"
```

当前状态是历史的投影（fold），不独立维护。

### 拓扑 — 资源间关系是一等公民

```
物理拓扑:  变电站 →[contains]→ 馈线 →[contains]→ 开关站 →[contains]→ 终端
组织拓扑:  区域 →[contains]→ 人员
动态绑定:  工单 →[assigned]→ 人员,  工单 →[located_at]→ 终端
```

## ROSIX 接口

POSIX 风格的统一资源访问接口，扩展了时态和拓扑操作：

| 接口 | 语义 |
|------|------|
| `Open(id)` | 打开资源，返回描述符 FD |
| `Close(fd)` | 关闭描述符 |
| `Read(fd)` | 读取资源孪生（当前状态 + 元数据） |
| `Write(fd, field, value, cause, actor)` | 写入状态变更，自动追加到孪生时间线 |
| `RCtl(fd, cmd, args)` | 资源控制命令 |
| `History(fd, since, until)` | 查询状态变化历史 |
| `Watch(fd, filter)` | 订阅状态变化事件流 |
| `Relate(fd1, fd2, rel)` | 建立资源间关系 |
| `Traverse(fd, dir, rel)` | 沿拓扑遍历关联资源 |

## 目录结构

```
uos-kernel/
├── kernel/              # 核心类型和接口定义
│   ├── resource.go      #   Resource, ResourceKind, ResourceID
│   ├── twin.go          #   DigitalTwin, StateEvent, TwinStore
│   ├── topology.go      #   Graph, Edge, RelationType
│   ├── rosix.go         #   ROSIX 接口
│   └── world.go         #   WorldModel 聚合
├── internal/
│   ├── twin/            #   内存版 TwinStore 实现
│   ├── topo/            #   内存版拓扑图实现（BFS 遍历）
│   ├── rosix/           #   ROSIX 完整实现
│   └── importer/        #   CSV 数据导入（消缺工单 → 世界模型）
├── cmd/demo/            # 端到端演示
├── data/                # 数据文件
│   └── defects.csv      #   4643条配电自动化消缺工单
└── docs/                # 设计文档
```

## 快速开始

```bash
# 运行端到端 demo
go run ./cmd/demo/

# 运行测试
go test ./...
```

Demo 会加载 4643 条真实消缺工单，构建包含 6000+ 资源和 15000+ 拓扑关系的世界模型，然后演示：
- 读取工单孪生及其完整状态变迁
- 沿拓扑查询人员负责的所有工单
- 沿拓扑查询变电站下辖的馈线
- 运行时写入新状态，验证孪生时间线增长

## 设计原则

1. **孪生内建** — 任何对资源的 Write 自动记录到时间线，不需要应用层维护
2. **历史不可变** — Timeline 只追加，不修改
3. **拓扑一等公民** — 关系独立于资源属性，支持高效图遍历
4. **接口可插拔** — TwinStore 和 Graph 都是 interface，可替换持久化实现
5. **世界模型即 Agent 基座** — Agent 通过 ROSIX 获取完整世界认知，无需自建模型

## 演进方向

1. **ROSIX as Agent Tools** — 将 ROSIX 暴露为 LLM tool_use 接口，Agent 用自然语言操作世界模型
2. **自治 Agent** — 资源订阅 Watch 事件流，自主执行 Perception-Plan-Act 闭环
3. **多 Agent 协同** — 多个 Agent 通过世界模型共享状态，协商完成跨资源编排
