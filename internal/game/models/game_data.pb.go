// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.21.4
// source: game_data.proto

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

type EducationInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title       string    `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	TitlePlural string    `protobuf:"bytes,2,opt,name=titlePlural,proto3" json:"titlePlural,omitempty"`
	Cost        int32     `protobuf:"varint,3,opt,name=cost,proto3" json:"cost,omitempty"`
	TrainTime   int32     `protobuf:"varint,4,opt,name=trainTime,proto3" json:"trainTime,omitempty"`
	Employer    *Building `protobuf:"varint,5,opt,name=employer,proto3,enum=pizzatribes.Building,oneof" json:"employer,omitempty"`
}

func (x *EducationInfo) Reset() {
	*x = EducationInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EducationInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EducationInfo) ProtoMessage() {}

func (x *EducationInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EducationInfo.ProtoReflect.Descriptor instead.
func (*EducationInfo) Descriptor() ([]byte, []int) {
	return file_game_data_proto_rawDescGZIP(), []int{0}
}

func (x *EducationInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *EducationInfo) GetTitlePlural() string {
	if x != nil {
		return x.TitlePlural
	}
	return ""
}

func (x *EducationInfo) GetCost() int32 {
	if x != nil {
		return x.Cost
	}
	return 0
}

func (x *EducationInfo) GetTrainTime() int32 {
	if x != nil {
		return x.TrainTime
	}
	return 0
}

func (x *EducationInfo) GetEmployer() Building {
	if x != nil && x.Employer != nil {
		return *x.Employer
	}
	return Building_KITCHEN
}

type GeniusFlashCost struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Coins  int32 `protobuf:"varint,1,opt,name=coins,proto3" json:"coins,omitempty"`
	Pizzas int32 `protobuf:"varint,2,opt,name=pizzas,proto3" json:"pizzas,omitempty"`
}

func (x *GeniusFlashCost) Reset() {
	*x = GeniusFlashCost{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeniusFlashCost) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeniusFlashCost) ProtoMessage() {}

func (x *GeniusFlashCost) ProtoReflect() protoreflect.Message {
	mi := &file_game_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeniusFlashCost.ProtoReflect.Descriptor instead.
func (*GeniusFlashCost) Descriptor() ([]byte, []int) {
	return file_game_data_proto_rawDescGZIP(), []int{1}
}

func (x *GeniusFlashCost) GetCoins() int32 {
	if x != nil {
		return x.Coins
	}
	return 0
}

func (x *GeniusFlashCost) GetPizzas() int32 {
	if x != nil {
		return x.Pizzas
	}
	return 0
}

type GameData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Buildings        map[int32]*BuildingInfo    `protobuf:"bytes,1,rep,name=buildings,proto3" json:"buildings,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Educations       map[int32]*EducationInfo   `protobuf:"bytes,2,rep,name=educations,proto3" json:"educations,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Research         map[int32]*ResearchInfo    `protobuf:"bytes,3,rep,name=research,proto3" json:"research,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Quests           []*Quest                   `protobuf:"bytes,4,rep,name=quests,proto3" json:"quests,omitempty"`
	AppearanceParts  map[string]*AppearancePart `protobuf:"bytes,5,rep,name=appearanceParts,proto3" json:"appearanceParts,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	GeniusFlashCosts []*GeniusFlashCost         `protobuf:"bytes,6,rep,name=geniusFlashCosts,proto3" json:"geniusFlashCosts,omitempty"`
}

func (x *GameData) Reset() {
	*x = GameData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_data_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameData) ProtoMessage() {}

