package actors

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/uos-projects/uos-kernel/actors/capacities"
)

// LifecycleState 生命周期状态
type LifecycleState string

const (
	StateOnline    LifecycleState = "online"    // 在线
	StateOffline   LifecycleState = "offline"   // 离线
	StateSuspended LifecycleState = "suspended" // 挂起
	StateBusy      LifecycleState = "busy"      // 忙碌
)

// StateTransition 状态转换记录
type StateTransition struct {
	FromState LifecycleState
	ToState   LifecycleState
	Timestamp int64
	Trigger   string
}

// BehaviorState 行为状态（正在执行的命令/场景步骤）
type BehaviorState struct {
	// CommandID 正在执行的命令 ID
	CommandID string
	// CommandType 命令类型
	CommandType string
	// ScenarioStep 场景步骤（如果有）
	ScenarioStep string
	// StartedAt 开始时间
	StartedAt int64
}

// BaseResourceActor 资源 Actor 基类
// 提供能力管理（Capacity）、外部绑定（Binding）和事件管理功能，是所有资源 Actor 的基础
type BaseResourceActor struct {
	*BaseActor
	resourceID   string
	resourceType string
	capabilities map[string]capacities.Capacity // 能力名称 -> Capacity 实现
	bindings     map[BindingType]Binding        // 绑定类型 -> Binding 实现
	events       map[string]*EventDescriptor    // 事件名称 -> EventDescriptor

	// 生命周期状态机
	currentState LifecycleState
	stateHistory []StateTransition // 状态变化历史
	stateMu      sync.RWMutex      // 状态锁

	// 行为状态（正在执行的命令/场景步骤）
	behaviorState *BehaviorState
	behaviorMu    sync.RWMutex // 行为状态锁

	// 事件发射器
	eventEmitter *BaseEventEmitter
	eventsMu     sync.RWMutex // 事件注册表锁
}

// NewBaseResourceActor 创建一个新的基础资源 Actor
func NewBaseResourceActor(
	id string,
	resourceType string,
) *BaseResourceActor {
	baseActor := NewBaseActor(id)
	actor := &BaseResourceActor{
		BaseActor:     baseActor,
		resourceID:    id,
		resourceType:  resourceType,
		capabilities:  make(map[string]capacities.Capacity),
		bindings:      make(map[BindingType]Binding),
		events:        make(map[string]*EventDescriptor),
		currentState:  StateOffline, // 初始状态为离线
		stateHistory:  make([]StateTransition, 0),
		behaviorState: nil, // 初始无行为状态
		eventEmitter:  nil, // 将在 SetSystem 时初始化
	}

	// 注册标准生命周期事件
	actor.registerStandardEvents()

	return actor
}

// ID 实现 Actor 接口，返回资源 ID
func (a *BaseResourceActor) ID() string {
	return a.resourceID
}

// ResourceID 返回资源 ID
func (a *BaseResourceActor) ResourceID() string {
	return a.resourceID
}

// ResourceType 返回资源类型
func (a *BaseResourceActor) ResourceType() string {
	return a.resourceType
}

// AddCapacity 添加一个能力
func (a *BaseResourceActor) AddCapacity(capacity capacities.Capacity) {
	a.capabilities[capacity.Name()] = capacity
}

// RemoveCapacity 移除一个能力
func (a *BaseResourceActor) RemoveCapacity(capacityName string) {
	delete(a.capabilities, capacityName)
}

// HasCapacity 检查是否具有某种能力
func (a *BaseResourceActor) HasCapacity(capacityName string) bool {
	_, exists := a.capabilities[capacityName]
	return exists
}

// GetCapacity 获取指定名称的能力
func (a *BaseResourceActor) GetCapacity(capacityName string) (capacities.Capacity, bool) {
	cap, exists := a.capabilities[capacityName]
	return cap, exists
}

// ListCapabilities 返回所有能力名称列表
func (a *BaseResourceActor) ListCapabilities() []string {
	names := make([]string, 0, len(a.capabilities))
	for name := range a.capabilities {
		names = append(names, name)
	}
	return names
}

