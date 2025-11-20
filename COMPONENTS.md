# 核心组件详解：职责与协作关系

## 📦 组件概览

本系统使用四个核心组件构建数据湖架构：

1. **SparkSession with Iceberg Extensions** - 计算引擎
2. **Nessie Catalog** - 元数据目录（版本管理）
3. **Apache Iceberg** - 表格式（数据组织）
4. **MinIO** - 对象存储（数据存储）

## 🔍 各组件详细职责

### 1. SparkSession with Iceberg Extensions

#### 职责

**SparkSession**:
- **SQL 查询引擎**: 解析和执行 SQL 查询
- **分布式计算**: 将查询任务分发到多个节点（本项目中为单机模式）
- **数据转换**: DataFrame API 和 SQL 之间的转换
- **资源管理**: 管理计算资源和内存

**Iceberg Extensions** (`IcebergSparkSessionExtensions`):
- **扩展 SQL 语法**: 支持 Iceberg 特有的 SQL 语法
  ```sql
  -- 时间旅行查询
  SELECT * FROM table VERSION AS OF 123456
  SELECT * FROM table TIMESTAMP AS OF '2025-01-01'
  
  -- 表管理
  CALL system.snapshots('table')
  ```
- **优化器集成**: 将 Iceberg 的优化规则集成到 Spark 优化器中
- **元数据访问**: 提供访问 Iceberg 元数据表的接口（如 `.snapshots`）

#### 配置示例

```python
.config("spark.sql.catalog.ontology", "org.apache.iceberg.spark.SparkCatalog")
.config("spark.sql.extensions", "org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions")
```

**作用**:
- `SparkCatalog`: 告诉 Spark 使用 Iceberg Catalog 接口
- `IcebergSparkSessionExtensions`: 启用 Iceberg SQL 扩展语法

### 2. Nessie Catalog

#### 职责

**元数据目录服务**:
- **表注册**: 存储表的元数据（schema、分区信息、位置等）
- **版本管理**: Git-like 的版本控制系统
  - 支持分支（branch）
  - 支持标签（tag）
  - 支持提交历史（commit history）
- **引用管理**: 管理表的引用（当前使用哪个版本）

**API 服务**:
- RESTful API: 提供 HTTP 接口供 Spark/Iceberg 访问
- 元数据查询: 查询表列表、schema、snapshot 信息等

#### 配置示例

```python
.config("spark.sql.catalog.ontology.catalog-impl", "org.apache.iceberg.nessie.NessieCatalog")
.config("spark.sql.catalog.ontology.uri", "http://localhost:19120/api/v2")
.config("spark.sql.catalog.ontology.ref", "main")
```

**作用**:
- `catalog-impl`: 指定使用 NessieCatalog 实现
- `uri`: Nessie 服务的地址
- `ref`: 使用哪个分支（类似 Git 的 branch）

#### 存储内容

Nessie 存储的是**元数据**，不是实际数据：
- 表名和命名空间
- Schema 定义
- 分区规范
- 表的当前 snapshot ID
- 表的存储位置（指向 MinIO）

**不存储**:
- 实际数据文件（存储在 MinIO）
- 数据文件内容

#### ⚠️ 当前使用情况

**重要说明**: 在当前系统中，**Nessie 主要作为元数据存储使用**，**并未充分利用其版本管理功能**。

**当前实际使用**:
- ✅ **表注册**: 存储表的元数据（表名、schema、当前 snapshot ID）
- ✅ **元数据查询**: 通过 Nessie API 查询表信息
- ❌ **分支管理**: 只使用固定的 `main` 分支，没有创建/切换分支
- ❌ **标签管理**: 没有使用标签功能
- ❌ **合并操作**: 没有使用分支合并功能
- ❌ **提交历史**: 没有查询或管理提交历史

**版本管理的实际实现**:
- 真正的版本管理（snapshots）由 **Iceberg 自己管理**
- Iceberg 的 snapshot 机制提供了数据的时间旅行功能
- Nessie 只存储"表名 → 当前 snapshot ID"的映射

**如果要使用 Nessie 的版本管理功能**，需要：
1. 创建多个分支（如 `dev`、`prod`）
2. 在不同分支上进行表操作
3. 使用分支合并功能
4. 使用标签标记重要版本

