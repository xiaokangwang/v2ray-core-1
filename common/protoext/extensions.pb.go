// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.15.6
// source: common/protoext/extensions.proto

package protoext

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MessageOpt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      []string `protobuf:"bytes,1,rep,name=type,proto3" json:"type,omitempty"`
	ShortName []string `protobuf:"bytes,2,rep,name=short_name,json=shortName,proto3" json:"short_name,omitempty"`
}

func (x *MessageOpt) Reset() {
	*x = MessageOpt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_protoext_extensions_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageOpt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOpt) ProtoMessage() {}

func (x *MessageOpt) ProtoReflect() protoreflect.Message {
	mi := &file_common_protoext_extensions_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageOpt.ProtoReflect.Descriptor instead.
func (*MessageOpt) Descriptor() ([]byte, []int) {
	return file_common_protoext_extensions_proto_rawDescGZIP(), []int{0}
}

func (x *MessageOpt) GetType() []string {
	if x != nil {
		return x.Type
	}
	return nil
}

func (x *MessageOpt) GetShortName() []string {
	if x != nil {
		return x.ShortName
	}
	return nil
}

type FieldOpt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AnyWants          []string `protobuf:"bytes,1,rep,name=any_wants,json=anyWants,proto3" json:"any_wants,omitempty"`
	AllowedValues     []string `protobuf:"bytes,2,rep,name=allowed_values,json=allowedValues,proto3" json:"allowed_values,omitempty"`
	AllowedValueTypes []string `protobuf:"bytes,3,rep,name=allowed_value_types,json=allowedValueTypes,proto3" json:"allowed_value_types,omitempty"`
}

func (x *FieldOpt) Reset() {
	*x = FieldOpt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_protoext_extensions_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldOpt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldOpt) ProtoMessage() {}

func (x *FieldOpt) ProtoReflect() protoreflect.Message {
	mi := &file_common_protoext_extensions_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldOpt.ProtoReflect.Descriptor instead.
func (*FieldOpt) Descriptor() ([]byte, []int) {
	return file_common_protoext_extensions_proto_rawDescGZIP(), []int{1}
}

func (x *FieldOpt) GetAnyWants() []string {
	if x != nil {
		return x.AnyWants
	}
	return nil
}

func (x *FieldOpt) GetAllowedValues() []string {
	if x != nil {
		return x.AllowedValues
	}
	return nil
}

func (x *FieldOpt) GetAllowedValueTypes() []string {
	if x != nil {
		return x.AllowedValueTypes
	}
	return nil
}

var file_common_protoext_extensions_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*MessageOpt)(nil),
		Field:         50000,
		Name:          "v2ray.core.common.protoext.message_opt",
		Tag:           "bytes,50000,opt,name=message_opt",
		Filename:      "common/protoext/extensions.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*FieldOpt)(nil),
		Field:         50000,
		Name:          "v2ray.core.common.protoext.field_opt",
		Tag:           "bytes,50000,opt,name=field_opt",
		Filename:      "common/protoext/extensions.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional v2ray.core.common.protoext.MessageOpt message_opt = 50000;
	E_MessageOpt = &file_common_protoext_extensions_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional v2ray.core.common.protoext.FieldOpt field_opt = 50000;
	E_FieldOpt = &file_common_protoext_extensions_proto_extTypes[1]
)

var File_common_protoext_extensions_proto protoreflect.FileDescriptor

var file_common_protoext_extensions_proto_rawDesc = []byte{
	0x0a, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x65, 0x78,
	0x74, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1a, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x65, 0x78, 0x74, 0x1a, 0x20,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x3f, 0x0a, 0x0a, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x22, 0x7e, 0x0a, 0x08, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x61, 0x6e, 0x79, 0x5f, 0x77, 0x61, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x08, 0x61, 0x6e, 0x79, 0x57, 0x61, 0x6e, 0x74, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x6c,
	0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0d, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x12, 0x2e, 0x0a, 0x13, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x11,
	0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x73, 0x3a, 0x6a, 0x0a, 0x0b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x6f, 0x70, 0x74,
	0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x76, 0x32, 0x72, 0x61,
	0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x65, 0x78, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70,
	0x74, 0x52, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x3a, 0x62, 0x0a,
	0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6f, 0x70, 0x74, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x24, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x65, 0x78, 0x74, 0x2e, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x52, 0x08, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70,
	0x74, 0x42, 0x6f, 0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x65, 0x78, 0x74, 0x50, 0x01, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x76, 0x32, 0x66, 0x6c, 0x79, 0x2f, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2d, 0x63, 0x6f,
	0x72, 0x65, 0x2f, 0x76, 0x34, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x65, 0x78, 0x74, 0xaa, 0x02, 0x1a, 0x56, 0x32, 0x52, 0x61, 0x79, 0x2e, 0x43, 0x6f,
	0x72, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x45,
	0x78, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_protoext_extensions_proto_rawDescOnce sync.Once
	file_common_protoext_extensions_proto_rawDescData = file_common_protoext_extensions_proto_rawDesc
)

func file_common_protoext_extensions_proto_rawDescGZIP() []byte {
	file_common_protoext_extensions_proto_rawDescOnce.Do(func() {
		file_common_protoext_extensions_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_protoext_extensions_proto_rawDescData)
	})
	return file_common_protoext_extensions_proto_rawDescData
}

var file_common_protoext_extensions_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_common_protoext_extensions_proto_goTypes = []interface{}{
	(*MessageOpt)(nil),                  // 0: v2ray.core.common.protoext.MessageOpt
	(*FieldOpt)(nil),                    // 1: v2ray.core.common.protoext.FieldOpt
	(*descriptorpb.MessageOptions)(nil), // 2: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 3: google.protobuf.FieldOptions
}
var file_common_protoext_extensions_proto_depIdxs = []int32{
	2, // 0: v2ray.core.common.protoext.message_opt:extendee -> google.protobuf.MessageOptions
	3, // 1: v2ray.core.common.protoext.field_opt:extendee -> google.protobuf.FieldOptions
	0, // 2: v2ray.core.common.protoext.message_opt:type_name -> v2ray.core.common.protoext.MessageOpt
	1, // 3: v2ray.core.common.protoext.field_opt:type_name -> v2ray.core.common.protoext.FieldOpt
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	2, // [2:4] is the sub-list for extension type_name
	0, // [0:2] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_protoext_extensions_proto_init() }
func file_common_protoext_extensions_proto_init() {
	if File_common_protoext_extensions_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_protoext_extensions_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageOpt); i {
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
		file_common_protoext_extensions_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldOpt); i {
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
			RawDescriptor: file_common_protoext_extensions_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_common_protoext_extensions_proto_goTypes,
		DependencyIndexes: file_common_protoext_extensions_proto_depIdxs,
		MessageInfos:      file_common_protoext_extensions_proto_msgTypes,
		ExtensionInfos:    file_common_protoext_extensions_proto_extTypes,
	}.Build()
	File_common_protoext_extensions_proto = out.File
	file_common_protoext_extensions_proto_rawDesc = nil
	file_common_protoext_extensions_proto_goTypes = nil
	file_common_protoext_extensions_proto_depIdxs = nil
}
