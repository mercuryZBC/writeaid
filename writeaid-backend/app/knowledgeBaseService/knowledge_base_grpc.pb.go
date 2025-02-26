// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.3
// source: knowledge_base.proto

package __

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
	KnowledgeService_CreateKnowledgeBase_FullMethodName    = "/knowledgeBaseService.KnowledgeService/CreateKnowledgeBase"
	KnowledgeService_GetKnowledgeBaseList_FullMethodName   = "/knowledgeBaseService.KnowledgeService/GetKnowledgeBaseList"
	KnowledgeService_GetKnowledgeBaseDetail_FullMethodName = "/knowledgeBaseService.KnowledgeService/GetKnowledgeBaseDetail"
	KnowledgeService_DeleteKnowledgeBase_FullMethodName    = "/knowledgeBaseService.KnowledgeService/DeleteKnowledgeBase"
	KnowledgeService_UpdateKnowledgeBase_FullMethodName    = "/knowledgeBaseService.KnowledgeService/UpdateKnowledgeBase"
)

// KnowledgeServiceClient is the client API for KnowledgeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KnowledgeServiceClient interface {
	CreateKnowledgeBase(ctx context.Context, in *CreateKnowledgeBaseReq, opts ...grpc.CallOption) (*CreateKnowledgeBaseResp, error)
	GetKnowledgeBaseList(ctx context.Context, in *GetKnowledgeBaseListReq, opts ...grpc.CallOption) (*GetKnowledgeBaseListResp, error)
	GetKnowledgeBaseDetail(ctx context.Context, in *GetKnowledgeBaseDetailReq, opts ...grpc.CallOption) (*GetKnowledgeBaseDetailResp, error)
	DeleteKnowledgeBase(ctx context.Context, in *DeleteKnowledgeBaseReq, opts ...grpc.CallOption) (*DeleteKnowledgeBaseResp, error)
	UpdateKnowledgeBase(ctx context.Context, in *UpdateKnowledgeBaseReq, opts ...grpc.CallOption) (*UpdateKnowledgeBaseResp, error)
}

type knowledgeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKnowledgeServiceClient(cc grpc.ClientConnInterface) KnowledgeServiceClient {
	return &knowledgeServiceClient{cc}
}

