package actors

import (
	"context"
	"testing"
	"time"

	cimcapacities "github.com/uos-projects/uos-kernel/actors/cim/capacities"
	"github.com/uos-projects/uos-kernel/actors/cim"
)

func TestBaseResourceActor(t *testing.T) {
	ctx := context.Background()
	actor := NewBaseResourceActor("BE-Line_2", "ACLineSegment", nil)

	if actor.ResourceID() != "BE-Line_2" {
		t.Errorf("Expected ResourceID 'BE-Line_2', got '%s'", actor.ResourceID())
	}

	if actor.ResourceType() != "ACLineSegment" {
		t.Errorf("Expected ResourceType 'ACLineSegment', got '%s'", actor.ResourceType())
	}

	// 添加能力（使用 CIM 特定的 Capacity 进行测试）
	capacity := cimcapacities.NewAccumulatorResetCapacity("ACC_RESET_1")
	actor.AddCapacity(capacity)

	if !actor.HasCapacity("AccumulatorResetCapacity") {
		t.Error("Actor should have AccumulatorResetCapacity")
	}

	retrieved, exists := actor.GetCapacity("AccumulatorResetCapacity")
	if !exists {
		t.Error("Should be able to get AccumulatorResetCapacity")
	}
	if retrieved != capacity {
		t.Error("Retrieved capacity should be the same instance")
	}

	// 测试能力列表
	capabilities := actor.ListCapabilities()
	if len(capabilities) != 1 {
		t.Errorf("Expected 1 capability, got %d", len(capabilities))
	}

	// 启动 Actor
	if err := actor.Start(ctx); err != nil {
		t.Fatalf("Failed to start actor: %v", err)
	}

	// 发送消息（CIM 特定的消息类型）
	msg := &cimcapacities.AccumulatorResetMessage{Value: 100}
	actor.Send(msg)

	// 等待处理
	time.Sleep(50 * time.Millisecond)

	// 移除能力
	actor.RemoveCapacity("AccumulatorResetCapacity")
	if actor.HasCapacity("AccumulatorResetCapacity") {
		t.Error("Actor should not have AccumulatorResetCapacity after removal")
	}

	// 停止
	if err := actor.Stop(); err != nil {
		t.Fatalf("Failed to stop actor: %v", err)
	}
}

func TestBaseResourceActorReceive(t *testing.T) {
	ctx := context.Background()
	actor := NewBaseResourceActor("BE-G4", "SynchronousMachine", nil)

	// 添加能力（使用 CIM 特定的 Capacity）
	capacity := cimcapacities.NewRaiseLowerCommandCapacity("CMD_RL_1")
	actor.AddCapacity(capacity)

	// 测试消息路由（CIM 特定的消息类型）
	msg := &cimcapacities.RaiseLowerCommandMessage{Delta: 10.0}
	if err := actor.Receive(ctx, msg); err != nil {
		t.Errorf("Failed to receive message: %v", err)
	}

	// 测试不支持的消息类型（使用简单的测试消息）
	type testMessage struct {
		Value string
	}
	unsupportedMsg := &testMessage{Value: "test"}
	err := actor.Receive(ctx, unsupportedMsg)
	if err == nil {
		t.Error("Expected error for unsupported message type")
	}
}

func TestCapacityFactory(t *testing.T) {
	// 注意：CapacityFactory 现在是 CIM 特定的，应该使用 cim.NewCapacityFactory()
	// 这里保留测试以验证 CIM Capacity 工厂功能
	factory := cim.NewCapacityFactory()

	// 测试创建各种 Capacity（CIM 特定的）
	testCases := []struct {
		controlType  string
		controlID    string
		expectedName string
	}{
		{"AccumulatorReset", "ACC_RESET_1", "AccumulatorResetCapacity"},
		{"Command", "CMD_1", "CommandCapacity"},
		{"RaiseLowerCommand", "CMD_RL_1", "RaiseLowerCommandCapacity"},
		{"SetPoint", "SET_PNT_1", "SetPointCapacity"},
	}

	for _, tc := range testCases {
		capacity, err := factory.CreateCapacity(tc.controlType, tc.controlID)
		if err != nil {
			t.Errorf("Failed to create capacity for %s: %v", tc.controlType, err)
			continue
		}

		if capacity.Name() != tc.expectedName {
			t.Errorf("Expected name %s, got %s", tc.expectedName, capacity.Name())
		}

		if capacity.ResourceID() != tc.controlID {
			t.Errorf("Expected ResourceID %s, got %s", tc.controlID, capacity.ResourceID())
		}
	}

	// 测试未知类型
	_, err := factory.CreateCapacity("UnknownType", "ID")
	if err == nil {
		t.Error("Expected error for unknown control type")
	}
}

func TestCapacityCanHandle(t *testing.T) {
	// 测试 AccumulatorResetCapacity（CIM 特定的）
	accCapacity := cimcapacities.NewAccumulatorResetCapacity("ACC_RESET_1")
	if !accCapacity.CanHandle(&cimcapacities.AccumulatorResetMessage{Value: 100}) {
		t.Error("AccumulatorResetCapacity should handle AccumulatorResetMessage")
	}
	if accCapacity.CanHandle(&cimcapacities.CommandMessage{}) {
		t.Error("AccumulatorResetCapacity should not handle CommandMessage")
	}

	// 测试 CommandCapacity（CIM 特定的）
	cmdCapacity := cimcapacities.NewCommandCapacity("CMD_1")
	if !cmdCapacity.CanHandle(&cimcapacities.CommandMessage{Command: "open"}) {
		t.Error("CommandCapacity should handle CommandMessage")
	}
	if cmdCapacity.CanHandle(&cimcapacities.AccumulatorResetMessage{}) {
		t.Error("CommandCapacity should not handle AccumulatorResetMessage")
	}
}
