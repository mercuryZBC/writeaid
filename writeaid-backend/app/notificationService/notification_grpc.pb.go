// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.3
// source: notification.proto

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
	NotificationService_PushKnowledgeBaseOperate_FullMethodName = "/notificationService.NotificationService/PushKnowledgeBaseOperate"
	NotificationService_PushDocumentOperate_FullMethodName      = "/notificationService.NotificationService/PushDocumentOperate"
	NotificationService_PushComment_FullMethodName              = "/notificationService.NotificationService/PushComment"
	NotificationService_Pull_FullMethodName                     = "/notificationService.NotificationService/Pull"
	NotificationService_PullNotificationCount_FullMethodName    = "/notificationService.NotificationService/PullNotificationCount"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	PushKnowledgeBaseOperate(ctx context.Context, in *PushKnowledgeBaseOperateReq, opts ...grpc.CallOption) (*PushKnowledgeBaseOperateResp, error)
	PushDocumentOperate(ctx context.Context, in *PushDocumentOperateReq, opts ...grpc.CallOption) (*PushDocumentOperateResp, error)
	PushComment(ctx context.Context, in *PushCommentReq, opts ...grpc.CallOption) (*PushCommentResp, error)
	Pull(ctx context.Context, in *PullReq, opts ...grpc.CallOption) (*PullResp, error)
	PullNotificationCount(ctx context.Context, in *PullNotificationCountReq, opts ...grpc.CallOption) (*PullNotificationCountResp, error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) PushKnowledgeBaseOperate(ctx context.Context, in *PushKnowledgeBaseOperateReq, opts ...grpc.CallOption) (*PushKnowledgeBaseOperateResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PushKnowledgeBaseOperateResp)
	err := c.cc.Invoke(ctx, NotificationService_PushKnowledgeBaseOperate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) PushDocumentOperate(ctx context.Context, in *PushDocumentOperateReq, opts ...grpc.CallOption) (*PushDocumentOperateResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PushDocumentOperateResp)
	err := c.cc.Invoke(ctx, NotificationService_PushDocumentOperate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) PushComment(ctx context.Context, in *PushCommentReq, opts ...grpc.CallOption) (*PushCommentResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PushCommentResp)
	err := c.cc.Invoke(ctx, NotificationService_PushComment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) Pull(ctx context.Context, in *PullReq, opts ...grpc.CallOption) (*PullResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PullResp)
	err := c.cc.Invoke(ctx, NotificationService_Pull_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) PullNotificationCount(ctx context.Context, in *PullNotificationCountReq, opts ...grpc.CallOption) (*PullNotificationCountResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PullNotificationCountResp)
	err := c.cc.Invoke(ctx, NotificationService_PullNotificationCount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility.
type NotificationServiceServer interface {
	PushKnowledgeBaseOperate(context.Context, *PushKnowledgeBaseOperateReq) (*PushKnowledgeBaseOperateResp, error)
	PushDocumentOperate(context.Context, *PushDocumentOperateReq) (*PushDocumentOperateResp, error)
	PushComment(context.Context, *PushCommentReq) (*PushCommentResp, error)
	Pull(context.Context, *PullReq) (*PullResp, error)
	PullNotificationCount(context.Context, *PullNotificationCountReq) (*PullNotificationCountResp, error)
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNotificationServiceServer struct{}

func (UnimplementedNotificationServiceServer) PushKnowledgeBaseOperate(context.Context, *PushKnowledgeBaseOperateReq) (*PushKnowledgeBaseOperateResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushKnowledgeBaseOperate not implemented")
}
func (UnimplementedNotificationServiceServer) PushDocumentOperate(context.Context, *PushDocumentOperateReq) (*PushDocumentOperateResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushDocumentOperate not implemented")
}
func (UnimplementedNotificationServiceServer) PushComment(context.Context, *PushCommentReq) (*PushCommentResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushComment not implemented")
}
func (UnimplementedNotificationServiceServer) Pull(context.Context, *PullReq) (*PullResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pull not implemented")
}
func (UnimplementedNotificationServiceServer) PullNotificationCount(context.Context, *PullNotificationCountReq) (*PullNotificationCountResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PullNotificationCount not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}
func (UnimplementedNotificationServiceServer) testEmbeddedByValue()                             {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	// If the following call pancis, it indicates UnimplementedNotificationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_PushKnowledgeBaseOperate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushKnowledgeBaseOperateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).PushKnowledgeBaseOperate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_PushKnowledgeBaseOperate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).PushKnowledgeBaseOperate(ctx, req.(*PushKnowledgeBaseOperateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_PushDocumentOperate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushDocumentOperateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).PushDocumentOperate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_PushDocumentOperate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).PushDocumentOperate(ctx, req.(*PushDocumentOperateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_PushComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushCommentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).PushComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_PushComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).PushComment(ctx, req.(*PushCommentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_Pull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PullReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).Pull(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_Pull_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).Pull(ctx, req.(*PullReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_PullNotificationCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PullNotificationCountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).PullNotificationCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_PullNotificationCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).PullNotificationCount(ctx, req.(*PullNotificationCountReq))
	}
	return interceptor(ctx, in, info, handler)
}

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notificationService.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushKnowledgeBaseOperate",
			Handler:    _NotificationService_PushKnowledgeBaseOperate_Handler,
		},
		{
			MethodName: "PushDocumentOperate",
			Handler:    _NotificationService_PushDocumentOperate_Handler,
		},
		{
			MethodName: "PushComment",
			Handler:    _NotificationService_PushComment_Handler,
		},
		{
			MethodName: "Pull",
			Handler:    _NotificationService_Pull_Handler,
		},
		{
			MethodName: "PullNotificationCount",
			Handler:    _NotificationService_PullNotificationCount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notification.proto",
}
