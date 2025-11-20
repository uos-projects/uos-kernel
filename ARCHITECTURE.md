# 代码架构详解

## 📐 代码组织架构

### 1. 后端服务 (`server/app.py`)

#### 1.1 SparkSession 管理

**延迟初始化模式**:
```python
_spark_session = None

def get_spark() -> SparkSession:
    global _spark_session
    if _spark_session is None:
        _spark_session = build_spark()
    return _spark_session
```

**设计原因**:
- SparkSession 初始化耗时（需要下载 JAR、连接 Nessie 等）
- 延迟加载避免阻塞 FastAPI 启动
- 单例模式确保只有一个 SparkSession 实例

**配置要点**:
- Iceberg Catalog: 使用 NessieCatalog，支持 Git-like 版本管理
- S3 FileIO: 配置 MinIO 作为对象存储
- Spark Extensions: 启用 IcebergSparkSessionExtensions 支持时间旅行语法

#### 1.2 时间旅行查询构建

**`build_clause()` 函数**:
```python
def build_clause(mode: TimeMode, value: str) -> str:
    if mode == TimeMode.biz:
        return f"WHERE '{value}' BETWEEN valid_from AND valid_to"
    if mode == TimeMode.snapshot:
        return f"VERSION AS OF {value}"
    if mode == TimeMode.timestamp:
        return f"TIMESTAMP AS OF '{value}'"
```

**三种模式对比**:

| 模式 | SQL 语法 | 用途 | 示例 |
|------|---------|------|------|
| `biz` | `WHERE 'time' BETWEEN valid_from AND valid_to` | 业务时间查询 | 查询 2025-01-15 的业务状态 |
| `snapshot` | `VERSION AS OF <snapshot_id>` | 技术版本查询 | 查询特定 snapshot 的数据 |
| `timestamp` | `TIMESTAMP AS OF '<timestamp>'` | 技术时间查询 | 查询特定时间点的数据 |

**设计考虑**:
- 业务时间：基于 `valid_from`/`valid_to`，适合业务场景
- 技术时间：基于 Iceberg snapshot，适合数据审计和恢复

#### 1.3 查询执行流程

**`get_view()` 函数流程**:
```
1. 获取 SparkSession（延迟初始化）
2. 构建查询子句（根据 mode）
3. 遍历所有实体表查询
   - 使用 try-except 处理表不存在的情况
   - 为每个实体添加 entity_type 字段
4. 遍历所有关系表查询
   - 同样使用异常处理
5. 返回合并结果
```

**错误处理策略**:
- 表不存在：跳过，不中断整个查询
- 查询失败：记录警告，继续处理其他表
- 保证部分数据可用性

#### 1.4 Snapshots 查询

**设计特点**:
- 支持单表查询和全表查询
- 自动按时间排序
- 包含表名信息，便于前端显示

**实现**:
```python
# 单表查询
if table:
    metadata_table = f"{table}.snapshots"
    df = spark.table(metadata_table)
    
# 全表查询
all_tables = [...]
for tbl in all_tables:
    metadata_table = f"{tbl}.snapshots"
    # 收集所有 snapshots
```

### 2. 前端应用 (`web/app.js`)

#### 2.1 数据获取与渲染

**数据流**:
```
用户操作
  ↓
fetchView(mode, value)
  ↓
GET /api/view
  ↓
renderEntities(data.entities)
renderRelations(data.relations)
renderGraph(entities, relations)
```

**关键函数**:

**`fetchView()`**:
- 发送 API 请求
- 处理错误和加载状态
- 同步时间轴位置（业务时间模式）

**`renderEntities()` / `renderRelations()`**:
- 动态创建表格行
- 应用中文映射
- 显示时间范围信息

#### 2.2 拓扑图渲染 (`renderGraph()`)

**D3.js 力导向图实现**:

**数据结构构建**:
```javascript
// 节点映射
const nodesMap = new Map();
entities.forEach((e) => {
    nodesMap.set(e.entity_id, {
        id: e.entity_id,
        label: e.name || e.entity_id,
        type: e.entity_type,
    });
});

// 边数据
const links = relations.map((r) => ({
    source: r.source_id,
    target: r.target_id,
    relation_type: r.relation_type,
}));
```