**当前架构的优势**:
- 简化了系统复杂度
- Iceberg snapshot 已满足时间旅行需求
- 未来可以轻松启用 Nessie 版本管理功能

### 3. Apache Iceberg

#### 职责

**表格式（Table Format）**:
- **数据组织**: 定义数据文件如何组织和存储
- **元数据文件**: 管理 manifest 文件（数据文件清单）
- **Schema 演进**: 支持添加/删除列而不重写数据
- **分区演进**: 支持改变分区方式

**ACID 事务**:
- **原子性**: 写入操作要么全部成功，要么全部失败
- **一致性**: 保证数据一致性
- **隔离性**: 并发写入的隔离
- **持久性**: 写入的数据持久保存

**时间旅行（Time Travel）**:
- **Snapshot 管理**: 每次写入创建一个 snapshot
- **版本查询**: 支持查询历史版本的数据
- **元数据表**: 提供 `.snapshots`、`.files` 等元数据表

#### 文件结构

```
warehouse/
└── ontology/
    └── grid/
        └── substation/
            ├── metadata/
            │   ├── v1.metadata.json      # 表元数据（schema、分区等）
            │   ├── v2.metadata.json
            │   └── ...
            ├── data/
            │   ├── 00000-0-xxx.parquet  # 实际数据文件
            │   ├── 00001-0-xxx.parquet
            │   └── ...
            └── metadata/
                └── snap-xxx-xxx.avro    # Snapshot 元数据
```

**关键文件**:
- `metadata.json`: 表的 schema、分区规范、当前 snapshot
- `manifest-list.avro`: 列出所有 manifest 文件
- `manifest.avro`: 列出数据文件及其统计信息
- `data/*.parquet`: 实际数据文件

#### 配置示例

```python
.config("spark.sql.catalog.ontology.warehouse", "s3a://iceberg/warehouse")
.config("spark.sql.catalog.ontology.io-impl", "org.apache.iceberg.aws.s3.S3FileIO")
```

**作用**:
- `warehouse`: Iceberg 数据仓库的根路径（在 MinIO 中）
- `io-impl`: 使用 S3FileIO 来读写文件（兼容 MinIO）

### 4. MinIO

#### 职责

**对象存储服务**:
- **数据文件存储**: 存储 Iceberg 的数据文件（Parquet 格式）
- **元数据文件存储**: 存储 Iceberg 的元数据文件（JSON、Avro）
- **S3 兼容 API**: 提供与 AWS S3 兼容的 API
- **数据持久化**: 保证数据不丢失

**存储内容**:
- Parquet 数据文件
- Iceberg 元数据文件（metadata.json、manifest 等）
- Nessie 的版本存储（如果配置）

#### 配置示例

```python
.config("spark.sql.catalog.ontology.s3.endpoint", "http://localhost:19000")
.config("spark.sql.catalog.ontology.s3.access-key-id", "iceberg")
.config("spark.sql.catalog.ontology.s3.secret-access-key", "iceberg_password")
.config("spark.hadoop.fs.s3a.endpoint", "http://localhost:19000")
```

**作用**:
- `s3.endpoint`: MinIO 的 API 端点
- `access-key-id` / `secret-access-key`: 认证信息
- `fs.s3a.endpoint`: Hadoop 文件系统的 S3 端点

## 🔄 组件协作流程

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    SparkSession                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Iceberg Extensions                                   │  │
│  │  - 解析 VERSION AS OF 语法                           │  │
│  │  - 优化查询计划                                       │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────┬────────────────────────────────────────┘
                     │ SQL 查询
                     │
        ┌────────────┴────────────┐
        │                         │
┌───────▼──────────┐    ┌─────────▼──────────┐
│  Nessie Catalog  │    │  Apache Iceberg    │
│  (元数据目录)     │    │  (表格式)          │
│                  │    │                    │
│  1. 查询表信息   │───▶│  2. 读取元数据文件 │
│  2. 获取当前     │    │  3. 确定数据文件   │
│     snapshot ID  │    │  4. 构建查询计划   │
└──────────────────┘    └─────────┬──────────┘
                                   │
                                   │ 读取文件路径
                                   │
                          ┌────────▼────────┐
                          │     MinIO       │
                          │  (对象存储)     │
                          │                 │
                          │  读取 Parquet   │
                          │  数据文件       │
                          └─────────────────┘
