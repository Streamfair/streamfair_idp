// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: role/rpc_get_role_by_id.proto

package role

import (
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

type GetRoleByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetRoleByIDRequest) Reset() {
	*x = GetRoleByIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_role_rpc_get_role_by_id_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRoleByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRoleByIDRequest) ProtoMessage() {}

func (x *GetRoleByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_role_rpc_get_role_by_id_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRoleByIDRequest.ProtoReflect.Descriptor instead.
func (*GetRoleByIDRequest) Descriptor() ([]byte, []int) {
	return file_role_rpc_get_role_by_id_proto_rawDescGZIP(), []int{0}
}

func (x *GetRoleByIDRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type GetRoleByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role *Role `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *GetRoleByIDResponse) Reset() {
	*x = GetRoleByIDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_role_rpc_get_role_by_id_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRoleByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRoleByIDResponse) ProtoMessage() {}

func (x *GetRoleByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_role_rpc_get_role_by_id_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRoleByIDResponse.ProtoReflect.Descriptor instead.
func (*GetRoleByIDResponse) Descriptor() ([]byte, []int) {
	return file_role_rpc_get_role_by_id_proto_rawDescGZIP(), []int{1}
}

func (x *GetRoleByIDResponse) GetRole() *Role {
	if x != nil {
		return x.Role
	}
	return nil
}

var File_role_rpc_get_role_by_id_proto protoreflect.FileDescriptor

var file_role_rpc_get_role_by_id_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x72, 0x6f, 0x6c, 0x65, 0x2f, 0x72, 0x70, 0x63, 0x5f, 0x67, 0x65, 0x74, 0x5f, 0x72,
	0x6f, 0x6c, 0x65, 0x5f, 0x62, 0x79, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x1a, 0x0f, 0x72, 0x6f, 0x6c, 0x65, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x24, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x42,
	0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x33, 0x0a, 0x13, 0x47, 0x65,
	0x74, 0x52, 0x6f, 0x6c, 0x65, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1c, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x08, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x42,
	0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x66, 0x61, 0x69, 0x72, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x66,
	0x61, 0x69, 0x72, 0x5f, 0x69, 0x64, 0x70, 0x5f, 0x73, 0x76, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x72,
	0x6f, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_role_rpc_get_role_by_id_proto_rawDescOnce sync.Once
	file_role_rpc_get_role_by_id_proto_rawDescData = file_role_rpc_get_role_by_id_proto_rawDesc
)

func file_role_rpc_get_role_by_id_proto_rawDescGZIP() []byte {
	file_role_rpc_get_role_by_id_proto_rawDescOnce.Do(func() {
		file_role_rpc_get_role_by_id_proto_rawDescData = protoimpl.X.CompressGZIP(file_role_rpc_get_role_by_id_proto_rawDescData)
	})
	return file_role_rpc_get_role_by_id_proto_rawDescData
}

var file_role_rpc_get_role_by_id_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_role_rpc_get_role_by_id_proto_goTypes = []interface{}{
	(*GetRoleByIDRequest)(nil),  // 0: pb.GetRoleByIDRequest
	(*GetRoleByIDResponse)(nil), // 1: pb.GetRoleByIDResponse
	(*Role)(nil),                // 2: pb.Role
}
var file_role_rpc_get_role_by_id_proto_depIdxs = []int32{
	2, // 0: pb.GetRoleByIDResponse.role:type_name -> pb.Role
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_role_rpc_get_role_by_id_proto_init() }
func file_role_rpc_get_role_by_id_proto_init() {
	if File_role_rpc_get_role_by_id_proto != nil {
		return
	}
	file_role_role_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_role_rpc_get_role_by_id_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRoleByIDRequest); i {
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
		file_role_rpc_get_role_by_id_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRoleByIDResponse); i {
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
			RawDescriptor: file_role_rpc_get_role_by_id_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_role_rpc_get_role_by_id_proto_goTypes,
		DependencyIndexes: file_role_rpc_get_role_by_id_proto_depIdxs,
		MessageInfos:      file_role_rpc_get_role_by_id_proto_msgTypes,
	}.Build()
	File_role_rpc_get_role_by_id_proto = out.File
	file_role_rpc_get_role_by_id_proto_rawDesc = nil
	file_role_rpc_get_role_by_id_proto_goTypes = nil
	file_role_rpc_get_role_by_id_proto_depIdxs = nil
}