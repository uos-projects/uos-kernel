"""变电站停电检修流程 — Python SDK 示例

镜像 Go 版本 examples/substation_maintenance/ 的完整逻辑。
前提：Go gRPC 服务端已启动（examples/substation_maintenance_grpc）

架构：
  Actor 层（Go）：BreakerActor 持续监测设备状态，暴露开/合能力
  Kernel 层（Go）：Open / Watch / Ioctl / Read 等 POSIX 风格 API
  gRPC 层：将 Kernel API 1:1 映射为 gRPC 服务
  App 层（Python）：MaintenanceProcess 通过 KernelClient 编排检修工作流
"""

import time
import sys
import os

# 添加 SDK 路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))

from uos_kernel import KernelClient
from uos_kernel.process import Process
from uos_kernel.helpers import wait_for_property, execute_capacity


# ---------------------------------------------------------------------------
# 设备状态显示
# ---------------------------------------------------------------------------

def display_device_status(client: KernelClient, fd: int) -> None:
    """格式化显示设备状态，镜像 Go 版 displayDeviceStatus。"""
    state = client.read(fd)
    props = state["properties"]

    rid = state["resource_id"]
    name = props.get("name", "未知")
    is_open = props.get("isOpen", False)
    voltage = props.get("voltage", 0.0)
    current = props.get("current", 0.0)
    temperature = props.get("temperature", 0.0)
    hours = props.get("operationHours", 0)
    last_maint = props.get("lastMaintenanceTime", "未知")

    status_text = "打开" if is_open else "关闭"
    print(f"{rid} ({name}):")
    print(f"  状态：{status_text}")
    print(f"  电压：{voltage:.2f} kV")
    print(f"  电流：{current:.2f} A")
    print(f"  温度：{temperature:.2f} °C")
    # operationHours 通过 protobuf Struct 传输时可能是 float
    print(f"  运行小时数：{int(hours)} 小时")
    print(f"  上次检修：{last_maint}")
    print()


def display_capabilities(client: KernelClient, fd: int) -> None:
    """显示资源的能力和事件定义。"""
    info = client.stat(fd)
    caps = info.get("capabilities", [])
    events = info.get("events", [])

    if caps:
        print("  能力列表：")
        for c in caps:
            ops = ", ".join(c.get("operations", []))
            desc = c.get("description", "")
            print(f"    - {c['name']} [{ops}] {desc}")

    if events:
        print("  事件列表：")
        for e in events:
            print(f"    - {e['name']} ({e.get('event_type', '')}) {e.get('description', '')}")


# ---------------------------------------------------------------------------
# 业务流程
# ---------------------------------------------------------------------------

class MaintenanceProcess(Process):
    """变电站停电检修流程

    实现 Process 接口，业务逻辑全部在应用层。
    通过 KernelClient (gRPC) 操控设备，不直接接触 Actor。
    """

    def __init__(self, breaker_ids: list[str]):
        self.breaker_ids = breaker_ids

    def name(self) -> str:
        return "变电站停电检修"

    def run(self, client: KernelClient) -> None:
        # 打开所有断路器资源
        fds: dict[str, int] = {}
        for bid in self.breaker_ids:
            fd = client.open("Breaker", bid, 0)
            fds[bid] = fd

        try:
            primary_id = self.breaker_ids[0]
            primary_fd = fds[primary_id]

            # 显示初始设备状态
            print("\n=== 初始状态 ===")
            for bid, fd in fds.items():
                display_device_status(client, fd)

            # 显示设备能力和事件
            print("=== 设备能力与事件 ===")
            display_capabilities(client, primary_fd)

            print(f"\n[{self.name()}] 开始监听设备事件...")

            # Watch 所有事件，过滤出自定义事件
            for event in client.watch(primary_fd):
                if event["type"] != "custom":
                    continue
                data = event.get("data", {})
                if data.get("_type") == "DeviceAbnormalEvent":
                    print(f"\n[{self.name()}] 收到设备异常事件，启动检修工作流...")
                    print(f"  设备: {data.get('device_id')}")
                    print(f"  类型: {data.get('event_type')}")
                    print(f"  严重性: {data.get('severity')}")
                    details = data.get("details", {})
                    if details:
                        print(f"  详情: 温度 {details.get('temperature', '?')}°C"
                              f"（阈值 {details.get('threshold', '?')}°C）")
                    self._execute_workflow(client, fds)
                    return
        finally:
            for fd in fds.values():
                client.close(fd)

    def _execute_workflow(self, client: KernelClient, fds: dict[str, int]) -> None:
        """执行停电检修工作流（四步）"""
        print("\n========== 开始停电检修流程 ==========")

        # 步骤 1：停电 — 打开所有断路器
        print("\n【步骤 1】停电操作")
        for bid, fd in fds.items():
            print(f"  打开断路器 {bid}...")
            execute_capacity(client, fd, "OpenBreakerCommand", {
                "reason": "停电检修",
                "operator": "python-maintenance-process",
            })
            wait_for_property(client, fd, "isOpen", True)
            print(f"  ✓ 断路器 {bid} 已打开")

        # 步骤 2：执行检修
        print("\n【步骤 2】执行检修操作")
        print("  检修中...")
        time.sleep(1)  # 模拟检修耗时
        print("  ✓ 检修完成")

        # 步骤 3：复电 — 关闭所有断路器
        print("\n【步骤 3】恢复供电")
        for bid, fd in fds.items():
            print(f"  关闭断路器 {bid}...")
            execute_capacity(client, fd, "CloseBreakerCommand", {
                "reason": "检修完成，恢复供电",
                "operator": "python-maintenance-process",
            })
            wait_for_property(client, fd, "isOpen", False)
            print(f"  ✓ 断路器 {bid} 已关闭")

        # 步骤 4：完成检修记录
        print("\n【步骤 4】更新检修记录")
        for bid, fd in fds.items():
            result = execute_capacity(client, fd, "CompleteMaintenanceCommand", {
                "operator": "python-maintenance-process",
                "result": "success",
            })
            if result is None:
                print(f"  ✓ 断路器 {bid} 检修记录已更新")
            else:
                print(f"  ⚠ 断路器 {bid} 更新记录返回: {result}")

        time.sleep(0.2)  # 等待异步命令处理

        # 显示最终状态
        print("\n=== 最终状态 ===")
        for bid, fd in fds.items():
            display_device_status(client, fd)

        print("========== 停电检修流程完成 ==========")


# ---------------------------------------------------------------------------
# 入口
# ---------------------------------------------------------------------------

def main():
    addr = "localhost:50051"

    print("=== 变电站停电检修操作场景演示（Python SDK）===")
    print()
    print("架构说明：")
    print("  Actor 层：BreakerActor（设备资源）持续监测状态，暴露开/合能力")
    print("  Kernel 层：Open / Watch / Ioctl / Read 等 POSIX 风格 API")
    print("  gRPC 层：将 Kernel API 映射为远程服务")
    print("  App 层：MaintenanceProcess（Python）通过 KernelClient 编排检修工作流")
    print()

    print(f"连接 gRPC 服务端 {addr}...")

    with KernelClient(addr) as client:
        process = MaintenanceProcess(["BREAKER-001"])
        process.run(client)

    print("\n=== 演示完成 ===")


if __name__ == "__main__":
    main()
