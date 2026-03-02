package server

import (
	"context"
	"fmt"
	"net"

	"github.com/uos-projects/uos-kernel/kernel"
	kernelpb "github.com/uos-projects/uos-kernel/server/gen/uos/kernel/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// KernelServer implements the gRPC KernelService.
type KernelServer struct {
	kernelpb.UnimplementedKernelServiceServer
	kernel   *kernel.Kernel
	registry *MessageRegistry
}

// NewKernelServer creates a gRPC server wrapping the given Kernel.
func NewKernelServer(k *kernel.Kernel, reg *MessageRegistry) *KernelServer {
	return &KernelServer{
		kernel:   k,
		registry: reg,
	}
}

// Open opens a resource and returns a file descriptor.
func (s *KernelServer) Open(ctx context.Context, req *kernelpb.OpenRequest) (*kernelpb.OpenResponse, error) {
	fd, err := s.kernel.Open(req.ResourceType, req.ResourceId, int(req.Flags))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "open failed: %v", err)
	}
	return &kernelpb.OpenResponse{Fd: int32(fd)}, nil
}

// Close closes a resource descriptor.
func (s *KernelServer) Close(ctx context.Context, req *kernelpb.CloseRequest) (*kernelpb.CloseResponse, error) {
	err := s.kernel.Close(kernel.ResourceDescriptor(req.Fd))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "close failed: %v", err)
	}
	return &kernelpb.CloseResponse{}, nil
}

// Read reads resource state.
func (s *KernelServer) Read(ctx context.Context, req *kernelpb.ReadRequest) (*kernelpb.ReadResponse, error) {
	state, err := s.kernel.Read(ctx, kernel.ResourceDescriptor(req.Fd))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "read failed: %v", err)
	}

	props, err := mapToStruct(state.Properties)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encode properties: %v", err)
	}

	return &kernelpb.ReadResponse{
		ResourceId:   state.ResourceID,
		ResourceType: state.ResourceType,
		Capabilities: state.Capabilities,
		Properties:   props,
	}, nil
}

// Write writes property updates to a resource.
func (s *KernelServer) Write(ctx context.Context, req *kernelpb.WriteRequest) (*kernelpb.WriteResponse, error) {
	updates := structToMap(req.Updates)
	err := s.kernel.Write(ctx, kernel.ResourceDescriptor(req.Fd), &kernel.WriteRequest{
		Updates: updates,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "write failed: %v", err)
	}
	return &kernelpb.WriteResponse{}, nil
}

// Stat returns resource metadata.
func (s *KernelServer) Stat(ctx context.Context, req *kernelpb.StatRequest) (*kernelpb.StatResponse, error) {
	st, err := s.kernel.Stat(ctx, kernel.ResourceDescriptor(req.Fd))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "stat failed: %v", err)
	}

	caps := make([]*kernelpb.CapabilityInfo, len(st.Capabilities))
	for i, c := range st.Capabilities {
		caps[i] = &kernelpb.CapabilityInfo{
			Name:        c.Name,
			Operations:  c.Operations,
			Description: c.Description,
		}
	}

	events := make([]*kernelpb.EventInfo, len(st.Events))
	for i, e := range st.Events {
		events[i] = &kernelpb.EventInfo{
			Name:        e.Name,
			EventType:   e.EventType,
			Description: e.Description,
		}
	}

	return &kernelpb.StatResponse{
		ResourceId:   st.ResourceID,
		ResourceType: st.ResourceType,
		Capabilities: caps,
		Events:       events,
	}, nil
}

// Ioctl executes a control command on a resource.
func (s *KernelServer) Ioctl(ctx context.Context, req *kernelpb.IoctlRequest) (*kernelpb.IoctlResponse, error) {
	fd := kernel.ResourceDescriptor(req.Fd)
	request := int(req.Request)

	var argp interface{}

	// For CMD_EXECUTE_CAPACITY, build the Go message from _type + fields
	if request == 0x1003 { // CMD_EXECUTE_CAPACITY
		fields := structToMap(req.Args)
		if fields == nil {
			return nil, status.Errorf(codes.InvalidArgument, "args required for CMD_EXECUTE_CAPACITY")
		}

		typeName, ok := fields["_type"].(string)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "_type field required in args")
		}
		delete(fields, "_type")

		// Extract optional capacity name
		var capacityName string
		if cap, ok := fields["capacity"].(string); ok {
			capacityName = cap
			delete(fields, "capacity")
		}

		msg, err := s.registry.CreateMessage(typeName, fields)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "create message: %v", err)
		}

		arg := map[string]interface{}{"message": msg}
		if capacityName != "" {
			arg["capacity"] = capacityName
		}
		argp = arg
	} else if req.Args != nil {
		argp = structToMap(req.Args)
	}

	result, err := s.kernel.Ioctl(ctx, fd, request, argp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ioctl failed: %v", err)
	}

	resp := &kernelpb.IoctlResponse{}
	if result != nil {
		resultMap, err := toResultMap(result)
		if err == nil && resultMap != nil {
			resp.Result, _ = mapToStruct(resultMap)
		}
	}

	return resp, nil
}

