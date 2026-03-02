package server

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/uos-projects/uos-kernel/actor"
)

// MessageFactory creates a Go Message from decoded fields.
type MessageFactory func(fields map[string]interface{}) (actor.Message, error)

// EventSerializer converts a Go event data to a map for proto Struct serialization.
type EventSerializer func(data interface{}) (map[string]interface{}, error)

// MessageRegistry maps type names to factory/serializer functions.
// Used to bridge between proto Struct and concrete Go types.
type MessageRegistry struct {
	mu          sync.RWMutex
	factories   map[string]MessageFactory
	serializers map[string]EventSerializer
}

// NewMessageRegistry creates an empty registry.
func NewMessageRegistry() *MessageRegistry {
	return &MessageRegistry{
		factories:   make(map[string]MessageFactory),
		serializers: make(map[string]EventSerializer),
	}
}

// RegisterCommand registers a command factory.
// typeName is the string used in proto's _type field (e.g., "OpenBreakerCommand").
func (r *MessageRegistry) RegisterCommand(typeName string, factory MessageFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[typeName] = factory
}

// RegisterEvent registers an event serializer.
// goTypeName is reflect.TypeOf(data).String() (e.g., "*main.DeviceAbnormalEvent").
func (r *MessageRegistry) RegisterEvent(goTypeName string, serializer EventSerializer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.serializers[goTypeName] = serializer
}

// CreateMessage creates a Go Message from a type name and fields map.
func (r *MessageRegistry) CreateMessage(typeName string, fields map[string]interface{}) (actor.Message, error) {
	r.mu.RLock()
	factory, ok := r.factories[typeName]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown message type: %s", typeName)
	}
	return factory(fields)
}

// SerializeEvent converts a Go event data to a map with _type field.
// Falls back to a generic serializer using fmt.Sprintf if no specific serializer is registered.
func (r *MessageRegistry) SerializeEvent(data interface{}) (map[string]interface{}, error) {
	if data == nil {
		return nil, nil
	}

	typeName := reflect.TypeOf(data).String()

	r.mu.RLock()
	serializer, ok := r.serializers[typeName]
	r.mu.RUnlock()

	if ok {
		return serializer(data)
	}

	// Fallback: try map[string]interface{} directly
	if m, ok := data.(map[string]interface{}); ok {
		result := make(map[string]interface{}, len(m)+1)
		for k, v := range m {
			result[k] = v
		}
		return result, nil
	}

	// Generic fallback: stringify
	return map[string]interface{}{
		"_type": typeName,
		"value": fmt.Sprintf("%+v", data),
	}, nil
}