// AddBinding 添加外部绑定
func (a *BaseResourceActor) AddBinding(binding Binding) error {
	if binding.ResourceID() != a.resourceID {
		return fmt.Errorf("binding resource ID mismatch: expected %s, got %s", a.resourceID, binding.ResourceID())
	}

	// 设置 Actor 引用（如果绑定支持）
	if setter, ok := binding.(ActorRefSetter); ok {
		setter.SetActorRef(a)
	}

	a.bindings[binding.Type()] = binding
	return nil
}

// RemoveBinding 移除外部绑定
func (a *BaseResourceActor) RemoveBinding(bindingType BindingType) error {
	binding, exists := a.bindings[bindingType]
	if !exists {
		return fmt.Errorf("binding type %s not found", bindingType)
	}

	if err := binding.Stop(); err != nil {
		return fmt.Errorf("failed to stop binding: %w", err)
	}

	delete(a.bindings, bindingType)
	return nil
}

// GetBinding 获取指定类型的绑定
func (a *BaseResourceActor) GetBinding(bindingType BindingType) (Binding, bool) {
	binding, exists := a.bindings[bindingType]
	return binding, exists
}

// ListBindings 返回所有绑定类型列表
func (a *BaseResourceActor) ListBindings() []BindingType {
	types := make([]BindingType, 0, len(a.bindings))
	for bindingType := range a.bindings {
		types = append(types, bindingType)
	}
	return types
}

// CurrentState 获取当前生命周期状态
func (a *BaseResourceActor) CurrentState() LifecycleState {
	a.stateMu.RLock()
	defer a.stateMu.RUnlock()
	return a.currentState
}

// CanTransition 检查是否可以转换到指定状态
func (a *BaseResourceActor) CanTransition(to LifecycleState) bool {
	a.stateMu.RLock()
	current := a.currentState
	a.stateMu.RUnlock()

	// 定义允许的状态转换
	allowedTransitions := map[LifecycleState][]LifecycleState{
		StateOffline:   {StateOnline, StateSuspended},
		StateOnline:    {StateOffline, StateSuspended, StateBusy},
		StateSuspended: {StateOnline, StateOffline},
		StateBusy:      {StateOnline, StateOffline},
	}

	allowed, exists := allowedTransitions[current]
	if !exists {
		return false
	}

	for _, state := range allowed {
		if state == to {
			return true
		}
	}

	return false
}

// SetSystem 设置 System 引用（重写基类方法，同时初始化事件发射器）
func (a *BaseResourceActor) SetSystem(system *System) {
	a.BaseActor.SetSystem(system)
	if a.eventEmitter == nil {
		a.eventEmitter = NewBaseEventEmitter(a.resourceID, system)
	}
}

// transition 转换状态（内部方法，只能通过消息调用）
func (a *BaseResourceActor) transition(to LifecycleState, trigger string) error {
	a.stateMu.Lock()
	defer a.stateMu.Unlock()

	current := a.currentState

	// 定义允许的状态转换
	allowedTransitions := map[LifecycleState][]LifecycleState{
		StateOffline:   {StateOnline, StateSuspended},
		StateOnline:    {StateOffline, StateSuspended, StateBusy},
		StateSuspended: {StateOnline, StateOffline},
		StateBusy:      {StateOnline, StateOffline},
	}

	allowed, exists := allowedTransitions[current]
	if !exists {
		return fmt.Errorf("invalid current state: %s", current)
	}

	// 检查是否允许转换
	canTransition := false
	for _, state := range allowed {
		if state == to {
			canTransition = true
			break
		}
	}

	if !canTransition {
		return fmt.Errorf("invalid state transition from %s to %s", current, to)
	}

	fromState := current

	// 记录状态转换
	transition := StateTransition{
		FromState: fromState,
		ToState:   to,
		Timestamp: 0, // TODO: 使用实际时间戳
		Trigger:   trigger,
	}

	// 添加到历史记录（保留最近100条）
	a.stateHistory = append(a.stateHistory, transition)
	if len(a.stateHistory) > 100 {
		a.stateHistory = a.stateHistory[1:]
	}

	// 更新当前状态
	a.currentState = to

	return nil
}

