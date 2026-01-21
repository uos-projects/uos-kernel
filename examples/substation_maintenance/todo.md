## 调度中心（DispatcherActor）与调度操作员（DispatcherOperatorActor）重构 TODO（对齐“两阶段 + 现实驱动事件”设计）

目标：像断路器 Actor 一样，把 **“命令发起”** 与 **“现实结果”** 分离开来：

- Command / Capability：表达“想让系统做什么”（意图）
- World Event / Feedback：通过 Binding 或协同消息感知“现实中确实发生了什么”
- Actor 只在“看到结果”后更新状态、发业务事件

---

### 1. DispatcherActor：从“直接处理事件+改任务状态”到“世界事件驱动的调度中枢”

**现状简述：**
- `Receive` 里直接 `switch` 处理：
  - `DeviceAbnormalEvent`
  - `MaintenanceRequiredEvent`
  - `MaintenanceCompletedEvent`
- 直接在 Handler 中：
  - 创建 `MaintenanceTask`
  - 更新 `pendingTasks`
  - 调用 `assignTaskToOperator` 立即把任务发给操作员

**重构思路：**

1. **划分命令 vs 事件：**
   - 命令（Capability Commands）：
     - `CreateMaintenanceTaskCommand`（由上层编排或自动规则触发）
     - `AssignTaskCommand`（将任务分配给某个操作员）
   - 事件（Coordination / World Events）：
     - `DeviceAbnormalEvent`（来自设备侧的世界事件）
     - `MaintenanceRequiredEvent`（定期计划触发的世界事件）
     - `MaintenanceCompletedEvent`（操作员反馈结果）

2. **Dispatcher 的 Capacity 层：**
   - 增加一个或多个 Capacity（逻辑上）：
     - `TaskPlanningCapacity`：
       - 输入：`DeviceAbnormalEvent`、`MaintenanceRequiredEvent`
       - 输出：内部命令 `CreateMaintenanceTaskCommand`
     - `TaskDispatchCapacity`：
       - 输入：`CreateMaintenanceTaskCommand`
       - 输出：选择合适操作员 → 生成 `AssignTaskCommand`
   - 这两个 Capacity 只“决定做什么”，不直接 mutate `pendingTasks`，而是：
     - 通过 Actor 的领域方法（如 `createTask`、`dispatchTask`）执行；
     - 或产生统一的内部消息，由 Actor 处理。

3. **Dispatcher 的 Actor 层：**
   - 引入更细粒度的领域方法（内部）：
     - `createEmergencyTaskFrom(event *DeviceAbnormalEvent) MaintenanceTask`
     - `createScheduledTaskFrom(event *MaintenanceRequiredEvent) MaintenanceTask`
     - `applyTaskCreated(task MaintenanceTask)`（更新 `pendingTasks`）
     - `applyTaskAssigned(task MaintenanceTask)`（更新 `pendingTasks` 中状态）
   - `Receive` 逻辑调整：
     - **世界事件路径**：
       - `DeviceAbnormalEvent` / `MaintenanceRequiredEvent` 作为 CoordinationEvent 进入
       - 通过 Capacity/领域方法创建任务 → 更新内部状态 →（可选）发出“任务创建”事件
     - **操作员反馈路径**：
       - `MaintenanceCompletedEvent` 只更新任务状态为 `completed/failed`，不直接驱动其他副作用

4. **事件输出显式化：**
   - 定义 Dispatcher 可能发出的高层事件：
     - `MaintenanceTaskCreated`
     - `MaintenanceTaskAssigned`
     - `MaintenanceTaskUpdated`（状态变化）
   - 在任务创建 / 分配 / 状态变更的领域方法中，通过统一的事件发射器输出，而不是在多个 handler 里散布。

---

### 2. DispatcherOperatorActor：从“直接执行检修 + 通知”到“基于任务 + 设备反馈的流程执行者”

**现状简述：**
- `Receive` 直接接收 `*MaintenanceTask`，立即：
  - 顺序发送 `OpenBreakerCommand` / `CloseBreakerCommand` 到各设备 Actor
  - 本地 sleep 模拟执行时间
  - 调用 `breaker.CompleteMaintenance()` 更新设备内部状态
  - 直接发送 `MaintenanceCompletedEvent` 给 Dispatcher

