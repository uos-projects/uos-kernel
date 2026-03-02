package server

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
)

// structToMap converts a protobuf Struct to map[string]interface{}.
func structToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	result := make(map[string]interface{}, len(s.Fields))
	for k, v := range s.Fields {
		result[k] = valueToInterface(v)
	}
	return result
}

// mapToStruct converts map[string]interface{} to protobuf Struct.
func mapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	if m == nil {
		return nil, nil
	}
	fields := make(map[string]*structpb.Value, len(m))
	for k, v := range m {
		val, err := interfaceToValue(v)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", k, err)
		}
		fields[k] = val
	}
	return &structpb.Struct{Fields: fields}, nil
}

// valueToInterface converts a protobuf Value to a Go value.
func valueToInterface(v *structpb.Value) interface{} {
	if v == nil {
		return nil
	}
	switch kind := v.Kind.(type) {
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_NumberValue:
		return kind.NumberValue
	case *structpb.Value_StringValue:
		return kind.StringValue
	case *structpb.Value_BoolValue:
		return kind.BoolValue
	case *structpb.Value_StructValue:
		return structToMap(kind.StructValue)
	case *structpb.Value_ListValue:
		if kind.ListValue == nil {
			return nil
		}
		list := make([]interface{}, len(kind.ListValue.Values))
		for i, item := range kind.ListValue.Values {
			list[i] = valueToInterface(item)
		}
		return list
	default:
		return nil
	}
}

// interfaceToValue converts a Go value to a protobuf Value.
func interfaceToValue(v interface{}) (*structpb.Value, error) {
	if v == nil {
		return structpb.NewNullValue(), nil
	}
	switch val := v.(type) {
	case bool:
		return structpb.NewBoolValue(val), nil
	case float64:
		return structpb.NewNumberValue(val), nil
	case float32:
		return structpb.NewNumberValue(float64(val)), nil
	case int:
		return structpb.NewNumberValue(float64(val)), nil
	case int32:
		return structpb.NewNumberValue(float64(val)), nil
	case int64:
		return structpb.NewNumberValue(float64(val)), nil
	case string:
		return structpb.NewStringValue(val), nil
	case time.Time:
		return structpb.NewStringValue(val.Format(time.RFC3339)), nil
	case map[string]interface{}:
		s, err := mapToStruct(val)
		if err != nil {
			return nil, err
		}
		if s == nil {
			return structpb.NewNullValue(), nil
		}
		return structpb.NewStructValue(s), nil
	case []interface{}:
		list := make([]*structpb.Value, len(val))
		for i, item := range val {
			v, err := interfaceToValue(item)
			if err != nil {
				return nil, fmt.Errorf("index %d: %w", i, err)
			}
			list[i] = v
		}
		return structpb.NewListValue(&structpb.ListValue{Values: list}), nil
	case []string:
		list := make([]*structpb.Value, len(val))
		for i, s := range val {
			list[i] = structpb.NewStringValue(s)
		}
		return structpb.NewListValue(&structpb.ListValue{Values: list}), nil
	default:
		return structpb.NewStringValue(fmt.Sprintf("%v", val)), nil
	}
}