// GetStateHistory 获取状态变化历史
func (a *BaseResourceActor) GetStateHistory() []StateTransition {
	a.stateMu.RLock()
	defer a.stateMu.RUnlock()

	// 返回副本
	result := make([]StateTransition, len(a.stateHistory))
	copy(result, a.stateHistory)
	return result
}

// Receive 重写消息接收逻辑，按照消息分类路由
// 核心执行逻辑：
//
//	if CapabilityCommand: check guard → execute or reject
//	else if CoordinationEvent: update state / notify
//	else if Internal: internal handling
func (a *BaseResourceActor) Receive(ctx context.Context, msg Message) error {
	// 根据消息类型分类处理
	switch msg.MessageType() {
	case MessageCategoryCapabilityCommand:
		return a.handleCapabilityCommand(ctx, msg)
	case MessageCategoryCoordinationEvent:
		return a.handleCoordinationEvent(ctx, msg)
	case MessageCategoryInternal:
		return a.handleInternalMessage(ctx, msg)
	default:
		// 兼容旧代码：如果消息没有实现 Message 接口，尝试按旧方式处理
		return a.handleLegacyMessage(ctx, msg)
	}
}

// handleCapabilityCommand 处理能力命令
func (a *BaseResourceActor) handleCapabilityCommand(ctx context.Context, msg Message) error {
	// 转换为 Capacity Message
	var capMsg capacities.Message = msg

	// 找到能处理此消息的 Capacity
	var targetCapacity capacities.Capacity
	for _, capacity := range a.capabilities {
		if capacity.CanHandle(capMsg) {
			targetCapacity = capacity
			break
		}
	}

	if targetCapacity == nil {
		return fmt.Errorf("no capacity can handle command %T", msg)
	}

	// 检查 Guards（如果 Capacity 实现了 Guards）
	if guards, ok := targetCapacity.(capacities.Guards); ok {
		satisfied, failed, err := guards.CheckGuards(ctx, a)
		if err != nil {
			return fmt.Errorf("error checking guards: %w", err)
		}
		if !satisfied {
			// 发射命令被拒绝事件
			if cmd, ok := msg.(CapabilityCommand); ok && a.eventEmitter != nil {
				_ = a.eventEmitter.EmitCommandFailed(cmd.CommandID(), fmt.Errorf("guards not satisfied: %v", failed))
			}
			return fmt.Errorf("guards not satisfied: %v", failed)
		}
	}

	// 记录行为状态（如果命令实现了 CapabilityCommand）
	if cmd, ok := msg.(CapabilityCommand); ok {
		a.setBehaviorState(&BehaviorState{
			CommandID:   cmd.CommandID(),
			CommandType: fmt.Sprintf("%T", msg),
			StartedAt:   0, // TODO: 使用实际时间戳
		})
		defer a.clearBehaviorState()
	}

	// 执行 Capacity
	err := targetCapacity.Execute(ctx, capMsg)

	// 发射执行结果事件
	if cmd, ok := msg.(CapabilityCommand); ok && a.eventEmitter != nil {
		if err != nil {
			_ = a.eventEmitter.EmitCommandFailed(cmd.CommandID(), err)
		} else {
			_ = a.eventEmitter.EmitCommandCompleted(cmd.CommandID(), nil)
		}
	}

	// 应用 Effects（如果 Capacity 实现了 Effects）
	if err == nil {
		if effects, ok := targetCapacity.(capacities.Effects); ok {
			a.applyEffects(effects.Effects())
		}
	}

	return err
}

// handleCoordinationEvent 处理协同事件
func (a *BaseResourceActor) handleCoordinationEvent(ctx context.Context, msg Message) error {
	switch m := msg.(type) {
	case *ExternalEventMessage:
		// 外部事件：更新状态或通知
		// 这里可以根据具体业务逻辑处理
		return nil
	case *ExecuteExternalCommandMessage:
		// 执行外部命令：通过 Binding 执行
		binding, exists := a.bindings[m.BindingType]
		if !exists {
			return fmt.Errorf("binding type %s not found", m.BindingType)
		}
		return binding.ExecuteExternal(ctx, m.Command)
	default:
		// 其他协同事件：默认处理
		return nil
	}
}

