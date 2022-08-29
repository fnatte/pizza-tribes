// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.21.1
// source: timeseries.proto

package models

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type DataPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int64   `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Coins     float64 `protobuf:"fixed64,2,opt,name=coins,proto3" json:"coins,omitempty"`
	Pizzas    float64 `protobuf:"fixed64,3,opt,name=pizzas,proto3" json:"pizzas,omitempty"`
}

func (x *DataPoint) Reset() {
	*x = DataPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_timeseries_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataPoint) ProtoMessage() {}

func (x *DataPoint) ProtoReflect() protoreflect.Message {
	mi := &file_timeseries_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataPoint.ProtoReflect.Descriptor instead.
func (*DataPoint) Descriptor() ([]byte, []int) {
	return file_timeseries_proto_rawDescGZIP(), []int{0}
}

func (x *DataPoint) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *DataPoint) GetCoins() float64 {
	if x != nil {
		return x.Coins
	}
	return 0
}

func (x *DataPoint) GetPizzas() float64 {
	if x != nil {
		return x.Pizzas
	}
	return 0
}

type TimeseriesData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataPoints []*DataPoint `protobuf:"bytes,1,rep,name=dataPoints,proto3" json:"dataPoints,omitempty"`
}

func (x *TimeseriesData) Reset() {
	*x = TimeseriesData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_timeseries_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeseriesData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeseriesData) ProtoMessage() {}

func (x *TimeseriesData) ProtoReflect() protoreflect.Message {
	mi := &file_timeseries_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeseriesData.ProtoReflect.Descriptor instead.
func (*TimeseriesData) Descriptor() ([]byte, []int) {
	return file_timeseries_proto_rawDescGZIP(), []int{1}
}

func (x *TimeseriesData) GetDataPoints() []*DataPoint {
	if x != nil {
		return x.DataPoints
	}
	return nil
}

var File_timeseries_proto protoreflect.FileDescriptor

var file_timeseries_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0b, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x22,
	0x57, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x69, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6f, 0x69, 0x6e, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x06, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x73, 0x22, 0x48, 0x0a, 0x0e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x36, 0x0a, 0x0a, 0x64, 0x61,
	0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e,
	0x74, 0x73, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x66, 0x6e, 0x61, 0x74, 0x74, 0x65, 0x2f, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x2d, 0x74, 0x72,
	0x69, 0x62, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x61,
	0x6d, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_timeseries_proto_rawDescOnce sync.Once
	file_timeseries_proto_rawDescData = file_timeseries_proto_rawDesc
)

func file_timeseries_proto_rawDescGZIP() []byte {
	file_timeseries_proto_rawDescOnce.Do(func() {
		file_timeseries_proto_rawDescData = protoimpl.X.CompressGZIP(file_timeseries_proto_rawDescData)
	})
	return file_timeseries_proto_rawDescData
}

var file_timeseries_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_timeseries_proto_goTypes = []interface{}{
	(*DataPoint)(nil),      // 0: pizzatribes.DataPoint
	(*TimeseriesData)(nil), // 1: pizzatribes.TimeseriesData
}
var file_timeseries_proto_depIdxs = []int32{
	0, // 0: pizzatribes.TimeseriesData.dataPoints:type_name -> pizzatribes.DataPoint
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_timeseries_proto_init() }
func file_timeseries_proto_init() {
	if File_timeseries_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_timeseries_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataPoint); i {
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
		file_timeseries_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeseriesData); i {
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
			RawDescriptor: file_timeseries_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_timeseries_proto_goTypes,
		DependencyIndexes: file_timeseries_proto_depIdxs,
		MessageInfos:      file_timeseries_proto_msgTypes,
	}.Build()
	File_timeseries_proto = out.File
	file_timeseries_proto_rawDesc = nil
	file_timeseries_proto_goTypes = nil
	file_timeseries_proto_depIdxs = nil
}