```

### 详细协作流程

#### 场景 1: 创建表

```
1. SparkSession 执行 CREATE TABLE 语句
   ↓
2. Iceberg Extensions 解析 SQL，识别 USING ICEBERG
   ↓
3. SparkCatalog 调用 NessieCatalog.createTable()
   ↓
4. NessieCatalog:
   - 生成表元数据（schema、分区等）
   - 创建初始 snapshot
   - 将元数据写入 Nessie（表名、schema、初始 snapshot ID）
   ↓
5. Iceberg:
   - 在 MinIO 中创建表目录结构
   - 写入初始 metadata.json 文件
   - 创建空的 manifest 文件
   ↓
6. 返回表对象给 SparkSession
```

**关键点**:
- Nessie 存储：表名、schema、当前 snapshot ID
- MinIO 存储：metadata.json、manifest 文件
- 此时还没有数据文件

#### 场景 2: 写入数据

```
1. SparkSession 执行 INSERT INTO 或 writeTo().append()
   ↓
2. Spark 将数据写入临时 Parquet 文件（本地）
   ↓
3. Iceberg:
   - 创建新的 manifest 文件（列出新数据文件）
   - 创建新的 manifest-list 文件
   - 创建新的 snapshot 元数据
   - 更新 metadata.json（指向新 snapshot）
   ↓
4. 将 Parquet 文件上传到 MinIO
   - 路径: s3a://iceberg/warehouse/ontology/grid/substation/data/xxx.parquet
   ↓
5. 将 manifest 和 metadata 文件上传到 MinIO
   ↓
6. NessieCatalog:
   - 更新 Nessie 中的表元数据
   - 更新当前 snapshot ID
   - 创建新的 commit（版本记录）
   ↓
7. 事务提交完成
```

**关键点**:
- 数据文件存储在 MinIO
- 元数据文件存储在 MinIO
- Nessie 只存储指向最新 snapshot 的引用

#### 场景 3: 时间旅行查询

```
1. SparkSession 执行: SELECT * FROM table VERSION AS OF 123456
   ↓
2. Iceberg Extensions 解析 VERSION AS OF 语法
   ↓
3. SparkCatalog 调用 NessieCatalog.loadTable()
   ↓
4. NessieCatalog:
   - 查询 Nessie 获取表信息
   - 返回表名和当前引用
   ↓
5. Iceberg:
   - 从 MinIO 读取 metadata.json
   - 根据 snapshot ID 查找对应的 snapshot 元数据
   - 读取 manifest-list 文件
   - 读取 manifest 文件，获取数据文件列表
   ↓
6. Spark:
   - 从 MinIO 读取对应的 Parquet 文件
   - 执行查询（过滤、聚合等）
   ↓
7. 返回结果给 SparkSession
```

**关键点**:
- Nessie: 提供表的基本信息
- Iceberg: 管理 snapshot 到数据文件的映射
- MinIO: 存储和读取实际数据

#### 场景 4: 查询当前数据

```
1. SparkSession 执行: SELECT * FROM table
   ↓
2. NessieCatalog.loadTable():
   - 查询 Nessie 获取表信息
   - 获取当前 snapshot ID（latest）
   ↓
3. Iceberg:
   - 从 MinIO 读取最新的 metadata.json
   - 获取当前 snapshot 的信息
   - 读取 manifest-list 和 manifest
   - 确定需要读取的数据文件
   ↓
4. Spark:
   - 从 MinIO 并行读取 Parquet 文件
   - 执行查询
   ↓
5. 返回结果
```

## 📊 数据流向图

### 写入流程

```
应用代码
  ↓
SparkSession.writeTo().append()
  ↓
┌─────────────────────────────────┐
│ Iceberg Writer                  │
│ 1. 创建 Parquet 文件（本地）    │
│ 2. 生成 manifest                │
│ 3. 创建新 snapshot              │
└─────────────────────────────────┘
  ↓                    ↓
