package sendas

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Method int32

const (
	Method_FAIL         Method = 0
	Method_RAW          Method = 1
	Method_RAW_ETHERNET Method = 2
)

// Enum value maps for Method.
var (
	Method_name = map[int32]string{
		0: "FAIL",
		1: "RAW",
		2: "RAW_ETHERNET",
	}
	Method_value = map[string]int32{
		"FAIL":         0,
		"RAW":          1,
		"RAW_ETHERNET": 2,
	}
)

func (x Method) Enum() *Method {
	p := new(Method)
	*p = x
	return p
}

func (x Method) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Method) Descriptor() protoreflect.EnumDescriptor {
	return file_common_sendas_sendas_proto_enumTypes[0].Descriptor()
}

func (Method) Type() protoreflect.EnumType {
	return &file_common_sendas_sendas_proto_enumTypes[0]
}

func (x Method) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Method.Descriptor instead.
func (Method) EnumDescriptor() ([]byte, []int) {
	return file_common_sendas_sendas_proto_rawDescGZIP(), []int{0}
}

type Config struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Interface     string                 `protobuf:"bytes,1,opt,name=interface,proto3" json:"interface,omitempty"`
	Method        Method                 `protobuf:"varint,2,opt,name=method,proto3,enum=v2ray.core.common.sendas.Method" json:"method,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Config) Reset() {
	*x = Config{}
	mi := &file_common_sendas_sendas_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_common_sendas_sendas_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_common_sendas_sendas_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetInterface() string {
	if x != nil {
		return x.Interface
	}
	return ""
}

func (x *Config) GetMethod() Method {
	if x != nil {
		return x.Method
	}
	return Method_FAIL
}

var File_common_sendas_sendas_proto protoreflect.FileDescriptor

const file_common_sendas_sendas_proto_rawDesc = "" +
	"\n" +
	"\x1acommon/sendas/sendas.proto\x12\x18v2ray.core.common.sendas\"`\n" +
	"\x06Config\x12\x1c\n" +
	"\tinterface\x18\x01 \x01(\tR\tinterface\x128\n" +
	"\x06method\x18\x02 \x01(\x0e2 .v2ray.core.common.sendas.MethodR\x06method*-\n" +
	"\x06Method\x12\b\n" +
	"\x04FAIL\x10\x00\x12\a\n" +
	"\x03RAW\x10\x01\x12\x10\n" +
	"\fRAW_ETHERNET\x10\x02Bi\n" +
	"\x1ccom.v2ray.core.common.sendasP\x01Z,github.com/v2fly/v2ray-core/v5/common/sendas\xaa\x02\x18V2Ray.Core.Common.Sendasb\x06proto3"

var (
	file_common_sendas_sendas_proto_rawDescOnce sync.Once
	file_common_sendas_sendas_proto_rawDescData []byte
)

func file_common_sendas_sendas_proto_rawDescGZIP() []byte {
	file_common_sendas_sendas_proto_rawDescOnce.Do(func() {
		file_common_sendas_sendas_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_common_sendas_sendas_proto_rawDesc), len(file_common_sendas_sendas_proto_rawDesc)))
	})
	return file_common_sendas_sendas_proto_rawDescData
}

var file_common_sendas_sendas_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_common_sendas_sendas_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_sendas_sendas_proto_goTypes = []any{
	(Method)(0),    // 0: v2ray.core.common.sendas.Method
	(*Config)(nil), // 1: v2ray.core.common.sendas.Config
}
var file_common_sendas_sendas_proto_depIdxs = []int32{
	0, // 0: v2ray.core.common.sendas.Config.method:type_name -> v2ray.core.common.sendas.Method
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_common_sendas_sendas_proto_init() }
func file_common_sendas_sendas_proto_init() {
	if File_common_sendas_sendas_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_common_sendas_sendas_proto_rawDesc), len(file_common_sendas_sendas_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_sendas_sendas_proto_goTypes,
		DependencyIndexes: file_common_sendas_sendas_proto_depIdxs,
		EnumInfos:         file_common_sendas_sendas_proto_enumTypes,
		MessageInfos:      file_common_sendas_sendas_proto_msgTypes,
	}.Build()
	File_common_sendas_sendas_proto = out.File
	file_common_sendas_sendas_proto_goTypes = nil
	file_common_sendas_sendas_proto_depIdxs = nil
}