// handleInternalMessage 处理内部消息
func (a *BaseResourceActor) handleInternalMessage(ctx context.Context, msg Message) error {
	switch m := msg.(type) {
	case *TransitionStateMessage:
		from := a.currentState
		err := a.transition(m.ToState, m.Trigger)
		if err == nil && a.eventEmitter != nil {
			// 发射状态变化事件
			_ = a.eventEmitter.EmitStateChanged(from, m.ToState, m.Trigger)
		}
		return err
	case *GetStateMessage:
		// 状态查询消息（这里简化处理，实际可以通过响应消息返回）
		return nil
	case *SuspendMessage:
		return a.transition(StateSuspended, m.Reason)
	case *ResumeMessage:
		return a.transition(StateOnline, m.Reason)
	default:
		return nil
	}
}

// handleLegacyMessage 处理旧格式消息（向后兼容）
// 注意：这个方法接收 interface{} 而不是 Message，用于处理未实现 Message 接口的旧消息
func (a *BaseResourceActor) handleLegacyMessage(ctx context.Context, msg interface{}) error {
	// 首先尝试将 msg 转换为 Message（如果实现了 Message 接口）
	if msgMsg, ok := msg.(Message); ok {
		// 如果实现了 Message 接口，应该已经被上面的 switch 处理了
		// 这里不应该到达，但为了安全起见，我们再次尝试分类处理
		return a.Receive(ctx, msgMsg)
	}

	// 尝试找到能处理此消息的 Capacity（作为 capacities.Message）
	var capMsg capacities.Message = msg
	for _, capacity := range a.capabilities {
		if capacity.CanHandle(capMsg) {
			return capacity.Execute(ctx, capMsg)
		}
	}

	// 如果没有找到对应的 Capacity，返回错误
	return fmt.Errorf("no capacity can handle message type %T (legacy message)", msg)
}

// setBehaviorState 设置行为状态
func (a *BaseResourceActor) setBehaviorState(state *BehaviorState) {
	a.behaviorMu.Lock()
	defer a.behaviorMu.Unlock()
	a.behaviorState = state
}

// clearBehaviorState 清除行为状态
func (a *BaseResourceActor) clearBehaviorState() {
	a.behaviorMu.Lock()
	defer a.behaviorMu.Unlock()
	a.behaviorState = nil
}

// GetBehaviorState 获取当前行为状态
func (a *BaseResourceActor) GetBehaviorState() *BehaviorState {
	a.behaviorMu.RLock()
	defer a.behaviorMu.RUnlock()
	if a.behaviorState == nil {
		return nil
	}
	// 返回副本
	state := *a.behaviorState
	return &state
}

// applyEffects 应用效果（状态变化和事件）
func (a *BaseResourceActor) applyEffects(effects []capacities.Effect) {
	for _, effect := range effects {
		// 应用状态变化
		for _, stateChange := range effect.StateChanges {
			setPropMsg := &SetPropertyMessage{
				Name:  stateChange.Property,
				Value: stateChange.ToValue,
			}
			_ = a.Send(setPropMsg)
		}

		// 发射事件
		for _, event := range effect.Events {
			_ = a.eventEmitter.Emit(Event{
				Type:    EventType(event.Type),
				Payload: event.Payload,
			})
		}
	}
}