func (c *knowledgeServiceClient) CreateKnowledgeBase(ctx context.Context, in *CreateKnowledgeBaseReq, opts ...grpc.CallOption) (*CreateKnowledgeBaseResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateKnowledgeBaseResp)
	err := c.cc.Invoke(ctx, KnowledgeService_CreateKnowledgeBase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *knowledgeServiceClient) GetKnowledgeBaseList(ctx context.Context, in *GetKnowledgeBaseListReq, opts ...grpc.CallOption) (*GetKnowledgeBaseListResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetKnowledgeBaseListResp)
	err := c.cc.Invoke(ctx, KnowledgeService_GetKnowledgeBaseList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *knowledgeServiceClient) GetKnowledgeBaseDetail(ctx context.Context, in *GetKnowledgeBaseDetailReq, opts ...grpc.CallOption) (*GetKnowledgeBaseDetailResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetKnowledgeBaseDetailResp)
	err := c.cc.Invoke(ctx, KnowledgeService_GetKnowledgeBaseDetail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *knowledgeServiceClient) DeleteKnowledgeBase(ctx context.Context, in *DeleteKnowledgeBaseReq, opts ...grpc.CallOption) (*DeleteKnowledgeBaseResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteKnowledgeBaseResp)
	err := c.cc.Invoke(ctx, KnowledgeService_DeleteKnowledgeBase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *knowledgeServiceClient) UpdateKnowledgeBase(ctx context.Context, in *UpdateKnowledgeBaseReq, opts ...grpc.CallOption) (*UpdateKnowledgeBaseResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateKnowledgeBaseResp)
	err := c.cc.Invoke(ctx, KnowledgeService_UpdateKnowledgeBase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KnowledgeServiceServer is the server API for KnowledgeService service.
// All implementations must embed UnimplementedKnowledgeServiceServer
// for forward compatibility.
type KnowledgeServiceServer interface {
	CreateKnowledgeBase(context.Context, *CreateKnowledgeBaseReq) (*CreateKnowledgeBaseResp, error)
	GetKnowledgeBaseList(context.Context, *GetKnowledgeBaseListReq) (*GetKnowledgeBaseListResp, error)
	GetKnowledgeBaseDetail(context.Context, *GetKnowledgeBaseDetailReq) (*GetKnowledgeBaseDetailResp, error)
	DeleteKnowledgeBase(context.Context, *DeleteKnowledgeBaseReq) (*DeleteKnowledgeBaseResp, error)
	UpdateKnowledgeBase(context.Context, *UpdateKnowledgeBaseReq) (*UpdateKnowledgeBaseResp, error)
	mustEmbedUnimplementedKnowledgeServiceServer()
}

// UnimplementedKnowledgeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedKnowledgeServiceServer struct{}

func (UnimplementedKnowledgeServiceServer) CreateKnowledgeBase(context.Context, *CreateKnowledgeBaseReq) (*CreateKnowledgeBaseResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateKnowledgeBase not implemented")
}
func (UnimplementedKnowledgeServiceServer) GetKnowledgeBaseList(context.Context, *GetKnowledgeBaseListReq) (*GetKnowledgeBaseListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKnowledgeBaseList not implemented")
}
func (UnimplementedKnowledgeServiceServer) GetKnowledgeBaseDetail(context.Context, *GetKnowledgeBaseDetailReq) (*GetKnowledgeBaseDetailResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKnowledgeBaseDetail not implemented")
}
func (UnimplementedKnowledgeServiceServer) DeleteKnowledgeBase(context.Context, *DeleteKnowledgeBaseReq) (*DeleteKnowledgeBaseResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKnowledgeBase not implemented")
}
func (UnimplementedKnowledgeServiceServer) UpdateKnowledgeBase(context.Context, *UpdateKnowledgeBaseReq) (*UpdateKnowledgeBaseResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateKnowledgeBase not implemented")
}
func (UnimplementedKnowledgeServiceServer) mustEmbedUnimplementedKnowledgeServiceServer() {}
func (UnimplementedKnowledgeServiceServer) testEmbeddedByValue()                          {}

// UnsafeKnowledgeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KnowledgeServiceServer will
// result in compilation errors.
type UnsafeKnowledgeServiceServer interface {
	mustEmbedUnimplementedKnowledgeServiceServer()
}

func RegisterKnowledgeServiceServer(s grpc.ServiceRegistrar, srv KnowledgeServiceServer) {
	// If the following call pancis, it indicates UnimplementedKnowledgeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&KnowledgeService_ServiceDesc, srv)
}

func _KnowledgeService_CreateKnowledgeBase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateKnowledgeBaseReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KnowledgeServiceServer).CreateKnowledgeBase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KnowledgeService_CreateKnowledgeBase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KnowledgeServiceServer).CreateKnowledgeBase(ctx, req.(*CreateKnowledgeBaseReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _KnowledgeService_GetKnowledgeBaseList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKnowledgeBaseListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KnowledgeServiceServer).GetKnowledgeBaseList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KnowledgeService_GetKnowledgeBaseList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KnowledgeServiceServer).GetKnowledgeBaseList(ctx, req.(*GetKnowledgeBaseListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _KnowledgeService_GetKnowledgeBaseDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKnowledgeBaseDetailReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KnowledgeServiceServer).GetKnowledgeBaseDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KnowledgeService_GetKnowledgeBaseDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KnowledgeServiceServer).GetKnowledgeBaseDetail(ctx, req.(*GetKnowledgeBaseDetailReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _KnowledgeService_DeleteKnowledgeBase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteKnowledgeBaseReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KnowledgeServiceServer).DeleteKnowledgeBase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KnowledgeService_DeleteKnowledgeBase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KnowledgeServiceServer).DeleteKnowledgeBase(ctx, req.(*DeleteKnowledgeBaseReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _KnowledgeService_UpdateKnowledgeBase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateKnowledgeBaseReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KnowledgeServiceServer).UpdateKnowledgeBase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KnowledgeService_UpdateKnowledgeBase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KnowledgeServiceServer).UpdateKnowledgeBase(ctx, req.(*UpdateKnowledgeBaseReq))
	}
	return interceptor(ctx, in, info, handler)
}

// KnowledgeService_ServiceDesc is the grpc.ServiceDesc for KnowledgeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KnowledgeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "knowledgeBaseService.KnowledgeService",
	HandlerType: (*KnowledgeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateKnowledgeBase",
			Handler:    _KnowledgeService_CreateKnowledgeBase_Handler,
		},
		{
			MethodName: "GetKnowledgeBaseList",
			Handler:    _KnowledgeService_GetKnowledgeBaseList_Handler,
		},
		{
			MethodName: "GetKnowledgeBaseDetail",
			Handler:    _KnowledgeService_GetKnowledgeBaseDetail_Handler,
		},
		{
			MethodName: "DeleteKnowledgeBase",
			Handler:    _KnowledgeService_DeleteKnowledgeBase_Handler,
		},
		{
			MethodName: "UpdateKnowledgeBase",
			Handler:    _KnowledgeService_UpdateKnowledgeBase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "knowledge_base.proto",
}
