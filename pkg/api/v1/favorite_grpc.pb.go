// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: favorite.proto

package v1

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

// FavoriteServiceClient is the client API for FavoriteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FavoriteServiceClient interface {
	AddItemToFavorites(ctx context.Context, in *AddItemToFavoritesRequest, opts ...grpc.CallOption) (*AddItemToFavoritesResponse, error)
	RemoveItemFromFavorites(ctx context.Context, in *RemoveItemFromFavoritesRequest, opts ...grpc.CallOption) (*RemoveItemFromFavoritesResponse, error)
	GetFavorites(ctx context.Context, in *GetFavoritesRequest, opts ...grpc.CallOption) (*GetFavoritesResponse, error)
}

type favoriteServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFavoriteServiceClient(cc grpc.ClientConnInterface) FavoriteServiceClient {
	return &favoriteServiceClient{cc}
}

func (c *favoriteServiceClient) AddItemToFavorites(ctx context.Context, in *AddItemToFavoritesRequest, opts ...grpc.CallOption) (*AddItemToFavoritesResponse, error) {
	out := new(AddItemToFavoritesResponse)
	err := c.cc.Invoke(ctx, "/v1.FavoriteService/AddItemToFavorites", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteServiceClient) RemoveItemFromFavorites(ctx context.Context, in *RemoveItemFromFavoritesRequest, opts ...grpc.CallOption) (*RemoveItemFromFavoritesResponse, error) {
	out := new(RemoveItemFromFavoritesResponse)
	err := c.cc.Invoke(ctx, "/v1.FavoriteService/RemoveItemFromFavorites", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteServiceClient) GetFavorites(ctx context.Context, in *GetFavoritesRequest, opts ...grpc.CallOption) (*GetFavoritesResponse, error) {
	out := new(GetFavoritesResponse)
	err := c.cc.Invoke(ctx, "/v1.FavoriteService/GetFavorites", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FavoriteServiceServer is the server API for FavoriteService service.
// All implementations must embed UnimplementedFavoriteServiceServer
// for forward compatibility
type FavoriteServiceServer interface {
	AddItemToFavorites(context.Context, *AddItemToFavoritesRequest) (*AddItemToFavoritesResponse, error)
	RemoveItemFromFavorites(context.Context, *RemoveItemFromFavoritesRequest) (*RemoveItemFromFavoritesResponse, error)
	GetFavorites(context.Context, *GetFavoritesRequest) (*GetFavoritesResponse, error)
	mustEmbedUnimplementedFavoriteServiceServer()
}

// UnimplementedFavoriteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFavoriteServiceServer struct {
}

func (UnimplementedFavoriteServiceServer) AddItemToFavorites(context.Context, *AddItemToFavoritesRequest) (*AddItemToFavoritesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItemToFavorites not implemented")
}
func (UnimplementedFavoriteServiceServer) RemoveItemFromFavorites(context.Context, *RemoveItemFromFavoritesRequest) (*RemoveItemFromFavoritesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveItemFromFavorites not implemented")
}
func (UnimplementedFavoriteServiceServer) GetFavorites(context.Context, *GetFavoritesRequest) (*GetFavoritesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFavorites not implemented")
}
func (UnimplementedFavoriteServiceServer) mustEmbedUnimplementedFavoriteServiceServer() {}

// UnsafeFavoriteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FavoriteServiceServer will
// result in compilation errors.
type UnsafeFavoriteServiceServer interface {
	mustEmbedUnimplementedFavoriteServiceServer()
}

func RegisterFavoriteServiceServer(s grpc.ServiceRegistrar, srv FavoriteServiceServer) {
	s.RegisterService(&FavoriteService_ServiceDesc, srv)
}

func _FavoriteService_AddItemToFavorites_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddItemToFavoritesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServiceServer).AddItemToFavorites(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.FavoriteService/AddItemToFavorites",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServiceServer).AddItemToFavorites(ctx, req.(*AddItemToFavoritesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FavoriteService_RemoveItemFromFavorites_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveItemFromFavoritesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServiceServer).RemoveItemFromFavorites(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.FavoriteService/RemoveItemFromFavorites",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServiceServer).RemoveItemFromFavorites(ctx, req.(*RemoveItemFromFavoritesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FavoriteService_GetFavorites_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFavoritesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServiceServer).GetFavorites(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.FavoriteService/GetFavorites",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServiceServer).GetFavorites(ctx, req.(*GetFavoritesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FavoriteService_ServiceDesc is the grpc.ServiceDesc for FavoriteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FavoriteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.FavoriteService",
	HandlerType: (*FavoriteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddItemToFavorites",
			Handler:    _FavoriteService_AddItemToFavorites_Handler,
		},
		{
			MethodName: "RemoveItemFromFavorites",
			Handler:    _FavoriteService_RemoveItemFromFavorites_Handler,
		},
		{
			MethodName: "GetFavorites",
			Handler:    _FavoriteService_GetFavorites_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "favorite.proto",
}