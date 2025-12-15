# 创建 Breaker 表的完整步骤

## 前置条件

1. **安装 PySpark**（已完成）
   ```bash
   source .venv/bin/activate
   pip install pyspark
   ```

2. **启动必要的服务**（Nessie Catalog 和 MinIO）
   ```bash
   cd infra
   docker compose up -d
   # 或
   docker-compose up -d
   ```

   验证服务运行：
   ```bash
   docker compose ps
   # 应该看到：
   # - iceberg-nessie (端口 19120)
   # - iceberg-minio (端口 19000)
   ```

## 执行步骤

### 方法1: 使用 Python 脚本（推荐）

```bash
cd /home/chun/Develop/uos-projects/uos-kernel
source .venv/bin/activate
python3 scripts/execute_breaker_table_sql.py
```

### 方法2: 使用 spark-sql（如果已安装完整 Spark）

```bash
spark-sql -f scripts/create_breaker_table.sql
```

### 方法3: 在 Spark Shell 中手动执行

启动 Spark Shell 并执行 `scripts/create_breaker_table.sql` 中的 SQL 语句。

## 验证

执行以下 SQL 验证表是否创建成功：

```sql
SHOW TABLES IN ontology.grid LIKE 'breaker*';
DESCRIBE EXTENDED ontology.grid.breaker_snapshots;
```

## 当前状态

- ✅ PySpark 已安装在 `.venv` 中
- ✅ SQL 文件已生成：`scripts/create_breaker_table.sql`
- ✅ Python 执行脚本已创建：`scripts/execute_breaker_table_sql.py`
- ⚠️ 需要启动 Nessie 和 MinIO 服务

## 故障排除

如果遇到连接错误：
1. 检查 Nessie 服务是否运行：`curl http://localhost:19120/api/v2/config`
2. 检查 MinIO 服务是否运行：`curl http://localhost:19000`
3. 启动服务：`cd infra && docker compose up -d`
