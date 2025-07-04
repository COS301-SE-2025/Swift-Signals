// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.31.1
// source: simulation.proto

package simulation

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SimulationService_GetSimulationResults_FullMethodName = "/swiftsignals.simulation.SimulationService/GetSimulationResults"
	SimulationService_GetSimulationOutput_FullMethodName  = "/swiftsignals.simulation.SimulationService/GetSimulationOutput"
)

// SimulationServiceClient is the client API for SimulationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SimulationServiceClient interface {
	GetSimulationResults(ctx context.Context, in *SimulationRequest, opts ...grpc.CallOption) (*SimulationResultsResponse, error)
	GetSimulationOutput(ctx context.Context, in *SimulationRequest, opts ...grpc.CallOption) (*SimulationOutputResponse, error)
}

type simulationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSimulationServiceClient(cc grpc.ClientConnInterface) SimulationServiceClient {
	return &simulationServiceClient{cc}
}

func (c *simulationServiceClient) GetSimulationResults(ctx context.Context, in *SimulationRequest, opts ...grpc.CallOption) (*SimulationResultsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SimulationResultsResponse)
	err := c.cc.Invoke(ctx, SimulationService_GetSimulationResults_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simulationServiceClient) GetSimulationOutput(ctx context.Context, in *SimulationRequest, opts ...grpc.CallOption) (*SimulationOutputResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SimulationOutputResponse)
	err := c.cc.Invoke(ctx, SimulationService_GetSimulationOutput_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SimulationServiceServer is the server API for SimulationService service.
// All implementations must embed UnimplementedSimulationServiceServer
// for forward compatibility.
type SimulationServiceServer interface {
	GetSimulationResults(context.Context, *SimulationRequest) (*SimulationResultsResponse, error)
	GetSimulationOutput(context.Context, *SimulationRequest) (*SimulationOutputResponse, error)
	mustEmbedUnimplementedSimulationServiceServer()
}

// UnimplementedSimulationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSimulationServiceServer struct{}

func (UnimplementedSimulationServiceServer) GetSimulationResults(context.Context, *SimulationRequest) (*SimulationResultsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSimulationResults not implemented")
}
func (UnimplementedSimulationServiceServer) GetSimulationOutput(context.Context, *SimulationRequest) (*SimulationOutputResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSimulationOutput not implemented")
}
func (UnimplementedSimulationServiceServer) mustEmbedUnimplementedSimulationServiceServer() {}
func (UnimplementedSimulationServiceServer) testEmbeddedByValue()                           {}

// UnsafeSimulationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SimulationServiceServer will
// result in compilation errors.
type UnsafeSimulationServiceServer interface {
	mustEmbedUnimplementedSimulationServiceServer()
}

func RegisterSimulationServiceServer(s grpc.ServiceRegistrar, srv SimulationServiceServer) {
	// If the following call pancis, it indicates UnimplementedSimulationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SimulationService_ServiceDesc, srv)
}

func _SimulationService_GetSimulationResults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimulationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimulationServiceServer).GetSimulationResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SimulationService_GetSimulationResults_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimulationServiceServer).GetSimulationResults(ctx, req.(*SimulationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimulationService_GetSimulationOutput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimulationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimulationServiceServer).GetSimulationOutput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SimulationService_GetSimulationOutput_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimulationServiceServer).GetSimulationOutput(ctx, req.(*SimulationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SimulationService_ServiceDesc is the grpc.ServiceDesc for SimulationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SimulationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "swiftsignals.simulation.SimulationService",
	HandlerType: (*SimulationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSimulationResults",
			Handler:    _SimulationService_GetSimulationResults_Handler,
		},
		{
			MethodName: "GetSimulationOutput",
			Handler:    _SimulationService_GetSimulationOutput_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "simulation.proto",
}
