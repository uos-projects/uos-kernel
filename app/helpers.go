package app

import (
	"context"
	"fmt"

	"github.com/uos-projects/uos-kernel/kernel"
)

// WaitForProperty 等待资源属性达到预期值
// 先检查当前值，若已匹配则立即返回；否则通过 Watch 等待属性变更。
func WaitForProperty(ctx context.Context, k *kernel.Kernel, fd kernel.ResourceDescriptor, prop string, expected interface{}) error {
	// 先检查当前值
	state, err := k.Read(ctx, fd)
	if err == nil && state.Properties != nil {
		if val, ok := state.Properties[prop]; ok && val == expected {
			return nil
		}
	}

	// 通过 Watch 等待变更
	ch, cancel, err := k.Watch(fd, []kernel.EventType{kernel.EventAttributeChange})
	if err != nil {
		return fmt.Errorf("watch failed: %w", err)
	}
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-ch:
			if !ok {
				return fmt.Errorf("watch channel closed")
			}
			// 检查变更的属性是否匹配
			if data, ok := event.Data.(map[string]interface{}); ok {
				if data["name"] == prop && data["value"] == expected {
					return nil
				}
			}
		}
	}
}

// WaitForEvent 等待匹配条件的事件到达
func WaitForEvent(ctx context.Context, k *kernel.Kernel, fd kernel.ResourceDescriptor, match func(kernel.Event) bool) (kernel.Event, error) {
	ch, cancel, err := k.Watch(fd, []kernel.EventType{kernel.EventCustom})
	if err != nil {
		return kernel.Event{}, fmt.Errorf("watch failed: %w", err)
	}
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return kernel.Event{}, ctx.Err()
		case event, ok := <-ch:
			if !ok {
				return kernel.Event{}, fmt.Errorf("watch channel closed")
			}
			if match(event) {
				return event, nil
			}
		}
	}
}
