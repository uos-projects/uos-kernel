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

# 导入语义模型和映射模块
import sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).parent.parent))

from ontology.parser import get_semantic_model
from ontology.model import Entity
from mapping.mapper import EntityMapper
from table_gen.generator import SQLGenerator


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


# ==================== 语义模型管理 API ====================

@app.get("/api/ontology/packages")
def get_packages():
    """获取所有包"""
    model = get_semantic_model()
    return {
        "packages": [
            {
                "id": pkg.id,
                "name": pkg.name,
                "summary": pkg.summary,
                "entity_count": len(pkg.entities),
            }
            for pkg in model.packages
        ]
    }


@app.get("/api/ontology/entities")
def get_entities(package_id: str = Query(None, description="包 ID，可选")):
    """获取所有实体"""
    model = get_semantic_model()
    entities = model.get_all_entities()
    
    if package_id:
        entities = [e for e in entities if e.package_id == package_id]
    
    return {
        "entities": [
            {
                "name": entity.name,
                "package_id": entity.package_id,
                "package_name": entity.package_name,
                "inherits": entity.inherits,
                "key_attributes": entity.key_attributes,
                "relationships": [
                    {"type": rel.type, "target": rel.target}
                    for rel in entity.relationships
                ],
            }
            for entity in entities
        ]
    }


@app.get("/api/ontology/entity/{entity_name}")
def get_entity(entity_name: str):
    """获取实体详情"""
    model = get_semantic_model()
    entity = model.get_entity(entity_name)
    
    if entity is None:
        raise HTTPException(status_code=404, detail=f"实体 {entity_name} 不存在")
    
    entity_registry = model.get_entity_registry()
    all_attributes = entity.get_all_attributes(entity_registry)
    
    return {
        "name": entity.name,
        "package_id": entity.package_id,
        "package_name": entity.package_name,
        "inherits": entity.inherits,
        "key_attributes": entity.key_attributes,
        "all_attributes": all_attributes,
        "relationships": [
            {"type": rel.type, "target": rel.target}
            for rel in entity.relationships
        ],
    }


# ==================== 映射和表生成 API ====================

@app.get("/api/mapping/entity/{entity_name}")
def get_table_schema(entity_name: str):
    """获取实体的表结构（不创建表）"""
    model = get_semantic_model()
    entity = model.get_entity(entity_name)
    
    if entity is None:
        raise HTTPException(status_code=404, detail=f"实体 {entity_name} 不存在")
    
    mapper = EntityMapper(model)
    schema = mapper.map_entity_to_table(entity)
    
    # 生成 SQL
    sqls = SQLGenerator.generate_table_sqls(schema)
    
    return {
        "entity_name": entity_name,
        "table_name": schema.get_full_table_name(),
        "namespace": schema.namespace,
        "columns": [
            {"name": col.name, "type": col.type, "comment": col.comment}
            for col in schema.columns
        ],
        "partition_fields": schema.partition_fields,
        "sql": {
            "namespace": sqls[0],
            "table": sqls[1] if len(sqls) > 1 else "",
        },
    }


@app.post("/api/tables/create/{entity_name}")
def create_table(entity_name: str, dry_run: bool = Query(False, description="是否只是预览，不实际创建")):
    """根据实体创建 Iceberg 表"""
    model = get_semantic_model()
    entity = model.get_entity(entity_name)
    
    if entity is None:
        raise HTTPException(status_code=404, detail=f"实体 {entity_name} 不存在")
    
    mapper = EntityMapper(model)
    schema = mapper.map_entity_to_table(entity)
    sqls = SQLGenerator.generate_table_sqls(schema)
    
    if dry_run:
        return {
            "entity_name": entity_name,
            "table_name": schema.get_full_table_name(),
            "sql": sqls,
            "status": "preview",
        }
    
    # 实际创建表
    spark = get_spark()
    results = []
    
    try:
        for sql in sqls:
            spark.sql(sql)
            results.append({"sql": sql, "status": "success"})
        
        return {
            "entity_name": entity_name,
            "table_name": schema.get_full_table_name(),
            "status": "created",
            "results": results,
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"创建表失败: {str(e)}")


@app.get("/api/tables/mappings")
def get_all_mappings():
    """获取所有实体到表的映射关系"""
    model = get_semantic_model()
    mapper = EntityMapper(model)
    schemas = mapper.map_all_entities()
    
    return {
        "mappings": [
            {
                "entity_name": entity_name,
                "table_name": schema.get_full_table_name(),
                "namespace": schema.namespace,
                "column_count": len(schema.columns),
                "partition_fields": schema.partition_fields,
            }
            for entity_name, schema in schemas.items()
        ]
    }


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