// Watch streams resource events to the client.
func (s *KernelServer) Watch(req *kernelpb.WatchRequest, stream kernelpb.KernelService_WatchServer) error {
	fd := kernel.ResourceDescriptor(req.Fd)

	// Convert proto EventType enums to kernel EventType strings
	var eventTypes []kernel.EventType
	for _, et := range req.EventTypes {
		if kt, ok := protoToKernelEventType(et); ok {
			eventTypes = append(eventTypes, kt)
		}
	}

	ch, cancel, err := s.kernel.Watch(fd, eventTypes)
	if err != nil {
		return status.Errorf(codes.Internal, "watch failed: %v", err)
	}
	defer cancel()

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-ch:
			if !ok {
				return nil
			}
			protoEvent, err := s.convertEvent(event)
			if err != nil {
				continue // skip events that fail to serialize
			}
			if err := stream.Send(protoEvent); err != nil {
				return err
			}
		}
	}
}

// convertEvent converts a kernel.Event to proto WatchEvent.
func (s *KernelServer) convertEvent(event kernel.Event) (*kernelpb.WatchEvent, error) {
	protoType := kernelToProtoEventType(event.Type)

	var data map[string]interface{}
	var err error

	switch event.Type {
	case kernel.EventCustom:
		// Custom events: use registry to serialize typed event data
		data, err = s.registry.SerializeEvent(event.Data)
		if err != nil {
			return nil, fmt.Errorf("serialize event: %w", err)
		}
	default:
		// Property/state/capability changes: data is already a map
		if m, ok := event.Data.(map[string]interface{}); ok {
			data = m
		} else {
			data = map[string]interface{}{"value": fmt.Sprintf("%v", event.Data)}
		}
	}

	dataStruct, err := mapToStruct(data)
	if err != nil {
		return nil, fmt.Errorf("encode event data: %w", err)
	}

	return &kernelpb.WatchEvent{
		Type:       protoType,
		ResourceId: event.ResourceID,
		Data:       dataStruct,
	}, nil
}

// Serve starts the gRPC server on the given address.
func Serve(addr string, k *kernel.Kernel, reg *MessageRegistry) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("listen: %w", err)
	}
	grpcServer := grpc.NewServer()
	kernelpb.RegisterKernelServiceServer(grpcServer, NewKernelServer(k, reg))
	return grpcServer, lis, nil
}

// ── Helpers ──

func protoToKernelEventType(et kernelpb.EventType) (kernel.EventType, bool) {
	switch et {
	case kernelpb.EventType_EVENT_TYPE_STATE_CHANGE:
		return kernel.EventStateChange, true
	case kernelpb.EventType_EVENT_TYPE_ATTRIBUTE_CHANGE:
		return kernel.EventAttributeChange, true
	case kernelpb.EventType_EVENT_TYPE_CAPABILITY_CHANGE:
		return kernel.EventCapabilityChange, true
	case kernelpb.EventType_EVENT_TYPE_CUSTOM:
		return kernel.EventCustom, true
	default:
		return "", false
	}
}

func kernelToProtoEventType(et kernel.EventType) kernelpb.EventType {
	switch et {
	case kernel.EventStateChange:
		return kernelpb.EventType_EVENT_TYPE_STATE_CHANGE
	case kernel.EventAttributeChange:
		return kernelpb.EventType_EVENT_TYPE_ATTRIBUTE_CHANGE
	case kernel.EventCapabilityChange:
		return kernelpb.EventType_EVENT_TYPE_CAPABILITY_CHANGE
	case kernel.EventCustom:
		return kernelpb.EventType_EVENT_TYPE_CUSTOM
	default:
		return kernelpb.EventType_EVENT_TYPE_UNSPECIFIED
	}
}

// toResultMap converts an ioctl result to a map for Struct serialization.
func toResultMap(result interface{}) (map[string]interface{}, error) {
	if result == nil {
		return nil, nil
	}
	switch v := result.(type) {
	case map[string]interface{}:
		return v, nil
	case []string:
		return map[string]interface{}{"items": v}, nil
	case string:
		return map[string]interface{}{"value": v}, nil
	default:
		return map[string]interface{}{"value": fmt.Sprintf("%v", v)}, nil
	}
}
