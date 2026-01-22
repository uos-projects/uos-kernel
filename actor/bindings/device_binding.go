package bindings

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/actor"
)

// DeviceBinding 设备协议绑定示例
// 演示如何实现一个具体的 Binding
type DeviceBinding struct {
	*actor.BaseBinding
	deviceProtocol string // 设备协议类型（如 Modbus, OPC-UA 等）
	deviceAddress  string // 设备地址
}

// NewDeviceBinding 创建设备绑定
func NewDeviceBinding(resourceID, deviceProtocol, deviceAddress string) *DeviceBinding {
	base := actor.NewBaseBinding(actor.BindingTypeDevice, resourceID)
	return &DeviceBinding{
		BaseBinding:    base,
		deviceProtocol: deviceProtocol,
		deviceAddress:  deviceAddress,
	}
}

// Start 启动设备绑定
func (b *DeviceBinding) Start(ctx context.Context) error {
	// 连接设备
	// 启动事件监听
	// 启动命令执行器
	fmt.Printf("Device binding started: %s at %s\n", b.deviceProtocol, b.deviceAddress)
	return nil
}

// Stop 停止设备绑定
func (b *DeviceBinding) Stop() error {
	// 断开设备连接
	// 停止事件监听
	fmt.Printf("Device binding stopped: %s at %s\n", b.deviceProtocol, b.deviceAddress)
	return b.BaseBinding.Stop()
}

// OnExternalEvent 接收外部事件（设备数据到达）
func (b *DeviceBinding) OnExternalEvent(ctx context.Context, event interface{}) error {
	// 将设备事件转换为消息
	// 例如：设备状态变化、测量值更新等
	
	// 示例：转换为 ExternalEventMessage
		msg := &actor.ExternalEventMessage{
			BindingType: actor.BindingTypeDevice,
		Event:       event,
	}
	
	// 发送到 Actor
	return b.SendToActor(msg)
}

// ExecuteExternal 执行外部操作（发送命令到设备）
func (b *DeviceBinding) ExecuteExternal(ctx context.Context, command interface{}) error {
	// 将命令转换为设备协议格式
	// 发送到设备
	fmt.Printf("Executing command on device %s: %v\n", b.deviceAddress, command)
	return nil
}