MinIO              Nessie Catalog
(上传文件)         (更新元数据引用)
  ↓                    ↓
s3a://iceberg/     表名 → snapshot ID
  warehouse/       分支: main
  ontology/       提交历史
  grid/
  substation/
  data/
  metadata/
```

### 读取流程

```
应用代码
  ↓
SparkSession.sql("SELECT ...")
  ↓
┌─────────────────────────────────┐
│ Nessie Catalog                 │
│ 查询: 表名 → 当前 snapshot ID  │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Iceberg Reader                 │
│ 1. 读取 metadata.json          │
│ 2. 查找 snapshot 元数据        │
│ 3. 读取 manifest-list          │
│ 4. 读取 manifest               │
│ 5. 获取数据文件列表            │
└─────────────────────────────────┘
  ↓
┌─────────────────────────────────┐
│ Spark File Reader              │
│ 从 MinIO 读取 Parquet 文件     │
└─────────────────────────────────┘
  ↓
返回 DataFrame
```

## 🔗 组件间的接口

### 1. Spark ↔ Nessie Catalog

**接口**: Iceberg Catalog API

```python
# Spark 调用
catalog.loadTable(TableIdentifier.of("ontology", "grid", "substation"))

# Nessie 响应
TableMetadata {
    name: "ontology.grid.substation",
    currentSnapshotId: 123456,
    schema: {...},
    ...
}
```

**通信方式**: HTTP REST API
- Spark → Nessie: `GET /api/v2/trees/main?fetch=MINIMAL`
- Nessie → Spark: JSON 响应

### 2. Spark ↔ Iceberg

**接口**: Iceberg Table API

```python
# Spark 通过 Iceberg 读取表
table = catalog.loadTable(...)
table.newScan()
    .useSnapshot(snapshotId)
    .planFiles()
```

**通信方式**: 
- Spark 直接调用 Iceberg Java API
- Iceberg 返回文件列表和统计信息

### 3. Iceberg ↔ MinIO

**接口**: S3 FileIO API

```python
# Iceberg 读取文件
fileIO.newInputFile("s3a://iceberg/warehouse/.../metadata.json")

# MinIO 响应
返回文件内容流
```

**通信方式**: S3 REST API
- Iceberg → MinIO: `GET /bucket/path/to/file`
- MinIO → Iceberg: 文件内容流

### 4. Spark ↔ MinIO

**接口**: Hadoop S3A FileSystem

```python
# Spark 读取数据文件
spark.read.parquet("s3a://iceberg/warehouse/.../data/*.parquet")
```

**通信方式**: S3 REST API
- Spark → MinIO: `GET /bucket/path/to/data.parquet`
- MinIO → Spark: Parquet 文件流

## 🎯 关键配置解析

### 完整配置链

```python
# 1. 指定使用 Iceberg Catalog
.config("spark.sql.catalog.ontology", "org.apache.iceberg.spark.SparkCatalog")
#    ↑
#    告诉 Spark: 使用 Iceberg 的 Catalog 实现

# 2. 指定 Catalog 的具体实现为 Nessie
.config("spark.sql.catalog.ontology.catalog-impl", 
        "org.apache.iceberg.nessie.NessieCatalog")
#    ↑
#    告诉 Iceberg: 使用 Nessie 作为元数据存储

# 3. 配置 Nessie 服务地址
.config("spark.sql.catalog.ontology.uri", "http://localhost:19120/api/v2")
#    ↑
#    NessieCatalog 通过这个地址访问 Nessie API

# 4. 配置数据仓库路径（在 MinIO 中）
.config("spark.sql.catalog.ontology.warehouse", "s3a://iceberg/warehouse")
#    ↑
#    Iceberg 在这个路径下存储数据文件和元数据文件

# 5. 配置文件系统为 S3（兼容 MinIO）
.config("spark.sql.catalog.ontology.io-impl", 
        "org.apache.iceberg.aws.s3.S3FileIO")
#    ↑
#    Iceberg 使用 S3FileIO 来读写文件

# 6. 配置 MinIO 连接信息
.config("spark.sql.catalog.ontology.s3.endpoint", "http://localhost:19000")
.config("spark.sql.catalog.ontology.s3.access-key-id", "iceberg")
.config("spark.sql.catalog.ontology.s3.secret-access-key", "iceberg_password")
#    ↑
#    S3FileIO 使用这些配置连接 MinIO