**力导向图配置**:
```javascript
simulation = d3.forceSimulation(nodes)
    .force("link", d3.forceLink(links).distance(80))  // 连线距离
    .force("charge", d3.forceManyBody().strength(-150)) // 节点排斥力
    .force("center", d3.forceCenter(width/2, height/2)) // 中心力
    .force("collision", d3.forceCollide().radius(20))  // 碰撞检测
```

**拖拽功能**:
```javascript
const drag = d3.drag()
    .on("start", (event, d) => {
        d.fx = d.x;  // 固定位置
        d.fy = d.y;
    })
    .on("drag", (event, d) => {
        d.fx = event.x;  // 更新位置
        d.fy = event.y;
    })
    .on("end", (event, d) => {
        d.fx = null;  // 释放固定
        d.fy = null;
    });
```

**连线标签**:
```javascript
const linkLabels = linkGroup
    .selectAll("text")
    .data(links)
    .enter()
    .append("text")
    .text((d) => getRelationTypeName(d.relation_type));

// 在 tick 事件中更新位置
linkLabels
    .attr("x", (d) => (d.source.x + d.target.x) / 2)
    .attr("y", (d) => (d.source.y + d.target.y) / 2);
```

#### 2.3 时间轴功能

**时间轴标记点**:
```javascript
function updateTimelineMarkers(snapshots) {
    // 1. 提取唯一时间点（按天）
    const uniqueDates = new Set();
    snapshots.forEach(snap => {
        const dayOffset = Math.floor((snapDate - baseDate) / (1000 * 60 * 60 * 24));
        uniqueDates.add(dayOffset);
    });
    
    // 2. 创建标记点
    Array.from(uniqueDates).forEach(dayOffset => {
        const marker = document.createElement("div");
        marker.style.left = `${(dayOffset / maxDays) * 100}%`;
        // 点击事件：跳转到对应时间
    });
}
```

**双向同步**:
- 时间轴拖动 → 更新输入框 → 自动查询
- 手动输入时间 → 同步时间轴位置

### 3. 数据脚本

#### 3.1 `bootstrap_iceberg.py`

**功能**:
- 创建命名空间（`ontology.grid`, `ontology.relations` 等）
- 创建表结构（定义 schema）
- 写入初始示例数据

**表创建模式**:
```python
spark.sql("""
    CREATE TABLE IF NOT EXISTS ontology.grid.substation (
        entity_id STRING,
        name STRING,
        region STRING,
        nominal_voltage_kv DOUBLE,
        valid_from TIMESTAMP,
        valid_to TIMESTAMP,
        op_type STRING,
        ingestion_ts TIMESTAMP
    )
    USING ICEBERG
    PARTITIONED BY (region)
""")
```

**特点**:
- 使用 `IF NOT EXISTS` 避免重复创建
- 按业务字段分区（如 `region`）
- 包含时间字段用于时间旅行

#### 3.2 `generate_timeline_data.py`

**时间线数据生成**:

**数据结构**:
```python
timeline = [
    {
        "timestamp": datetime(2025, 1, 1),
        "entities": {
            "substation": [...],
            "voltage_level": [...],
        },
        "relations": [...],
    },
    # 更多时间点...
]
```

**写入逻辑**:
```python
def write_data_at_timestamp(spark, data, timestamp):
    # 1. 遍历所有实体类型
    for entity_type, entities in data.get("entities", {}).items():
        # 2. 确定命名空间
        namespace = determine_namespace(entity_type)
        # 3. 构建记录（添加时间字段）
        records = [add_time_fields(e, timestamp) for e in entities]
        # 4. 写入 Iceberg 表
        df.writeTo(table_name).append()
```

**Snapshot 创建**:
- 每次 `append()` 操作创建一个新的 snapshot
- 使用 `time.sleep(2)` 确保 snapshot 时间戳不同
- 支持跨表 snapshot 查询

## 🔍 关键设计模式

### 1. 单例模式
- SparkSession: 全局唯一实例
- 避免资源浪费和连接冲突

### 2. 延迟加载
- SparkSession: 首次使用时初始化
- 提升服务启动速度

### 3. 异常容错
- 表不存在：跳过，继续处理
- 查询失败：记录日志，返回部分结果

### 4. 响应式设计
- SVG viewBox: 自适应容器宽度
- 窗口大小改变：重新计算布局

