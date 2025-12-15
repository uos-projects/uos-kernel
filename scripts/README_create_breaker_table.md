# 在 Iceberg 中创建 Breaker 表

## 方法1: 使用 spark-sql 命令行工具（推荐）

如果已安装 Spark，可以直接使用 `spark-sql` 命令：

```bash
spark-sql -f scripts/create_breaker_table.sql
```

或者使用 Spark SQL shell：

```bash
spark-sql \
  --conf spark.sql.catalog.ontology=org.apache.iceberg.spark.SparkCatalog \
  --conf spark.sql.catalog.ontology.catalog-impl=org.apache.iceberg.nessie.NessieCatalog \
  --conf spark.sql.catalog.ontology.uri=http://localhost:19120/api/v2 \
  --conf spark.sql.catalog.ontology.ref=main \
  --conf spark.sql.catalog.ontology.warehouse=s3a://iceberg/warehouse \
  --conf spark.sql.catalog.ontology.io-impl=org.apache.iceberg.aws.s3.S3FileIO \
  --conf spark.sql.catalog.ontology.s3.endpoint=http://localhost:19000 \
  --conf spark.sql.catalog.ontology.s3.access-key-id=iceberg \
  --conf spark.sql.catalog.ontology.s3.secret-access-key=iceberg_password \
  --conf spark.sql.catalog.ontology.s3.path-style-access=true \
  --conf spark.sql.extensions=org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions \
  --packages org.apache.iceberg:iceberg-spark-runtime-3.5_2.12:1.5.0,org.apache.iceberg:iceberg-aws-bundle:1.5.0 \
  -f scripts/create_breaker_table.sql
```

## 方法2: 使用 PySpark（Python）

### 安装依赖

```bash
pip install pyspark
```

### 执行脚本

```bash
python3 scripts/execute_breaker_table_sql.py
```

## 方法3: 在 Spark Shell 中手动执行

启动 Spark Shell：

```bash
spark-shell \
  --packages org.apache.iceberg:iceberg-spark-runtime-3.5_2.12:1.5.0,org.apache.iceberg:iceberg-aws-bundle:1.5.0 \
  --conf spark.sql.catalog.ontology=org.apache.iceberg.spark.SparkCatalog \
  --conf spark.sql.catalog.ontology.catalog-impl=org.apache.iceberg.nessie.NessieCatalog \
  --conf spark.sql.catalog.ontology.uri=http://localhost:19120/api/v2 \
  --conf spark.sql.catalog.ontology.ref=main \
  --conf spark.sql.catalog.ontology.warehouse=s3a://iceberg/warehouse \
  --conf spark.sql.catalog.ontology.io-impl=org.apache.iceberg.aws.s3.S3FileIO \
  --conf spark.sql.catalog.ontology.s3.endpoint=http://localhost:19000 \
  --conf spark.sql.catalog.ontology.s3.access-key-id=iceberg \
  --conf spark.sql.catalog.ontology.s3.secret-access-key=iceberg_password \
  --conf spark.sql.catalog.ontology.s3.path-style-access=true \
  --conf spark.sql.extensions=org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions
```

然后在 Spark Shell 中执行：

```scala
spark.sql("CREATE NAMESPACE IF NOT EXISTS ontology.grid")

spark.sql("""
CREATE TABLE IF NOT EXISTS ontology.grid.breaker_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    locked STRING,
    normal_open STRING,
    open STRING,
    retained STRING,
    switch_on_count STRING,
    switch_on_date STRING,
    aggregate STRING,
    in_service STRING,
    network_analysis_enabled STRING,
    normally_in_service STRING,
    alias_name STRING,
    description STRING,
    m_rid STRING,
    name STRING,
valid_from TIMESTAMP,
valid_to TIMESTAMP,
op_type STRING,
ingestion_ts TIMESTAMP
)
USING ICEBERG
PARTITIONED BY (actor_id)
""")

spark.sql("SHOW TABLES IN ontology.grid LIKE 'breaker*'").show()
```

## 方法4: 使用项目中的 bootstrap_iceberg.py（修改版）

可以修改 `scripts/bootstrap_iceberg.py`，添加创建 Breaker 表的逻辑。

## 验证表创建

执行以下 SQL 验证表是否创建成功：

```sql
SHOW TABLES IN ontology.grid LIKE 'breaker*';
DESCRIBE EXTENDED ontology.grid.breaker_snapshots;
```

## 注意事项

1. **前置条件**：
   - Nessie Catalog 服务运行在 `http://localhost:19120`
   - MinIO/S3 服务运行在 `http://localhost:19000`
   - Spark 已安装并配置了 Iceberg 扩展

2. **如果使用 Docker Compose**：
   ```bash
   cd infra
   docker-compose up -d
   ```

3. **SQL 文件位置**：
   - `scripts/create_breaker_table.sql` - 包含创建命名空间和表的 SQL