**重构思路：**

1. **将“检修任务执行”拆成多个阶段：**
   - 阶段 A：接收任务并“接单”
     - 任务进入 `currentTask`，状态从 `pending` → `in_progress`
     - 可发出 `TaskAcceptedEvent`（供监控/审计使用）
   - 阶段 B：发起停电控制
     - 对每个设备发送 `OpenBreakerCommand`
     - 不假设立即成功，只记录“控制已下发”
   - 阶段 C：等待/感知设备状态（理论上应通过设备事件流）
     - 当前示例中可以先简化为“假设下发即成功”
     - 未来可以订阅来自断路器的状态事件来确认
   - 阶段 D：执行检修逻辑（本地耗时）
   - 阶段 E：恢复供电（发送 `CloseBreakerCommand`）
   - 阶段 F：提交结果（发送 `MaintenanceCompletedEvent` 给 Dispatcher）

2. **命令 vs 事件：**
   - 命令：
     - `StartMaintenanceCommand`（从 Dispatcher 发给 Operator，而不是直接发送 `*MaintenanceTask` 指针）
     - `CompleteMaintenanceCommand`（Operator 内部生成，用于明确结束任务）
   - 事件：
     - `MaintenanceTaskStarted`（任务进入执行状态）
     - `MaintenanceTaskStepCompleted`（例如：停电完成 / 检修完成 / 恢复供电完成）
     - `MaintenanceCompletedEvent`（最终结果，已存在）

3. **Operator 的领域方法：**
   - `acceptTask(task MaintenanceTask)`：设置 `currentTask`，更新本地状态
   - `startMaintenanceFlow(ctx)`：按步骤执行检修流程（发命令给设备）
   - `onDeviceStateFeedback(...)`（未来扩展）：收到设备状态事件后推进流程状态机
   - `finishTask(result string)`：结束任务，发 `MaintenanceCompletedEvent`

4. **与断路器 Actor 的协同：**
   - Operator 不再直接调用 `breaker.CompleteMaintenance()`，而是：
     - 给每个设备发一个“检修完成”的命令或事件
     - 由设备 Actor 自己决定如何更新 `lastMaintenanceTime` / `operationHours` 等（领域内聚）
   - 当前示例可保留直接调用作为“模拟”，但在 TODO 中明确未来要经由消息/事件完成。

---

### 3. 渐进式落地步骤（实现顺序建议）

1. **DispatcherActor**
   - [ ] 把 `Receive` 中对 `DeviceAbnormalEvent` / `MaintenanceRequiredEvent` 的逻辑移动到领域方法（`createEmergencyTaskFrom` / `createScheduledTaskFrom`）
   - [ ] 引入 `MaintenanceTaskCreatedEvent` / `MaintenanceTaskAssignedEvent`，在任务变更时统一用事件发射器输出
   - [ ] 将“给操作员发任务”改为发送命令/消息（例如 `StartMaintenanceCommand`），而不是直接传 `*MaintenanceTask`

2. **DispatcherOperatorActor**
   - [ ] 将 `Receive(*MaintenanceTask)` 改为 `Receive(StartMaintenanceCommand)`，在内部通过命令携带的 TaskID 查询/构造任务
   - [ ] 将检修流程拆成若干内部方法（accept / startFlow / finish），并在关键步骤发出事件（`MaintenanceTaskStarted`、`MaintenanceCompletedEvent`）
   - [ ] 去掉直接调用 `breaker.CompleteMaintenance()`，改为发送“maintenance-done”类消息给设备 Actor（可先用简单命令表示）

3. **与设备反馈的整合（后续增强）**
   - [ ] 为 Operator 增加对设备状态事件的订阅（例如接收 `BreakerDeviceEvent` 或更通用的 Measurement/状态事件）
   - [ ] 用一个简单状态机驱动检修流程：只有当所有相关断路器状态确认为 `opened` 时，才进入检修阶段；所有恢复为 `closed` 时，才允许标记任务完成。

