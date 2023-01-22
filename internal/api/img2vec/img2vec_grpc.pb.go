// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: img2vec.proto

package img2vec

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// Img2VecClient is the client API for Img2Vec service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type Img2VecClient interface {
	Vectorize(ctx context.Context, in *Image, opts ...grpc.CallOption) (*Vector, error)
}

type img2VecClient struct {
	cc grpc.ClientConnInterface
}

func NewImg2VecClient(cc grpc.ClientConnInterface) Img2VecClient {
	return &img2VecClient{cc}
}

func (c *img2VecClient) Vectorize(ctx context.Context, in *Image, opts ...grpc.CallOption) (*Vector, error) {
	out := new(Vector)
	err := c.cc.Invoke(ctx, "/img2vec.Img2Vec/Detect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Img2VecServer is the server API for Img2Vec service.
// All implementations must embed UnimplementedImg2VecServer
// for forward compatibility
type Img2VecServer interface {
	Vectorize(context.Context, *Image) (*Vector, error)
	mustEmbedUnimplementedImg2VecServer()
}

// UnimplementedImg2VecServer must be embedded to have forward compatible implementations.
type UnimplementedImg2VecServer struct {
}

func (UnimplementedImg2VecServer) Vectorize(context.Context, *Image) (*Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Detect not implemented")
}
func (UnimplementedImg2VecServer) mustEmbedUnimplementedImg2VecServer() {}

// UnsafeImg2VecServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Img2VecServer will
// result in compilation errors.
type UnsafeImg2VecServer interface {
	mustEmbedUnimplementedImg2VecServer()
}

func RegisterImg2VecServer(s grpc.ServiceRegistrar, srv Img2VecServer) {
	s.RegisterService(&Img2Vec_ServiceDesc, srv)
}

func _Img2Vec_Vectorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Image)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Img2VecServer).Vectorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/img2vec.Img2Vec/Detect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Img2VecServer).Vectorize(ctx, req.(*Image))
	}
	return interceptor(ctx, in, info, handler)
}

// Img2Vec_ServiceDesc is the grpc.ServiceDesc for Img2Vec service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Img2Vec_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "img2vec.Img2Vec",
	HandlerType: (*Img2VecServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Detect",
			Handler:    _Img2Vec_Vectorize_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "img2vec.proto",
}
