#!/usr/bin/env python3
"""测试 Nessie API 访问"""

import requests
import json

NESSIE_URI = "http://localhost:19120/api/v2"

def test_nessie_api():
    """测试 Nessie API 的基本功能"""
    
    print("=" * 60)
    print("Nessie API 测试")
    print("=" * 60)
    
    # 1. 检查 Nessie 是否运行
    try:
        response = requests.get(f"{NESSIE_URI}/config")
        print(f"\n✅ Nessie 服务运行正常")
        print(f"   状态码: {response.status_code}")
        if response.status_code == 200:
            config = response.json()
            print(f"   配置: {json.dumps(config, indent=2, ensure_ascii=False)}")
    except Exception as e:
        print(f"\n❌ 无法连接到 Nessie: {e}")
        print(f"   请确保 Nessie 运行在 {NESSIE_URI}")
        return
    
    # 2. 查看所有分支
    print("\n" + "=" * 60)
    print("分支列表")
    print("=" * 60)
    try:
        response = requests.get(f"{NESSIE_URI}/trees")
        if response.status_code == 200:
            branches = response.json()
            print(f"   找到 {len(branches.get('references', []))} 个分支:")
            for ref in branches.get('references', []):
                print(f"   - {ref.get('name')} ({ref.get('type')})")
        else:
            print(f"   状态码: {response.status_code}")
            print(f"   响应: {response.text}")
    except Exception as e:
        print(f"   ❌ 错误: {e}")
    
    # 3. 查看 main 分支的表
    print("\n" + "=" * 60)
    print("main 分支的表列表")
    print("=" * 60)
    try:
        response = requests.get(f"{NESSIE_URI}/trees/main/entries")
        if response.status_code == 200:
            entries = response.json()
            if entries.get('entries'):
                print(f"   找到 {len(entries['entries'])} 个表:")
                for entry in entries['entries']:
                    name = entry.get('name', {}).get('elements', [])
                    name_str = '.'.join(name) if name else 'unknown'
                    print(f"   - {name_str}")
                    print(f"     类型: {entry.get('type')}")
            else:
                print("   (暂无表)")
        else:
            print(f"   状态码: {response.status_code}")
            print(f"   响应: {response.text}")
    except Exception as e:
        print(f"   ❌ 错误: {e}")
    
    # 4. 查看 API 文档端点
    print("\n" + "=" * 60)
    print("API 信息")
    print("=" * 60)
    print(f"   Nessie API 端点: {NESSIE_URI}")
    print(f"   OpenAPI 文档: {NESSIE_URI}/openapi.json")
    print(f"   在浏览器中访问: {NESSIE_URI}/trees")
    
    print("\n" + "=" * 60)
    print("测试完成")
    print("=" * 60)

if __name__ == "__main__":
    test_nessie_api()