## 📊 数据模型设计

### SCD Type 2 模式

**字段设计**:
- `valid_from`: 记录生效时间
- `valid_to`: 记录失效时间（9999-12-31 表示当前有效）
- `op_type`: 操作类型（insert/update/delete）
- `ingestion_ts`: 数据摄入时间（用于去重）

**查询模式**:
```sql
-- 查询特定业务时间点的数据
SELECT * FROM table
WHERE '2025-01-15' BETWEEN valid_from AND valid_to
```

### 关系表设计

**通用结构**:
```sql
CREATE TABLE ontology.relations.xxx (
    rel_id STRING,
    source_id STRING,
    target_id STRING,
    relation_type STRING,
    valid_from TIMESTAMP,
    valid_to TIMESTAMP,
    op_type STRING,
    ingestion_ts TIMESTAMP
)
PARTITIONED BY (relation_type)
```

**设计考虑**:
- 按关系类型分区，提升查询性能
- 支持关系的时间旅行
- 统一的关系模型，便于扩展

## 🎨 UI/UX 设计

### 布局结构

```
┌─────────────────────────────────────────┐
│  Header: 标题和描述                     │
├─────────────────────────────────────────┤
│  Controls: 时间模式、时间轴、查询按钮    │
├─────────────────────────────────────────┤
│  ┌──────────────────┐  ┌─────────────┐│
│  │  拓扑视图        │  │  实体列表    ││
│  │  (D3.js)        │  │  关系列表    ││
│  └──────────────────┘  └─────────────┘│
├─────────────────────────────────────────┤
│  Snapshot 列表                          │
└─────────────────────────────────────────┘
```

### 交互特性

1. **时间轴**:
   - 拖动滑块：实时查询
   - 标记点：快速跳转
   - 双向同步：输入框 ↔ 时间轴

2. **拓扑图**:
   - 节点拖拽：自定义布局
   - 悬停效果：显示详细信息
   - 连线标签：显示关系类型

3. **数据表格**:
   - 点击 snapshot：自动填充查询参数
   - 实时更新：查询后自动刷新

## 🔐 安全与性能

### 当前实现

**安全**:
- CORS: 允许所有来源（开发环境）
- 输入验证: 检查 snapshot ID 格式
- 异常处理: 避免信息泄露

**性能**:
- SparkSession 复用: 避免重复初始化
- 延迟加载: 按需创建资源
- 批量查询: 一次请求获取所有数据

### 改进建议

**安全**:
- 添加认证机制（JWT）
- 限制 CORS 来源
- 输入参数验证和清理

**性能**:
- 添加 Redis 缓存
- 查询结果分页
- 异步查询处理

## 🧪 测试策略

### 单元测试

**后端**:
- `build_clause()`: 测试不同 mode 的 SQL 生成
- `get_view()`: Mock SparkSession，测试查询逻辑

**前端**:
- `getEntityTypeName()`: 测试中文映射
- `updateTimelineMarkers()`: 测试标记点生成

### 集成测试

- 端到端测试: 从 API 到前端渲染
- 时间旅行测试: 验证不同时间点的数据正确性
- 性能测试: 大量数据下的查询性能

## 📈 监控与日志

### 当前日志

**后端**:
- FastAPI 自动日志（uvicorn）
- Spark 警告和错误
- 自定义警告（表不存在等）

**前端**:
- Console.log（调试信息）
- 错误提示（用户可见）

### 建议增强

- 结构化日志（JSON 格式）
- 性能指标（查询耗时）
- 错误追踪（Sentry 等）

## 🔄 部署流程

### 开发环境

```bash
# 1. 启动基础设施
make infra-up

# 2. 初始化数据
python scripts/bootstrap_iceberg.py
python scripts/generate_timeline_data.py

# 3. 启动服务
uvicorn server.app:app --host 0.0.0.0 --port 8000
```

### 生产环境建议

1. **容器化**:
   - Dockerfile for FastAPI
   - Docker Compose for 完整栈

2. **配置管理**:
   - 环境变量配置
   - 配置文件外部化

3. **监控**:
   - Prometheus + Grafana
   - 日志聚合（ELK）

4. **高可用**:
   - 负载均衡
   - 数据库备份
   - 灾难恢复计划