// Start 启动 Actor（扩展基类方法，启动所有需要订阅的 Capacity 和所有 Bindings）
func (a *BaseResourceActor) Start(ctx context.Context) error {
	// 初始化事件发射器（如果还没有）
	if a.eventEmitter == nil && a.BaseActor.system != nil {
		a.eventEmitter = NewBaseEventEmitter(a.resourceID, a.BaseActor.system)
	}

	// 先调用基类的 Start
	if err := a.BaseActor.Start(ctx); err != nil {
		return err
	}

	// 注意：事件注册应该由子类在初始化时调用 RegisterEvent 方法完成
	// 不再通过 EventRegistry 接口自动注册，因为 BaseResourceActor 已有 RegisterEvent 方法

	// 启动所有需要订阅的 Capacity
	for _, capacity := range a.capabilities {
		// 检查是否有 StartSubscription 方法（可选接口）
		if starter, ok := capacity.(interface {
			StartSubscription(context.Context) error
		}); ok {
			// 设置 Capacity 引用和 context
			if capacitySetter, ok := capacity.(interface {
				SetCapacityRef(capacities.Capacity, context.Context)
			}); ok {
				capacitySetter.SetCapacityRef(capacity, ctx)
			}

			// 设置 Actor 引用
			if setter, ok := capacity.(interface {
				SetActorRef(interface {
					Send(Message) bool
				})
			}); ok {
				setter.SetActorRef(a)
			}

			// 启动订阅
			if err := starter.StartSubscription(ctx); err != nil {
				return fmt.Errorf("failed to start subscription for capacity %s: %w",
					capacity.Name(), err)
			}
		}
	}

	// 启动所有外部绑定
	for bindingType, binding := range a.bindings {
		if err := binding.Start(ctx); err != nil {
			return fmt.Errorf("failed to start binding %s: %w", bindingType, err)
		}
	}

	// 转换状态为在线
	from := a.currentState
	if err := a.transition(StateOnline, "start"); err != nil {
		return fmt.Errorf("failed to transition to online state: %w", err)
	}

	// 发射启动事件
	if a.eventEmitter != nil {
		_ = a.eventEmitter.Emit(Event{
			Type:    EventTypeActorStarted,
			Payload: map[string]interface{}{"from": from},
		})
	}

	return nil
}

// Stop 停止 Actor（扩展基类方法，停止所有 Bindings）
func (a *BaseResourceActor) Stop() error {
	// 转换状态为离线
	from := a.currentState
	if err := a.transition(StateOffline, "stop"); err != nil {
		// 记录错误但不阻止停止
	}

	// 发射停止事件
	if a.eventEmitter != nil {
		_ = a.eventEmitter.Emit(Event{
			Type:    EventTypeActorStopped,
			Payload: map[string]interface{}{"from": from},
		})
	}

	// 停止所有外部绑定
	for bindingType, binding := range a.bindings {
		if err := binding.Stop(); err != nil {
			// 记录错误但不阻止停止
			_ = bindingType
		}
	}

	// 调用基类的 Stop
	return a.BaseActor.Stop()
}

// Suspend 挂起 Actor
func (a *BaseResourceActor) Suspend(reason string) error {
	from := a.currentState
	if err := a.transition(StateSuspended, reason); err != nil {
		return err
	}
	if a.eventEmitter != nil {
		_ = a.eventEmitter.Emit(Event{
			Type:    EventTypeActorSuspended,
			Payload: map[string]interface{}{"from": from, "reason": reason},
		})
	}
	return nil
}

// Resume 恢复 Actor
func (a *BaseResourceActor) Resume(reason string) error {
	from := a.currentState
	if err := a.transition(StateOnline, reason); err != nil {
		return err
	}
	if a.eventEmitter != nil {
		_ = a.eventEmitter.Emit(Event{
			Type:    EventTypeActorResumed,
			Payload: map[string]interface{}{"from": from, "reason": reason},
		})
	}
	return nil
}

// GetEventEmitter 获取事件发射器
func (a *BaseResourceActor) GetEventEmitter() EventEmitter {
	return a.eventEmitter
}

// ============================================================================
// 事件管理方法（参考 Capacity 管理）
// ============================================================================

// RegisterEvent 注册一个事件描述符
func (a *BaseResourceActor) RegisterEvent(eventDesc *EventDescriptor) {
	if eventDesc == nil {
		return
	}
	a.eventsMu.Lock()
	defer a.eventsMu.Unlock()
	a.events[eventDesc.Name] = eventDesc
}

// UnregisterEvent 注销一个事件
func (a *BaseResourceActor) UnregisterEvent(eventName string) {
	a.eventsMu.Lock()
	defer a.eventsMu.Unlock()
	delete(a.events, eventName)
}

