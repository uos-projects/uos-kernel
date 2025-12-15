#!/usr/bin/env python3
"""
使用 PySpark 执行 create_breaker_table.sql 在 Iceberg 中创建表
"""

import os
from pathlib import Path
from pyspark.sql import SparkSession

# Iceberg 配置
NESSIE_URI = "http://localhost:19120/api/v2"
S3_ENDPOINT = "http://localhost:19000"
WAREHOUSE = "s3a://iceberg/warehouse"
BRANCH = "main"
ICEBERG_PACKAGES = ",".join(
    [
        "org.apache.iceberg:iceberg-spark-runtime-3.5_2.12:1.5.0",
        "org.apache.iceberg:iceberg-aws-bundle:1.5.0",
    ]
)

os.environ.setdefault("AWS_REGION", "us-east-1")
os.environ.setdefault("AWS_ACCESS_KEY_ID", "iceberg")
os.environ.setdefault("AWS_SECRET_ACCESS_KEY", "iceberg_password")


def build_spark() -> SparkSession:
    """构建 SparkSession"""
    return (
        SparkSession.builder.appName("CreateBreakerTable")
        .config("spark.sql.catalog.ontology", "org.apache.iceberg.spark.SparkCatalog")
        .config("spark.jars.packages", ICEBERG_PACKAGES)
        .config(
            "spark.sql.catalog.ontology.catalog-impl",
            "org.apache.iceberg.nessie.NessieCatalog",
        )
        .config("spark.sql.catalog.ontology.uri", NESSIE_URI)
        .config("spark.sql.catalog.ontology.ref", BRANCH)
        .config("spark.sql.catalog.ontology.warehouse", WAREHOUSE)
        .config("spark.sql.catalog.ontology.io-impl", "org.apache.iceberg.aws.s3.S3FileIO")
        .config("spark.sql.catalog.ontology.s3.endpoint", S3_ENDPOINT)
        .config("spark.sql.catalog.ontology.s3.access-key-id", "iceberg")
        .config("spark.sql.catalog.ontology.s3.secret-access-key", "iceberg_password")
        .config("spark.sql.catalog.ontology.s3.path-style-access", "true")
        .config("spark.sql.catalog.ontology.s3.region", "us-east-1")
        .config("spark.hadoop.fs.s3a.endpoint", S3_ENDPOINT)
        .config("spark.hadoop.fs.s3a.path.style.access", "true")
        .config("spark.hadoop.fs.s3a.access.key", "iceberg")
        .config("spark.hadoop.fs.s3a.secret.key", "iceberg_password")
        .config("spark.sql.extensions", "org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions")
        .getOrCreate()
    )


def parse_sql_file(sql_file: str) -> list:
    """解析 SQL 文件，提取 SQL 语句"""
    with open(sql_file, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 按分号分割 SQL 语句
    statements = []
    current_statement = []
    
    for line in content.split('\n'):
        # 跳过注释行
        stripped = line.strip()
        if stripped.startswith('--'):
            continue
        
        # 如果行不为空，添加到当前语句
        if stripped:
            current_statement.append(line)
        
        # 如果行以分号结尾，表示一个完整的 SQL 语句
        if stripped.endswith(';'):
            # 合并当前语句
            sql = '\n'.join(current_statement).strip()
            if sql:
                # 移除末尾的分号（Spark SQL 不需要）
                sql = sql.rstrip(';').strip()
                statements.append(sql)
            current_statement = []
    
    # 处理最后一个语句（可能没有分号）
    if current_statement:
        sql = '\n'.join(current_statement).strip()
        if sql:
            sql = sql.rstrip(';').strip()
            statements.append(sql)
    
    return statements


def check_services():
    """检查必要的服务是否运行"""
    import socket
    import urllib.request
    
    services = {
        "Nessie": ("localhost", 19120),
        "MinIO": ("localhost", 19000),
    }
    
    print("\n[0/3] Checking services...")
    all_ok = True
    for name, (host, port) in services.items():
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(2)
            result = sock.connect_ex((host, port))
            sock.close()
            if result == 0:
                print(f"  ✓ {name} is running on {host}:{port}")
            else:
                print(f"  ✗ {name} is not accessible on {host}:{port}")
                all_ok = False
        except Exception as e:
            print(f"  ✗ {name} check failed: {e}")
            all_ok = False
    
    if not all_ok:
        print("\n⚠ Warning: Some services are not running.")
        print("Please start services using:")
        print("  cd infra && docker-compose up -d")
        print("\nContinuing anyway...")
    
    return all_ok


def main():
    print("=" * 60)
    print("Create Breaker Table in Iceberg")
    print("=" * 60)
    
    sql_file = "scripts/create_breaker_table.sql"
    
    if not Path(sql_file).exists():
        print(f"Error: SQL file not found: {sql_file}")
        return
    
    # 检查服务
    check_services()
    
    # 解析 SQL 文件
    print(f"\n[1/3] Parsing SQL file: {sql_file}")
    sql_statements = parse_sql_file(sql_file)
    print(f"  Found {len(sql_statements)} SQL statements")
    
    # 构建 SparkSession
    print("\n[2/3] Building SparkSession...")
    try:
        spark = build_spark()
        print("  ✓ SparkSession created")
    except Exception as e:
        print(f"  ✗ Failed to create SparkSession: {e}")
        print("\nNote: Make sure Spark and Iceberg dependencies are available.")
        print("You may need to install PySpark:")
        print("  pip install pyspark")
        return
    
    # 执行 SQL 语句
    print("\n[3/3] Executing SQL statements...")
    print("-" * 60)
    
    for i, sql in enumerate(sql_statements, 1):
        print(f"\n[{i}/{len(sql_statements)}] Executing:")
        # 显示 SQL 的前 100 个字符
        sql_preview = sql[:100] + "..." if len(sql) > 100 else sql
        print(f"  {sql_preview}")
        
        try:
            result = spark.sql(sql)
            
            # 如果是查询语句，显示结果
            if sql.strip().upper().startswith('SHOW') or sql.strip().upper().startswith('SELECT'):
                print("\n  Result:")
                result.show(truncate=False)
            else:
                print("  ✓ Success")
                
        except Exception as e:
            print(f"  ✗ Error: {e}")
            # 继续执行其他语句
            continue
    
    print("\n" + "=" * 60)
    print("Completed!")
    print("=" * 60)
    
    # 验证表是否创建成功
    print("\nVerifying table creation...")
    try:
        tables = spark.sql("SHOW TABLES IN ontology.grid LIKE 'breaker*'").collect()
        if tables:
            print("\n✓ Tables found:")
            for table in tables:
                print(f"  - {table['namespace']}.{table['tableName']}")
        else:
            print("\n⚠ No breaker tables found")
    except Exception as e:
        print(f"\n⚠ Could not verify: {e}")
    
    spark.stop()


if __name__ == "__main__":
    main()
