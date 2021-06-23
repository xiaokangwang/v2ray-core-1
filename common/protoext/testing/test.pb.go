// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.6
// source: common/protoext/testing/test.proto

package testing

import (
	_ "github.com/v2fly/v2ray-core/v4/common/protoext"
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

type TestingMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TestField string `protobuf:"bytes,1,opt,name=test_field,json=testField,proto3" json:"test_field,omitempty"`
}

func (x *TestingMessage) Reset() {
	*x = TestingMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_protoext_testing_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestingMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestingMessage) ProtoMessage() {}

func (x *TestingMessage) ProtoReflect() protoreflect.Message {
	mi := &file_common_protoext_testing_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestingMessage.ProtoReflect.Descriptor instead.
func (*TestingMessage) Descriptor() ([]byte, []int) {
	return file_common_protoext_testing_test_proto_rawDescGZIP(), []int{0}
}

func (x *TestingMessage) GetTestField() string {
	if x != nil {
		return x.TestField
	}
	return ""
}

var File_common_protoext_testing_test_proto protoreflect.FileDescriptor

var file_common_protoext_testing_test_proto_rawDesc = []byte{
	0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x65, 0x78,
	0x74, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x22, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x65, 0x78, 0x74,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x1a, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x65, 0x78, 0x74, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5d, 0x0a, 0x0e, 0x54, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x34, 0x0a, 0x0a,
	0x74, 0x65, 0x73, 0x74, 0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x15, 0x82, 0xb5, 0x18, 0x06, 0x12, 0x04, 0x74, 0x65, 0x73, 0x74, 0x82, 0xb5, 0x18, 0x07,
	0x12, 0x05, 0x74, 0x65, 0x73, 0x74, 0x32, 0x52, 0x09, 0x74, 0x65, 0x73, 0x74, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x3a, 0x15, 0x82, 0xb5, 0x18, 0x06, 0x0a, 0x04, 0x64, 0x65, 0x6d, 0x6f, 0x82, 0xb5,
	0x18, 0x07, 0x0a, 0x05, 0x64, 0x65, 0x6d, 0x6f, 0x32, 0x42, 0x84, 0x01, 0x0a, 0x26, 0x63, 0x6f,
	0x6d, 0x2e, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x65, 0x78, 0x74, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x69, 0x6e, 0x67, 0x50, 0x01, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x76, 0x32, 0x66, 0x6c, 0x79, 0x2f, 0x76, 0x32, 0x72, 0x61, 0x79, 0x2d, 0x63,
	0x6f, 0x72, 0x65, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x65, 0x78, 0x74, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0xaa, 0x02, 0x22, 0x56, 0x32,
	0x52, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x45, 0x78, 0x74, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_protoext_testing_test_proto_rawDescOnce sync.Once
	file_common_protoext_testing_test_proto_rawDescData = file_common_protoext_testing_test_proto_rawDesc
)

func file_common_protoext_testing_test_proto_rawDescGZIP() []byte {
	file_common_protoext_testing_test_proto_rawDescOnce.Do(func() {
		file_common_protoext_testing_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_protoext_testing_test_proto_rawDescData)
	})
	return file_common_protoext_testing_test_proto_rawDescData
}

var file_common_protoext_testing_test_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_protoext_testing_test_proto_goTypes = []interface{}{
	(*TestingMessage)(nil), // 0: v2ray.core.common.protoext.testing.TestingMessage
}
var file_common_protoext_testing_test_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_protoext_testing_test_proto_init() }
func file_common_protoext_testing_test_proto_init() {
	if File_common_protoext_testing_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_protoext_testing_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestingMessage); i {
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
			RawDescriptor: file_common_protoext_testing_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_protoext_testing_test_proto_goTypes,
		DependencyIndexes: file_common_protoext_testing_test_proto_depIdxs,
		MessageInfos:      file_common_protoext_testing_test_proto_msgTypes,
	}.Build()
	File_common_protoext_testing_test_proto = out.File
	file_common_protoext_testing_test_proto_rawDesc = nil
	file_common_protoext_testing_test_proto_goTypes = nil
	file_common_protoext_testing_test_proto_depIdxs = nil
}
