from __future__ import annotations

import os
from enum import Enum
from functools import lru_cache
from typing import Dict, List

from fastapi import FastAPI, HTTPException, Query
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse
from pyspark.sql import SparkSession, DataFrame


NESSIE_URI = "http://localhost:19120/api/v2"
S3_ENDPOINT = "http://localhost:19000"
WAREHOUSE = "s3a://iceberg/warehouse"
CATALOG_BRANCH = "main"
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
    return (
        SparkSession.builder.appName("OntologyTimeTravelAPI")
        .config("spark.sql.catalog.ontology", "org.apache.iceberg.spark.SparkCatalog")
        .config("spark.jars.packages", ICEBERG_PACKAGES)
        .config(
            "spark.sql.catalog.ontology.catalog-impl",
            "org.apache.iceberg.nessie.NessieCatalog",
        )
        .config("spark.sql.catalog.ontology.uri", NESSIE_URI)
        .config("spark.sql.catalog.ontology.ref", CATALOG_BRANCH)
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
        .config(
            "spark.sql.extensions",
            "org.apache.iceberg.spark.extensions.IcebergSparkSessionExtensions",
        )
        .getOrCreate()
    )


class TimeMode(str, Enum):
    biz = "biz"
    snapshot = "snapshot"
    timestamp = "timestamp"


def df_to_dicts(df: DataFrame) -> List[Dict]:
    return [row.asDict(recursive=True) for row in df.collect()]


_spark_session = None


def get_spark():
    global _spark_session
    if _spark_session is None:
        _spark_session = build_spark()
    return _spark_session


app = FastAPI(title="Ontology Time Travel API")
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Web 目录路径
WEB_DIR = os.path.join(os.path.dirname(os.path.dirname(__file__)), "web")

ENTITY_QUERIES = {
    "substation": "SELECT * FROM ontology.grid.substation {clause}",
    "voltage_level": "SELECT * FROM ontology.grid.voltage_level {clause}",
    "power_transformer": "SELECT * FROM ontology.grid.power_transformer {clause}",
    "breaker": "SELECT * FROM ontology.grid.breaker {clause}",
    "asset": "SELECT * FROM ontology.assets.asset {clause}",
    "work": "SELECT * FROM ontology.assets.work {clause}",
    "meter": "SELECT * FROM ontology.metering.meter {clause}",
    "usage_point": "SELECT * FROM ontology.metering.usage_point {clause}",
}

RELATION_QUERIES = [
    "SELECT * FROM ontology.relations.substation_voltage_level {clause}",
    "SELECT * FROM ontology.relations.transformer_substation {clause}",
    "SELECT * FROM ontology.relations.breaker_voltage_level {clause}",
    "SELECT * FROM ontology.relations.asset_resource {clause}",
    "SELECT * FROM ontology.relations.work_resource {clause}",
    "SELECT * FROM ontology.relations.meter_usage_point {clause}",
]


def build_clause(mode: TimeMode, value: str) -> str:
    if mode == TimeMode.biz:
        return f"WHERE '{value}' BETWEEN valid_from AND valid_to"
    if mode == TimeMode.snapshot:
        if not value.isdigit():
            raise HTTPException(status_code=400, detail="snapshot value 必须为数字")
        return f"VERSION AS OF {value}"
    if mode == TimeMode.timestamp:
        return f"TIMESTAMP AS OF '{value}'"
    raise HTTPException(status_code=400, detail="不支持的 mode")


@app.get("/api/view")
def get_view(
    mode: TimeMode = Query(TimeMode.biz, description="biz|snapshot|timestamp"),
    value: str = Query(..., description="biz/timestamp 使用 ISO8601，snapshot 使用 ID"),
):
    spark = get_spark()
    clause = build_clause(mode, value)
    entities: List[Dict] = []
    for entity, query in ENTITY_QUERIES.items():
        try:
            df = spark.sql(query.format(clause=clause))
            data = df_to_dicts(df)
            for item in data:
                item["entity_type"] = entity
            entities.extend(data)
        except Exception as e:
            # 如果表不存在或查询失败，跳过
            print(f"Warning: Failed to query {entity}: {e}")
            continue

    relations: List[Dict] = []
    for rel_query in RELATION_QUERIES:
        try:
            rel_df = spark.sql(rel_query.format(clause=clause))
            relations.extend(df_to_dicts(rel_df))
        except Exception as e:
            # 如果表不存在或查询失败，跳过
            print(f"Warning: Failed to query relation: {e}")
            continue
    
    return {"mode": mode, "value": value, "entities": entities, "relations": relations}


@lru_cache(maxsize=1)
def default_snapshot_table() -> str:
    return "ontology.grid.substation"


@app.get("/api/health")
def health():
    return {"status": "ok", "service": "ontology-time-travel-api"}


@app.get("/api/snapshots")
def list_snapshots(table: str = Query(None, description="Iceberg 表全名，不指定则返回所有表")):
    spark = get_spark()
    
    # 如果指定了表，只返回该表的 snapshots
    if table:
        metadata_table = f"{table}.snapshots"
        df = spark.table(metadata_table)
        snapshots = [
            {
                "snapshot_id": row["snapshot_id"],
                "committed_at": row["committed_at"],
                "operation": row["operation"],
                "table": table,
            }
            for row in df.collect()
        ]
        return {"table": table, "snapshots": snapshots}
    
    # 否则返回所有表的 snapshots
    all_tables = [
        "ontology.grid.substation",
        "ontology.grid.voltage_level",
        "ontology.grid.power_transformer",
        "ontology.grid.breaker",
        "ontology.assets.asset",
        "ontology.assets.work",
        "ontology.metering.meter",
        "ontology.metering.usage_point",
    ]
    
    all_snapshots = []
    for tbl in all_tables:
        try:
            metadata_table = f"{tbl}.snapshots"
            df = spark.table(metadata_table)
            for row in df.collect():
                all_snapshots.append({
                    "snapshot_id": row["snapshot_id"],
                    "committed_at": row["committed_at"],
                    "operation": row["operation"],
                    "table": tbl,
                })
        except Exception as e:
            # 如果表不存在或查询失败，跳过
            print(f"Warning: Failed to query {tbl}.snapshots: {e}")
            continue
    
    # 按时间排序
    all_snapshots.sort(key=lambda x: x["committed_at"])
    
    return {"table": "all", "snapshots": all_snapshots}


@app.on_event("shutdown")
def shutdown_spark():
    global _spark_session
    if _spark_session is not None:
        _spark_session.stop()
        _spark_session = None


# 静态文件服务（必须在所有 API 路由之后）
@app.get("/")
async def read_root():
    """返回前端页面"""
    return FileResponse(os.path.join(WEB_DIR, "index.html"))


@app.get("/{path:path}")
async def serve_static(path: str):
    """处理静态文件请求（CSS、JS 等），排除 API 路径"""
    if path.startswith("api/"):
        raise HTTPException(status_code=404, detail="Not found")
    # 空路径或根路径返回 index.html
    if not path or path == "/":
        return FileResponse(os.path.join(WEB_DIR, "index.html"))
    file_path = os.path.join(WEB_DIR, path)
    if os.path.isfile(file_path):
        return FileResponse(file_path)
    # 如果不是文件，返回 index.html（用于前端路由）
    return FileResponse(os.path.join(WEB_DIR, "index.html"))


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("server.app:app", host="0.0.0.0", port=8000, reload=False)

