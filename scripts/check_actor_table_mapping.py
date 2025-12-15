#!/usr/bin/env python3
"""
检查 Actor 属性与 Iceberg 表字段的对应关系
"""

import re

def to_snake_case(name: str) -> str:
    """转换为 snake_case（与生成器一致）"""
    s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
    return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()

# Actor 中设置的属性（CIM 原始命名）
actor_properties = {
    "mRID": "breaker-001",
    "name": "Main Breaker",
    "normalOpen": False,
    "open": False,
    "locked": False,
    "ratedCurrent": 1000.0,
    "description": "Main circuit breaker",
}

# Iceberg 表中的字段（从 SQL 提取）
table_fields = [
    "actor_id",
    "owl_class_uri",
    "sequence",
    "timestamp",
    "snapshot_id",
    "locked",
    "normal_open",
    "open",
    "retained",
    "switch_on_count",
    "switch_on_date",
    "aggregate",
    "in_service",
    "network_analysis_enabled",
    "normally_in_service",
    "alias_name",
    "description",
    "m_rid",
    "name",
    "valid_from",
    "valid_to",
    "op_type",
    "ingestion_ts",
]

print("=" * 60)
print("Actor 属性 vs Iceberg 表字段对应关系")
print("=" * 60)

print("\n1. Actor 属性（CIM 原始命名）：")
for prop_name in actor_properties.keys():
    print(f"   {prop_name}")

print("\n2. Iceberg 表字段（snake_case）：")
for field_name in table_fields:
    print(f"   {field_name}")

print("\n3. 对应关系分析：")
print("-" * 60)
matched = []
unmatched = []

for prop_name, prop_value in actor_properties.items():
    # 转换为 snake_case
    expected_field = to_snake_case(prop_name)
    
    if expected_field in table_fields:
        matched.append((prop_name, expected_field, prop_value))
        print(f"   ✓ {prop_name:20} → {expected_field:20} (匹配)")
    else:
        unmatched.append((prop_name, expected_field))
        print(f"   ✗ {prop_name:20} → {expected_field:20} (表中不存在)")

print("\n4. 匹配结果：")
print(f"   匹配: {len(matched)}/{len(actor_properties)}")
print(f"   不匹配: {len(unmatched)}/{len(actor_properties)}")

if unmatched:
    print("\n5. 不匹配的属性：")
    for prop_name, expected_field in unmatched:
        print(f"   - {prop_name} (期望字段: {expected_field})")
        # 检查是否有类似的字段
        similar = [f for f in table_fields if prop_name.lower() in f.lower() or expected_field in f]
        if similar:
            print(f"     可能的字段: {', '.join(similar)}")

print("\n6. 结论：")
print("-" * 60)
if len(matched) == len(actor_properties):
    print("   ✓ 所有 Actor 属性都有对应的 Iceberg 表字段")
    print("   ✓ 属性名通过 snake_case 转换可以匹配")
else:
    print("   ✗ 部分 Actor 属性没有对应的 Iceberg 表字段")
    print("   ⚠ 写入 Iceberg 时需要处理属性名到字段名的转换")
    print("   ⚠ 或者 Actor 应该使用 snake_case 属性名")

print("\n7. 建议：")
print("-" * 60)
print("   方案1: Actor 使用 CIM 原始属性名，写入 Iceberg 时转换为 snake_case")
print("   方案2: Actor 直接使用 snake_case 属性名（与表字段一致）")
print("   方案3: 在 CreateSnapshot 时自动转换属性名为表字段名")