// HasEvent 检查是否能发出某种事件
func (a *BaseResourceActor) HasEvent(eventName string) bool {
	a.eventsMu.RLock()
	defer a.eventsMu.RUnlock()
	_, exists := a.events[eventName]
	return exists
}

// GetEvent 获取指定名称的事件描述符
func (a *BaseResourceActor) GetEvent(eventName string) (*EventDescriptor, bool) {
	a.eventsMu.RLock()
	defer a.eventsMu.RUnlock()
	eventDesc, exists := a.events[eventName]
	return eventDesc, exists
}

// ListEvents 返回所有事件名称列表
func (a *BaseResourceActor) ListEvents() []string {
	a.eventsMu.RLock()
	defer a.eventsMu.RUnlock()
	names := make([]string, 0, len(a.events))
	for name := range a.events {
		names = append(names, name)
	}
	return names
}

// ListEventDescriptors 返回所有事件描述符列表
func (a *BaseResourceActor) ListEventDescriptors() []*EventDescriptor {
	a.eventsMu.RLock()
	defer a.eventsMu.RUnlock()
	descriptors := make([]*EventDescriptor, 0, len(a.events))
	for _, desc := range a.events {
		descriptors = append(descriptors, desc)
	}
	return descriptors
}

// CanEmitEvent 检查是否能发出指定类型的事件
func (a *BaseResourceActor) CanEmitEvent(eventType EventType, payload interface{}) bool {
	a.eventsMu.RLock()
	defer a.eventsMu.RUnlock()
	for _, desc := range a.events {
		if desc.CanEmit(eventType, payload) {
			return true
		}
	}
	return false
}

// registerStandardEvents 注册标准生命周期事件
func (a *BaseResourceActor) registerStandardEvents() {
	// 注册标准生命周期事件
	standardEvents := []*EventDescriptor{
		NewEventDescriptor(
			"ActorCreated",
			EventTypeActorCreated,
			reflect.TypeOf(map[string]interface{}{}),
			"Actor 创建事件",
			a.resourceID,
		),
		NewEventDescriptor(
			"ActorStarted",
			EventTypeActorStarted,
			reflect.TypeOf(map[string]interface{}{}),
			"Actor 启动事件",
			a.resourceID,
		),
		NewEventDescriptor(
			"ActorStopped",
			EventTypeActorStopped,
			reflect.TypeOf(map[string]interface{}{}),
			"Actor 停止事件",
			a.resourceID,
		),
		NewEventDescriptor(
			"ActorSuspended",
			EventTypeActorSuspended,
			reflect.TypeOf(map[string]interface{}{}),
			"Actor 挂起事件",
			a.resourceID,
		),
		NewEventDescriptor(
			"ActorResumed",
			EventTypeActorResumed,
			reflect.TypeOf(map[string]interface{}{}),
			"Actor 恢复事件",
			a.resourceID,
		),
		NewEventDescriptor(
			"StateChanged",
			EventTypeStateChanged,
			reflect.TypeOf(map[string]interface{}{}),
			"状态变化事件",
			a.resourceID,
		),
		NewEventDescriptor(
			"CommandCompleted",
			EventTypeCommandCompleted,
			reflect.TypeOf(map[string]interface{}{}),
			"命令完成事件",
			a.resourceID,
		),
		NewEventDescriptor(
			"CommandFailed",
			EventTypeCommandFailed,
			reflect.TypeOf(map[string]interface{}{}),
			"命令失败事件",
			a.resourceID,
		),
	}

	for _, eventDesc := range standardEvents {
		a.RegisterEvent(eventDesc)
	}
}

// GetRef 获取指定 Actor 的 ActorRef
func (a *BaseResourceActor) GetRef(targetID string) (*ActorRef, error) {
	return a.BaseActor.GetRef(targetID)
}

// SendTo 向指定 ID 的 Actor 发送消息
func (a *BaseResourceActor) SendTo(targetID string, msg Message) error {
	return a.BaseActor.SendTo(targetID, msg)
}

// Tell 是 SendTo 的别名
func (a *BaseResourceActor) Tell(targetID string, msg Message) error {
	return a.BaseActor.Tell(targetID, msg)
}
