// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: favorite.proto

package v1

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "github.com/mwitkow/go-proto-validators"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AddItemToFavoritesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ListName  string `protobuf:"bytes,1,opt,name=list_name,json=listName,proto3" json:"list_name,omitempty"`
	ProductId int32  `protobuf:"varint,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
}

func (x *AddItemToFavoritesRequest) Reset() {
	*x = AddItemToFavoritesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddItemToFavoritesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddItemToFavoritesRequest) ProtoMessage() {}

func (x *AddItemToFavoritesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddItemToFavoritesRequest.ProtoReflect.Descriptor instead.
func (*AddItemToFavoritesRequest) Descriptor() ([]byte, []int) {
	return file_favorite_proto_rawDescGZIP(), []int{0}
}

func (x *AddItemToFavoritesRequest) GetListName() string {
	if x != nil {
		return x.ListName
	}
	return ""
}

func (x *AddItemToFavoritesRequest) GetProductId() int32 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

type AddItemToFavoritesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *AddItemToFavoritesResponse) Reset() {
	*x = AddItemToFavoritesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddItemToFavoritesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddItemToFavoritesResponse) ProtoMessage() {}

func (x *AddItemToFavoritesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddItemToFavoritesResponse.ProtoReflect.Descriptor instead.
func (*AddItemToFavoritesResponse) Descriptor() ([]byte, []int) {
	return file_favorite_proto_rawDescGZIP(), []int{1}
}

func (x *AddItemToFavoritesResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type RemoveItemFromFavoritesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ListName  string `protobuf:"bytes,1,opt,name=list_name,json=listName,proto3" json:"list_name,omitempty"`
	ProductId int32  `protobuf:"varint,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
}

func (x *RemoveItemFromFavoritesRequest) Reset() {
	*x = RemoveItemFromFavoritesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveItemFromFavoritesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveItemFromFavoritesRequest) ProtoMessage() {}

func (x *RemoveItemFromFavoritesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveItemFromFavoritesRequest.ProtoReflect.Descriptor instead.
func (*RemoveItemFromFavoritesRequest) Descriptor() ([]byte, []int) {
	return file_favorite_proto_rawDescGZIP(), []int{2}
}

func (x *RemoveItemFromFavoritesRequest) GetListName() string {
	if x != nil {
		return x.ListName
	}
	return ""
}

func (x *RemoveItemFromFavoritesRequest) GetProductId() int32 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

type RemoveItemFromFavoritesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *RemoveItemFromFavoritesResponse) Reset() {
	*x = RemoveItemFromFavoritesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveItemFromFavoritesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveItemFromFavoritesResponse) ProtoMessage() {}

func (x *RemoveItemFromFavoritesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveItemFromFavoritesResponse.ProtoReflect.Descriptor instead.
func (*RemoveItemFromFavoritesResponse) Descriptor() ([]byte, []int) {
	return file_favorite_proto_rawDescGZIP(), []int{3}
}

func (x *RemoveItemFromFavoritesResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type GetFavoritesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ListName string `protobuf:"bytes,1,opt,name=list_name,json=listName,proto3" json:"list_name,omitempty"`
}

func (x *GetFavoritesRequest) Reset() {
	*x = GetFavoritesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFavoritesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFavoritesRequest) ProtoMessage() {}

func (x *GetFavoritesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFavoritesRequest.ProtoReflect.Descriptor instead.
func (*GetFavoritesRequest) Descriptor() ([]byte, []int) {
	return file_favorite_proto_rawDescGZIP(), []int{4}
}

func (x *GetFavoritesRequest) GetListName() string {
	if x != nil {
		return x.ListName
	}
	return ""
}

type GetFavoritesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Products []*Product `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
}

func (x *GetFavoritesResponse) Reset() {
	*x = GetFavoritesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_favorite_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFavoritesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFavoritesResponse) ProtoMessage() {}

func (x *GetFavoritesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_favorite_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFavoritesResponse.ProtoReflect.Descriptor instead.
func (*GetFavoritesResponse) Descriptor() ([]byte, []int) {
	return file_favorite_proto_rawDescGZIP(), []int{5}
}

func (x *GetFavoritesResponse) GetProducts() []*Product {
	if x != nil {
		return x.Products
	}
	return nil
}

var File_favorite_proto protoreflect.FileDescriptor

var file_favorite_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x2c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x73,
	0x77, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x77, 0x69,
	0x74, 0x6b, 0x6f, 0x77, 0x2f, 0x67, 0x6f, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2d, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x72, 0x0a, 0x19, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65,
	0x6d, 0x54, 0x6f, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x09, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x11, 0xe2, 0xdf, 0x1f, 0x0d, 0x0a, 0x0b, 0x5e, 0x66,
	0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x24, 0x52, 0x08, 0x6c, 0x69, 0x73, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x06, 0xe2, 0xdf, 0x1f, 0x02, 0x10, 0x00, 0x52,
	0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x1a, 0x41, 0x64,
	0x64, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x6f, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x22, 0x77, 0x0a, 0x1e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x49, 0x74, 0x65, 0x6d,
	0x46, 0x72, 0x6f, 0x6d, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x09, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x11, 0xe2, 0xdf, 0x1f, 0x0d, 0x0a, 0x0b, 0x5e,
	0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x24, 0x52, 0x08, 0x6c, 0x69, 0x73, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x06, 0xe2, 0xdf, 0x1f, 0x02, 0x10, 0x00,
	0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x22, 0x3b, 0x0a, 0x1f, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x45, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x46,
	0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2e, 0x0a, 0x09, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x11, 0xe2, 0xdf, 0x1f, 0x0d, 0x0a, 0x0b, 0x5e, 0x66, 0x61, 0x76, 0x6f, 0x72,
	0x69, 0x74, 0x65, 0x73, 0x24, 0x52, 0x08, 0x6c, 0x69, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x3f, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73,
	0x32, 0x84, 0x03, 0x0a, 0x0f, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x70, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x54,
	0x6f, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x12, 0x1d, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x6f, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x64, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x6f, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x15, 0x22, 0x10, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x72,
	0x69, 0x74, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x95, 0x01, 0x0a, 0x17, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x49, 0x74, 0x65, 0x6d, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x73, 0x12, 0x22, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x49, 0x74,
	0x65, 0x6d, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x46, 0x72, 0x6f, 0x6d, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x31, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x2b, 0x2a, 0x29, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x7b, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x7d, 0x2f, 0x7b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x67,
	0x0a, 0x0c, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x12, 0x17,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x12, 0x1c, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2f, 0x7b, 0x6c, 0x69, 0x73,
	0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x7d, 0x42, 0xec, 0x01, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x76,
	0x31, 0x92, 0x41, 0xe1, 0x01, 0x12, 0x1a, 0x0a, 0x13, 0x44, 0x69, 0x67, 0x69, 0x76, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x20, 0x41, 0x75, 0x74, 0x68, 0x20, 0x41, 0x50, 0x49, 0x32, 0x03, 0x31, 0x2e,
	0x30, 0x2a, 0x01, 0x02, 0x32, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x52, 0x3b, 0x0a, 0x03, 0x34, 0x30, 0x34, 0x12,
	0x34, 0x0a, 0x2a, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x65, 0x64, 0x20, 0x77, 0x68, 0x65, 0x6e,
	0x20, 0x74, 0x68, 0x65, 0x20, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x20, 0x64, 0x6f,
	0x65, 0x73, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x65, 0x78, 0x69, 0x73, 0x74, 0x2e, 0x12, 0x06, 0x0a,
	0x04, 0x9a, 0x02, 0x01, 0x07, 0x72, 0x5f, 0x0a, 0x1c, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x20, 0x61, 0x70, 0x69, 0x73, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x64, 0x69, 0x67, 0x69, 0x76,
	0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x3f, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x2d, 0x70, 0x72,
	0x6f, 0x67, 0x72, 0x61, 0x6d, 0x6d, 0x69, 0x6e, 0x67, 0x2d, 0x66, 0x61, 0x6c, 0x6c, 0x2d, 0x32,
	0x30, 0x32, 0x32, 0x2f, 0x64, 0x69, 0x67, 0x69, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x2d, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_favorite_proto_rawDescOnce sync.Once
	file_favorite_proto_rawDescData = file_favorite_proto_rawDesc
)

func file_favorite_proto_rawDescGZIP() []byte {
	file_favorite_proto_rawDescOnce.Do(func() {
		file_favorite_proto_rawDescData = protoimpl.X.CompressGZIP(file_favorite_proto_rawDescData)
	})
	return file_favorite_proto_rawDescData
}

var file_favorite_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_favorite_proto_goTypes = []interface{}{
	(*AddItemToFavoritesRequest)(nil),       // 0: v1.AddItemToFavoritesRequest
	(*AddItemToFavoritesResponse)(nil),      // 1: v1.AddItemToFavoritesResponse
	(*RemoveItemFromFavoritesRequest)(nil),  // 2: v1.RemoveItemFromFavoritesRequest
	(*RemoveItemFromFavoritesResponse)(nil), // 3: v1.RemoveItemFromFavoritesResponse
	(*GetFavoritesRequest)(nil),             // 4: v1.GetFavoritesRequest
	(*GetFavoritesResponse)(nil),            // 5: v1.GetFavoritesResponse
	(*Product)(nil),                         // 6: v1.Product
}
var file_favorite_proto_depIdxs = []int32{
	6, // 0: v1.GetFavoritesResponse.products:type_name -> v1.Product
	0, // 1: v1.FavoriteService.AddItemToFavorites:input_type -> v1.AddItemToFavoritesRequest
	2, // 2: v1.FavoriteService.RemoveItemFromFavorites:input_type -> v1.RemoveItemFromFavoritesRequest
	4, // 3: v1.FavoriteService.GetFavorites:input_type -> v1.GetFavoritesRequest
	1, // 4: v1.FavoriteService.AddItemToFavorites:output_type -> v1.AddItemToFavoritesResponse
	3, // 5: v1.FavoriteService.RemoveItemFromFavorites:output_type -> v1.RemoveItemFromFavoritesResponse
	5, // 6: v1.FavoriteService.GetFavorites:output_type -> v1.GetFavoritesResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_favorite_proto_init() }
func file_favorite_proto_init() {
	if File_favorite_proto != nil {
		return
	}
	file_search_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_favorite_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddItemToFavoritesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_favorite_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddItemToFavoritesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_favorite_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveItemFromFavoritesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_favorite_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveItemFromFavoritesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_favorite_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFavoritesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_favorite_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFavoritesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_favorite_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_favorite_proto_goTypes,
		DependencyIndexes: file_favorite_proto_depIdxs,
		MessageInfos:      file_favorite_proto_msgTypes,
	}.Build()
	File_favorite_proto = out.File
	file_favorite_proto_rawDesc = nil
	file_favorite_proto_goTypes = nil
	file_favorite_proto_depIdxs = nil
}