func (x *GameData) ProtoReflect() protoreflect.Message {
	mi := &file_game_data_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameData.ProtoReflect.Descriptor instead.
func (*GameData) Descriptor() ([]byte, []int) {
	return file_game_data_proto_rawDescGZIP(), []int{2}
}

func (x *GameData) GetBuildings() map[int32]*BuildingInfo {
	if x != nil {
		return x.Buildings
	}
	return nil
}

func (x *GameData) GetEducations() map[int32]*EducationInfo {
	if x != nil {
		return x.Educations
	}
	return nil
}

func (x *GameData) GetResearch() map[int32]*ResearchInfo {
	if x != nil {
		return x.Research
	}
	return nil
}

func (x *GameData) GetQuests() []*Quest {
	if x != nil {
		return x.Quests
	}
	return nil
}

func (x *GameData) GetAppearanceParts() map[string]*AppearancePart {
	if x != nil {
		return x.AppearanceParts
	}
	return nil
}

func (x *GameData) GetGeniusFlashCosts() []*GeniusFlashCost {
	if x != nil {
		return x.GeniusFlashCosts
	}
	return nil
}

var File_game_data_proto protoreflect.FileDescriptor

var file_game_data_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0b, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x1a, 0x10,
	0x61, 0x70, 0x70, 0x65, 0x61, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x0e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x0e, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x0b, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbe, 0x01,
	0x0a, 0x0d, 0x45, 0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x50, 0x6c,
	0x75, 0x72, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x50, 0x6c, 0x75, 0x72, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x72, 0x61, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x74, 0x72, 0x61, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x36, 0x0a, 0x08, 0x65, 0x6d, 0x70,
	0x6c, 0x6f, 0x79, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x70, 0x69,
	0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x69,
	0x6e, 0x67, 0x48, 0x00, 0x52, 0x08, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x88, 0x01,
	0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x65, 0x6d, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x72, 0x22, 0x3f,
	0x0a, 0x0f, 0x47, 0x65, 0x6e, 0x69, 0x75, 0x73, 0x46, 0x6c, 0x61, 0x73, 0x68, 0x43, 0x6f, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x69, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x63, 0x6f, 0x69, 0x6e, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x69, 0x7a, 0x7a, 0x61,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x73, 0x22,
	0x8f, 0x06, 0x0a, 0x08, 0x47, 0x61, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x42, 0x0a, 0x09,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x24, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x47, 0x61,
	0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x73,
	0x12, 0x45, 0x0a, 0x0a, 0x65, 0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62,
	0x65, 0x73, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x45, 0x64, 0x75, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x65, 0x64, 0x75,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3f, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x70, 0x69, 0x7a, 0x7a,
	0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x2e, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08,
	0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x2a, 0x0a, 0x06, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61,
	0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x51, 0x75, 0x65, 0x73, 0x74, 0x52, 0x06, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x73, 0x12, 0x54, 0x0a, 0x0f, 0x61, 0x70, 0x70, 0x65, 0x61, 0x72, 0x61, 0x6e,
	0x63, 0x65, 0x50, 0x61, 0x72, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e,
	0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x47, 0x61, 0x6d, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x61, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x50,
	0x61, 0x72, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0f, 0x61, 0x70, 0x70, 0x65, 0x61,
	0x72, 0x61, 0x6e, 0x63, 0x65, 0x50, 0x61, 0x72, 0x74, 0x73, 0x12, 0x48, 0x0a, 0x10, 0x67, 0x65,
	0x6e, 0x69, 0x75, 0x73, 0x46, 0x6c, 0x61, 0x73, 0x68, 0x43, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x06,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62,
	0x65, 0x73, 0x2e, 0x47, 0x65, 0x6e, 0x69, 0x75, 0x73, 0x46, 0x6c, 0x61, 0x73, 0x68, 0x43, 0x6f,
	0x73, 0x74, 0x52, 0x10, 0x67, 0x65, 0x6e, 0x69, 0x75, 0x73, 0x46, 0x6c, 0x61, 0x73, 0x68, 0x43,
	0x6f, 0x73, 0x74, 0x73, 0x1a, 0x57, 0x0a, 0x0e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x69, 0x6e, 0x67,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2f, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74,
	0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x59, 0x0a,
	0x0f, 0x45, 0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x30, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e,
	0x45, 0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x56, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2f, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x69, 0x7a,
	0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x1a, 0x5f, 0x0a, 0x14, 0x41, 0x70, 0x70, 0x65, 0x61, 0x72, 0x61, 0x6e, 0x63, 0x65, 0x50, 0x61,
	0x72, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x31, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x69, 0x7a, 0x7a,
	0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x61, 0x72, 0x61, 0x6e,
	0x63, 0x65, 0x50, 0x61, 0x72, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x66, 0x6e, 0x61, 0x74, 0x74, 0x65, 0x2f, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x2d, 0x74, 0x72, 0x69,
	0x62, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x61, 0x6d,
	0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_game_data_proto_rawDescOnce sync.Once
	file_game_data_proto_rawDescData = file_game_data_proto_rawDesc
)

