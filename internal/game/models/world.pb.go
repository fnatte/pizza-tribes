// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.21.1
// source: world.proto

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

type WorldEntry_LandType int32

const (
	WorldEntry_GRASS WorldEntry_LandType = 0
)

// Enum value maps for WorldEntry_LandType.
var (
	WorldEntry_LandType_name = map[int32]string{
		0: "GRASS",
	}
	WorldEntry_LandType_value = map[string]int32{
		"GRASS": 0,
	}
)

func (x WorldEntry_LandType) Enum() *WorldEntry_LandType {
	p := new(WorldEntry_LandType)
	*p = x
	return p
}

func (x WorldEntry_LandType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WorldEntry_LandType) Descriptor() protoreflect.EnumDescriptor {
	return file_world_proto_enumTypes[0].Descriptor()
}

func (WorldEntry_LandType) Type() protoreflect.EnumType {
	return &file_world_proto_enumTypes[0]
}

func (x WorldEntry_LandType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WorldEntry_LandType.Descriptor instead.
func (WorldEntry_LandType) EnumDescriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{1, 0}
}

type WorldState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTime int64 `protobuf:"varint,1,opt,name=startTime,proto3" json:"startTime,omitempty"`
	// Types that are assignable to Type:
	//	*WorldState_Starting_
	//	*WorldState_Started_
	//	*WorldState_Ended_
	Type isWorldState_Type `protobuf_oneof:"type"`
}

func (x *WorldState) Reset() {
	*x = WorldState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorldState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorldState) ProtoMessage() {}

func (x *WorldState) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorldState.ProtoReflect.Descriptor instead.
func (*WorldState) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{0}
}

func (x *WorldState) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (m *WorldState) GetType() isWorldState_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *WorldState) GetStarting() *WorldState_Starting {
	if x, ok := x.GetType().(*WorldState_Starting_); ok {
		return x.Starting
	}
	return nil
}

func (x *WorldState) GetStarted() *WorldState_Started {
	if x, ok := x.GetType().(*WorldState_Started_); ok {
		return x.Started
	}
	return nil
}

func (x *WorldState) GetEnded() *WorldState_Ended {
	if x, ok := x.GetType().(*WorldState_Ended_); ok {
		return x.Ended
	}
	return nil
}

type isWorldState_Type interface {
	isWorldState_Type()
}

type WorldState_Starting_ struct {
	Starting *WorldState_Starting `protobuf:"bytes,2,opt,name=starting,proto3,oneof"`
}

type WorldState_Started_ struct {
	Started *WorldState_Started `protobuf:"bytes,3,opt,name=started,proto3,oneof"`
}

type WorldState_Ended_ struct {
	Ended *WorldState_Ended `protobuf:"bytes,4,opt,name=ended,proto3,oneof"`
}

func (*WorldState_Starting_) isWorldState_Type() {}

func (*WorldState_Started_) isWorldState_Type() {}

func (*WorldState_Ended_) isWorldState_Type() {}

type WorldEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LandType WorldEntry_LandType `protobuf:"varint,1,opt,name=landType,proto3,enum=pizzatribes.WorldEntry_LandType" json:"landType,omitempty"`
	// Types that are assignable to Object:
	//	*WorldEntry_Town_
	Object isWorldEntry_Object `protobuf_oneof:"object"`
}

func (x *WorldEntry) Reset() {
	*x = WorldEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorldEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorldEntry) ProtoMessage() {}

func (x *WorldEntry) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorldEntry.ProtoReflect.Descriptor instead.
func (*WorldEntry) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{1}
}

func (x *WorldEntry) GetLandType() WorldEntry_LandType {
	if x != nil {
		return x.LandType
	}
	return WorldEntry_GRASS
}

func (m *WorldEntry) GetObject() isWorldEntry_Object {
	if m != nil {
		return m.Object
	}
	return nil
}

func (x *WorldEntry) GetTown() *WorldEntry_Town {
	if x, ok := x.GetObject().(*WorldEntry_Town_); ok {
		return x.Town
	}
	return nil
}

type isWorldEntry_Object interface {
	isWorldEntry_Object()
}

type WorldEntry_Town_ struct {
	Town *WorldEntry_Town `protobuf:"bytes,2,opt,name=town,proto3,oneof"`
}

func (*WorldEntry_Town_) isWorldEntry_Object() {}

type EntriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries map[string]*WorldEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *EntriesResponse) Reset() {
	*x = EntriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntriesResponse) ProtoMessage() {}

func (x *EntriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntriesResponse.ProtoReflect.Descriptor instead.
func (*EntriesResponse) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{2}
}

func (x *EntriesResponse) GetEntries() map[string]*WorldEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type World struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries map[string]*WorldEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	State   *WorldState            `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *World) Reset() {
	*x = World{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *World) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*World) ProtoMessage() {}

func (x *World) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use World.ProtoReflect.Descriptor instead.
func (*World) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{3}
}

func (x *World) GetEntries() map[string]*WorldEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

func (x *World) GetState() *WorldState {
	if x != nil {
		return x.State
	}
	return nil
}

type WorldState_Starting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WorldState_Starting) Reset() {
	*x = WorldState_Starting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorldState_Starting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorldState_Starting) ProtoMessage() {}

func (x *WorldState_Starting) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorldState_Starting.ProtoReflect.Descriptor instead.
func (*WorldState_Starting) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{0, 0}
}

type WorldState_Started struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WorldState_Started) Reset() {
	*x = WorldState_Started{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorldState_Started) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorldState_Started) ProtoMessage() {}

func (x *WorldState_Started) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorldState_Started.ProtoReflect.Descriptor instead.
func (*WorldState_Started) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{0, 1}
}

type WorldState_Ended struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WinnerUserId string `protobuf:"bytes,1,opt,name=winnerUserId,proto3" json:"winnerUserId,omitempty"`
	EndedAt      int64  `protobuf:"varint,2,opt,name=endedAt,proto3" json:"endedAt,omitempty"`
}

func (x *WorldState_Ended) Reset() {
	*x = WorldState_Ended{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorldState_Ended) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorldState_Ended) ProtoMessage() {}

func (x *WorldState_Ended) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorldState_Ended.ProtoReflect.Descriptor instead.
func (*WorldState_Ended) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{0, 2}
}

func (x *WorldState_Ended) GetWinnerUserId() string {
	if x != nil {
		return x.WinnerUserId
	}
	return ""
}

func (x *WorldState_Ended) GetEndedAt() int64 {
	if x != nil {
		return x.EndedAt
	}
	return 0
}

type WorldEntry_Town struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *WorldEntry_Town) Reset() {
	*x = WorldEntry_Town{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorldEntry_Town) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorldEntry_Town) ProtoMessage() {}

func (x *WorldEntry_Town) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorldEntry_Town.ProtoReflect.Descriptor instead.
func (*WorldEntry_Town) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{1, 0}
}

func (x *WorldEntry_Town) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

var File_world_proto protoreflect.FileDescriptor

var file_world_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x70,
	0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x22, 0xc4, 0x02, 0x0a, 0x0a, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3e, 0x0a, 0x08, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x69, 0x7a, 0x7a,
	0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x48, 0x00, 0x52, 0x08, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x3b, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61,
	0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x48, 0x00, 0x52, 0x07, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x65, 0x64, 0x12, 0x35, 0x0a, 0x05, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65,
	0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x45, 0x6e, 0x64,
	0x65, 0x64, 0x48, 0x00, 0x52, 0x05, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x1a, 0x0a, 0x0a, 0x08, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x69, 0x6e, 0x67, 0x1a, 0x09, 0x0a, 0x07, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x65, 0x64, 0x1a, 0x45, 0x0a, 0x05, 0x45, 0x6e, 0x64, 0x65, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x77,
	0x69, 0x6e, 0x6e, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x77, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x41, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x07, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x41, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x22, 0xbf, 0x01, 0x0a, 0x0a, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x3c, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x20, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73,
	0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x4c, 0x61, 0x6e, 0x64,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x32,
	0x0a, 0x04, 0x74, 0x6f, 0x77, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70,
	0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x54, 0x6f, 0x77, 0x6e, 0x48, 0x00, 0x52, 0x04, 0x74, 0x6f,
	0x77, 0x6e, 0x1a, 0x1e, 0x0a, 0x04, 0x54, 0x6f, 0x77, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x22, 0x15, 0x0a, 0x08, 0x4c, 0x61, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09,
	0x0a, 0x05, 0x47, 0x52, 0x41, 0x53, 0x53, 0x10, 0x00, 0x42, 0x08, 0x0a, 0x06, 0x6f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x22, 0xab, 0x01, 0x0a, 0x0f, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61,
	0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x1a, 0x53, 0x0a, 0x0c,
	0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2d,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6c,
	0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0xc6, 0x01, 0x0a, 0x05, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x12, 0x39, 0x0a, 0x07, 0x65,
	0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x70,
	0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64,
	0x2e, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65,
	0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x2d, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69,
	0x62, 0x65, 0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x1a, 0x53, 0x0a, 0x0c, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2d, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72,
	0x69, 0x62, 0x65, 0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6e, 0x61, 0x74, 0x74, 0x65, 0x2f,
	0x70, 0x69, 0x7a, 0x7a, 0x61, 0x2d, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_world_proto_rawDescOnce sync.Once
	file_world_proto_rawDescData = file_world_proto_rawDesc
)

func file_world_proto_rawDescGZIP() []byte {
	file_world_proto_rawDescOnce.Do(func() {
		file_world_proto_rawDescData = protoimpl.X.CompressGZIP(file_world_proto_rawDescData)
	})
	return file_world_proto_rawDescData
}

var file_world_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_world_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_world_proto_goTypes = []interface{}{
	(WorldEntry_LandType)(0),    // 0: pizzatribes.WorldEntry.LandType
	(*WorldState)(nil),          // 1: pizzatribes.WorldState
	(*WorldEntry)(nil),          // 2: pizzatribes.WorldEntry
	(*EntriesResponse)(nil),     // 3: pizzatribes.EntriesResponse
	(*World)(nil),               // 4: pizzatribes.World
	(*WorldState_Starting)(nil), // 5: pizzatribes.WorldState.Starting
	(*WorldState_Started)(nil),  // 6: pizzatribes.WorldState.Started
	(*WorldState_Ended)(nil),    // 7: pizzatribes.WorldState.Ended
	(*WorldEntry_Town)(nil),     // 8: pizzatribes.WorldEntry.Town
	nil,                         // 9: pizzatribes.EntriesResponse.EntriesEntry
	nil,                         // 10: pizzatribes.World.EntriesEntry
}
var file_world_proto_depIdxs = []int32{
	5,  // 0: pizzatribes.WorldState.starting:type_name -> pizzatribes.WorldState.Starting
	6,  // 1: pizzatribes.WorldState.started:type_name -> pizzatribes.WorldState.Started
	7,  // 2: pizzatribes.WorldState.ended:type_name -> pizzatribes.WorldState.Ended
	0,  // 3: pizzatribes.WorldEntry.landType:type_name -> pizzatribes.WorldEntry.LandType
	8,  // 4: pizzatribes.WorldEntry.town:type_name -> pizzatribes.WorldEntry.Town
	9,  // 5: pizzatribes.EntriesResponse.entries:type_name -> pizzatribes.EntriesResponse.EntriesEntry
	10, // 6: pizzatribes.World.entries:type_name -> pizzatribes.World.EntriesEntry
	1,  // 7: pizzatribes.World.state:type_name -> pizzatribes.WorldState
	2,  // 8: pizzatribes.EntriesResponse.EntriesEntry.value:type_name -> pizzatribes.WorldEntry
	2,  // 9: pizzatribes.World.EntriesEntry.value:type_name -> pizzatribes.WorldEntry
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_world_proto_init() }
func file_world_proto_init() {
	if File_world_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_world_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorldState); i {
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
		file_world_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorldEntry); i {
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
		file_world_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntriesResponse); i {
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
		file_world_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*World); i {
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
		file_world_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorldState_Starting); i {
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
		file_world_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorldState_Started); i {
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
		file_world_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorldState_Ended); i {
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
		file_world_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorldEntry_Town); i {
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
	file_world_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*WorldState_Starting_)(nil),
		(*WorldState_Started_)(nil),
		(*WorldState_Ended_)(nil),
	}
	file_world_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*WorldEntry_Town_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_world_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_world_proto_goTypes,
		DependencyIndexes: file_world_proto_depIdxs,
		EnumInfos:         file_world_proto_enumTypes,
		MessageInfos:      file_world_proto_msgTypes,
	}.Build()
	File_world_proto = out.File
	file_world_proto_rawDesc = nil
	file_world_proto_goTypes = nil
	file_world_proto_depIdxs = nil
}