# 7. 启用 Iceberg SQL 扩展
.config("spark.sql.extensions", 
        "org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions")
#    ↑
#    启用 VERSION AS OF、TIMESTAMP AS OF 等语法
```

## 💡 设计优势

### 1. 职责分离

- **Spark**: 专注于计算和查询执行
- **Nessie**: 专注于元数据管理和版本控制
- **Iceberg**: 专注于数据组织和 ACID 保证
- **MinIO**: 专注于数据持久化存储

### 2. 可扩展性

- **替换 Nessie**: 可以替换为 Hive Metastore、AWS Glue 等
- **替换 MinIO**: 可以替换为 AWS S3、Azure Blob 等
- **替换 Spark**: Iceberg 支持 Flink、Trino 等引擎

### 3. 版本管理

- **Nessie**: Git-like 的版本管理，支持分支和合并
- **Iceberg**: Snapshot 级别的版本管理，支持时间旅行

### 4. 性能优化

- **列式存储**: Parquet 格式，列式读取
- **分区剪枝**: 根据分区信息跳过不相关的文件
- **统计信息**: Manifest 中包含统计信息，优化查询

## 🔍 实际查询示例

### 示例 1: 查询当前数据

```sql
SELECT * FROM ontology.grid.substation
```

**执行流程**:
1. Spark 解析 SQL
2. 调用 `NessieCatalog.loadTable("ontology.grid.substation")`
3. Nessie 返回：当前 snapshot ID = 123456
4. Iceberg 从 MinIO 读取 `metadata.json`
5. 找到 snapshot 123456 的 manifest-list
6. 读取 manifest，获取数据文件列表
7. Spark 从 MinIO 读取 Parquet 文件
8. 返回结果

### 示例 2: 时间旅行查询

```sql
SELECT * FROM ontology.grid.substation VERSION AS OF 123456
```

**执行流程**:
1. Spark + Iceberg Extensions 解析 `VERSION AS OF`
2. 调用 `NessieCatalog.loadTable()`（获取表信息）
3. Iceberg 从 MinIO 读取 `metadata.json`
4. 查找 snapshot 123456 的元数据
5. 读取对应的 manifest-list 和 manifest
6. Spark 从 MinIO 读取历史版本的 Parquet 文件
7. 返回历史数据

### 示例 3: 业务时间查询

```sql
SELECT * FROM ontology.grid.substation 
WHERE '2025-01-15' BETWEEN valid_from AND valid_to
```

**执行流程**:
1. Spark 解析 SQL（普通 WHERE 子句）
2. 调用 `NessieCatalog.loadTable()`（获取当前数据）
3. Iceberg 读取当前 snapshot 的数据文件
4. Spark 读取 Parquet 文件并应用 WHERE 过滤
5. 返回符合业务时间条件的数据

## 📝 总结

### 组件职责矩阵

| 组件 | 存储内容 | 主要职责 | 访问方式 |
|------|---------|---------|---------|
| **SparkSession** | 无 | SQL 解析、查询执行、分布式计算 | 应用代码调用 |
| **Nessie Catalog** | 表元数据引用 | 表注册、版本管理、分支管理 | HTTP REST API |
| **Apache Iceberg** | 元数据文件 | 表格式、ACID、时间旅行、文件组织 | Java API |
| **MinIO** | 数据文件 + 元数据文件 | 对象存储、数据持久化 | S3 REST API |

### 协作关系

1. **Spark** 是入口，接收 SQL 查询
2. **Nessie** 提供表的元数据引用（当前 snapshot）
3. **Iceberg** 管理 snapshot 到文件的映射
4. **MinIO** 存储所有实际文件（数据和元数据）

### 关键理解

- **Nessie 不存储数据文件**，只存储元数据引用
- **Iceberg 不存储数据**，只管理数据文件的组织方式
- **MinIO 存储一切文件**，但不理解表结构
- **Spark 协调一切**，执行查询并返回结果

这种设计实现了**计算与存储分离**、**元数据与数据分离**，提供了强大的扩展性和灵活性。