func file_game_data_proto_rawDescGZIP() []byte {
	file_game_data_proto_rawDescOnce.Do(func() {
		file_game_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_game_data_proto_rawDescData)
	})
	return file_game_data_proto_rawDescData
}

var file_game_data_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_game_data_proto_goTypes = []interface{}{
	(*EducationInfo)(nil),   // 0: pizzatribes.EducationInfo
	(*GeniusFlashCost)(nil), // 1: pizzatribes.GeniusFlashCost
	(*GameData)(nil),        // 2: pizzatribes.GameData
	nil,                     // 3: pizzatribes.GameData.BuildingsEntry
	nil,                     // 4: pizzatribes.GameData.EducationsEntry
	nil,                     // 5: pizzatribes.GameData.ResearchEntry
	nil,                     // 6: pizzatribes.GameData.AppearancePartsEntry
	(Building)(0),           // 7: pizzatribes.Building
	(*Quest)(nil),           // 8: pizzatribes.Quest
	(*BuildingInfo)(nil),    // 9: pizzatribes.BuildingInfo
	(*ResearchInfo)(nil),    // 10: pizzatribes.ResearchInfo
	(*AppearancePart)(nil),  // 11: pizzatribes.AppearancePart
}
var file_game_data_proto_depIdxs = []int32{
	7,  // 0: pizzatribes.EducationInfo.employer:type_name -> pizzatribes.Building
	3,  // 1: pizzatribes.GameData.buildings:type_name -> pizzatribes.GameData.BuildingsEntry
	4,  // 2: pizzatribes.GameData.educations:type_name -> pizzatribes.GameData.EducationsEntry
	5,  // 3: pizzatribes.GameData.research:type_name -> pizzatribes.GameData.ResearchEntry
	8,  // 4: pizzatribes.GameData.quests:type_name -> pizzatribes.Quest
	6,  // 5: pizzatribes.GameData.appearanceParts:type_name -> pizzatribes.GameData.AppearancePartsEntry
	1,  // 6: pizzatribes.GameData.geniusFlashCosts:type_name -> pizzatribes.GeniusFlashCost
	9,  // 7: pizzatribes.GameData.BuildingsEntry.value:type_name -> pizzatribes.BuildingInfo
	0,  // 8: pizzatribes.GameData.EducationsEntry.value:type_name -> pizzatribes.EducationInfo
	10, // 9: pizzatribes.GameData.ResearchEntry.value:type_name -> pizzatribes.ResearchInfo
	11, // 10: pizzatribes.GameData.AppearancePartsEntry.value:type_name -> pizzatribes.AppearancePart
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_game_data_proto_init() }
func file_game_data_proto_init() {
	if File_game_data_proto != nil {
		return
	}
	file_appearance_proto_init()
	file_building_proto_init()
	file_research_proto_init()
	file_quest_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_game_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EducationInfo); i {
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
		file_game_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GeniusFlashCost); i {
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
		file_game_data_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameData); i {
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
	file_game_data_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_game_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_game_data_proto_goTypes,
		DependencyIndexes: file_game_data_proto_depIdxs,
		MessageInfos:      file_game_data_proto_msgTypes,
	}.Build()
	File_game_data_proto = out.File
	file_game_data_proto_rawDesc = nil
	file_game_data_proto_goTypes = nil
	file_game_data_proto_depIdxs = nil
}
