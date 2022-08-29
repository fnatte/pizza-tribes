// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.21.1
// source: research.proto

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

type ResearchDiscovery int32

const (
	ResearchDiscovery_WEBSITE                 ResearchDiscovery = 0
	ResearchDiscovery_DIGITAL_ORDERING_SYSTEM ResearchDiscovery = 1
	ResearchDiscovery_MOBILE_APP              ResearchDiscovery = 2
	ResearchDiscovery_MASONRY_OVEN            ResearchDiscovery = 3
	ResearchDiscovery_GAS_OVEN                ResearchDiscovery = 4
	ResearchDiscovery_HYBRID_OVEN             ResearchDiscovery = 5
	ResearchDiscovery_DURUM_WHEAT             ResearchDiscovery = 6
	ResearchDiscovery_DOUBLE_ZERO_FLOUR       ResearchDiscovery = 7
	ResearchDiscovery_SAN_MARZANO_TOMATOES    ResearchDiscovery = 8
	ResearchDiscovery_OCIMUM_BASILICUM        ResearchDiscovery = 9
	ResearchDiscovery_EXTRA_VIRGIN            ResearchDiscovery = 10
	ResearchDiscovery_WHITEBOARD              ResearchDiscovery = 11
	ResearchDiscovery_KITCHEN_STRATEGY        ResearchDiscovery = 12
	ResearchDiscovery_STRESS_HANDLING         ResearchDiscovery = 13
	ResearchDiscovery_SLAM                    ResearchDiscovery = 14
	ResearchDiscovery_HIT_IT                  ResearchDiscovery = 15
	ResearchDiscovery_GRAND_SLAM              ResearchDiscovery = 16
	ResearchDiscovery_GODS_TOUCH              ResearchDiscovery = 17
	ResearchDiscovery_CONSECUTIVE             ResearchDiscovery = 18
	ResearchDiscovery_ON_A_ROLL               ResearchDiscovery = 19
	ResearchDiscovery_BOOTS_OF_HASTE          ResearchDiscovery = 20
	ResearchDiscovery_TIP_TOE                 ResearchDiscovery = 21
	ResearchDiscovery_SHADOW_EXPERT           ResearchDiscovery = 22
	ResearchDiscovery_BIG_POCKETS             ResearchDiscovery = 23
	ResearchDiscovery_THIEVES_FAVORITE_BAG    ResearchDiscovery = 24
	ResearchDiscovery_COFFEE                  ResearchDiscovery = 25
	ResearchDiscovery_NIGHTS_WATCH            ResearchDiscovery = 26
	ResearchDiscovery_LASER_ALARM             ResearchDiscovery = 27
	ResearchDiscovery_TRIP_WIRE               ResearchDiscovery = 28
	ResearchDiscovery_CARDIO                  ResearchDiscovery = 29
)

// Enum value maps for ResearchDiscovery.
var (
	ResearchDiscovery_name = map[int32]string{
		0:  "WEBSITE",
		1:  "DIGITAL_ORDERING_SYSTEM",
		2:  "MOBILE_APP",
		3:  "MASONRY_OVEN",
		4:  "GAS_OVEN",
		5:  "HYBRID_OVEN",
		6:  "DURUM_WHEAT",
		7:  "DOUBLE_ZERO_FLOUR",
		8:  "SAN_MARZANO_TOMATOES",
		9:  "OCIMUM_BASILICUM",
		10: "EXTRA_VIRGIN",
		11: "WHITEBOARD",
		12: "KITCHEN_STRATEGY",
		13: "STRESS_HANDLING",
		14: "SLAM",
		15: "HIT_IT",
		16: "GRAND_SLAM",
		17: "GODS_TOUCH",
		18: "CONSECUTIVE",
		19: "ON_A_ROLL",
		20: "BOOTS_OF_HASTE",
		21: "TIP_TOE",
		22: "SHADOW_EXPERT",
		23: "BIG_POCKETS",
		24: "THIEVES_FAVORITE_BAG",
		25: "COFFEE",
		26: "NIGHTS_WATCH",
		27: "LASER_ALARM",
		28: "TRIP_WIRE",
		29: "CARDIO",
	}
	ResearchDiscovery_value = map[string]int32{
		"WEBSITE":                 0,
		"DIGITAL_ORDERING_SYSTEM": 1,
		"MOBILE_APP":              2,
		"MASONRY_OVEN":            3,
		"GAS_OVEN":                4,
		"HYBRID_OVEN":             5,
		"DURUM_WHEAT":             6,
		"DOUBLE_ZERO_FLOUR":       7,
		"SAN_MARZANO_TOMATOES":    8,
		"OCIMUM_BASILICUM":        9,
		"EXTRA_VIRGIN":            10,
		"WHITEBOARD":              11,
		"KITCHEN_STRATEGY":        12,
		"STRESS_HANDLING":         13,
		"SLAM":                    14,
		"HIT_IT":                  15,
		"GRAND_SLAM":              16,
		"GODS_TOUCH":              17,
		"CONSECUTIVE":             18,
		"ON_A_ROLL":               19,
		"BOOTS_OF_HASTE":          20,
		"TIP_TOE":                 21,
		"SHADOW_EXPERT":           22,
		"BIG_POCKETS":             23,
		"THIEVES_FAVORITE_BAG":    24,
		"COFFEE":                  25,
		"NIGHTS_WATCH":            26,
		"LASER_ALARM":             27,
		"TRIP_WIRE":               28,
		"CARDIO":                  29,
	}
)

func (x ResearchDiscovery) Enum() *ResearchDiscovery {
	p := new(ResearchDiscovery)
	*p = x
	return p
}

func (x ResearchDiscovery) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResearchDiscovery) Descriptor() protoreflect.EnumDescriptor {
	return file_research_proto_enumTypes[0].Descriptor()
}

func (ResearchDiscovery) Type() protoreflect.EnumType {
	return &file_research_proto_enumTypes[0]
}

func (x ResearchDiscovery) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResearchDiscovery.Descriptor instead.
func (ResearchDiscovery) EnumDescriptor() ([]byte, []int) {
	return file_research_proto_rawDescGZIP(), []int{0}
}

type ResearchTree int32

const (
	ResearchTree_PRODUCTION ResearchTree = 0
	ResearchTree_DEMAND     ResearchTree = 1
	ResearchTree_TAPPING    ResearchTree = 2
	ResearchTree_THIEVES    ResearchTree = 3
	ResearchTree_GUARDS     ResearchTree = 4
)

// Enum value maps for ResearchTree.
var (
	ResearchTree_name = map[int32]string{
		0: "PRODUCTION",
		1: "DEMAND",
		2: "TAPPING",
		3: "THIEVES",
		4: "GUARDS",
	}
	ResearchTree_value = map[string]int32{
		"PRODUCTION": 0,
		"DEMAND":     1,
		"TAPPING":    2,
		"THIEVES":    3,
		"GUARDS":     4,
	}
)

func (x ResearchTree) Enum() *ResearchTree {
	p := new(ResearchTree)
	*p = x
	return p
}

func (x ResearchTree) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResearchTree) Descriptor() protoreflect.EnumDescriptor {
	return file_research_proto_enumTypes[1].Descriptor()
}

func (ResearchTree) Type() protoreflect.EnumType {
	return &file_research_proto_enumTypes[1]
}

func (x ResearchTree) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResearchTree.Descriptor instead.
func (ResearchTree) EnumDescriptor() ([]byte, []int) {
	return file_research_proto_rawDescGZIP(), []int{1}
}

type ResearchInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title        string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Discovery    ResearchDiscovery      `protobuf:"varint,2,opt,name=discovery,proto3,enum=pizzatribes.ResearchDiscovery" json:"discovery,omitempty"`
	Requirements []ResearchDiscovery    `protobuf:"varint,3,rep,packed,name=requirements,proto3,enum=pizzatribes.ResearchDiscovery" json:"requirements,omitempty"`
	ResearchTime int32                  `protobuf:"varint,4,opt,name=researchTime,proto3" json:"researchTime,omitempty"`
	Tree         ResearchTree           `protobuf:"varint,5,opt,name=tree,proto3,enum=pizzatribes.ResearchTree" json:"tree,omitempty"`
	Description  string                 `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	Rewards      []*ResearchInfo_Reward `protobuf:"bytes,7,rep,name=rewards,proto3" json:"rewards,omitempty"`
	X            int32                  `protobuf:"varint,8,opt,name=x,proto3" json:"x,omitempty"`
	Y            int32                  `protobuf:"varint,9,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *ResearchInfo) Reset() {
	*x = ResearchInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_research_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResearchInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResearchInfo) ProtoMessage() {}

func (x *ResearchInfo) ProtoReflect() protoreflect.Message {
	mi := &file_research_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResearchInfo.ProtoReflect.Descriptor instead.
func (*ResearchInfo) Descriptor() ([]byte, []int) {
	return file_research_proto_rawDescGZIP(), []int{0}
}

func (x *ResearchInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ResearchInfo) GetDiscovery() ResearchDiscovery {
	if x != nil {
		return x.Discovery
	}
	return ResearchDiscovery_WEBSITE
}

func (x *ResearchInfo) GetRequirements() []ResearchDiscovery {
	if x != nil {
		return x.Requirements
	}
	return nil
}

func (x *ResearchInfo) GetResearchTime() int32 {
	if x != nil {
		return x.ResearchTime
	}
	return 0
}

func (x *ResearchInfo) GetTree() ResearchTree {
	if x != nil {
		return x.Tree
	}
	return ResearchTree_PRODUCTION
}

func (x *ResearchInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ResearchInfo) GetRewards() []*ResearchInfo_Reward {
	if x != nil {
		return x.Rewards
	}
	return nil
}

func (x *ResearchInfo) GetX() int32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *ResearchInfo) GetY() int32 {
	if x != nil {
		return x.Y
	}
	return 0
}

type ResearchInfo_Reward struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attribute string `protobuf:"bytes,1,opt,name=attribute,proto3" json:"attribute,omitempty"`
	Value     string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ResearchInfo_Reward) Reset() {
	*x = ResearchInfo_Reward{}
	if protoimpl.UnsafeEnabled {
		mi := &file_research_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResearchInfo_Reward) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResearchInfo_Reward) ProtoMessage() {}

func (x *ResearchInfo_Reward) ProtoReflect() protoreflect.Message {
	mi := &file_research_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResearchInfo_Reward.ProtoReflect.Descriptor instead.
func (*ResearchInfo_Reward) Descriptor() ([]byte, []int) {
	return file_research_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ResearchInfo_Reward) GetAttribute() string {
	if x != nil {
		return x.Attribute
	}
	return ""
}

func (x *ResearchInfo_Reward) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_research_proto protoreflect.FileDescriptor

var file_research_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x22, 0xb1, 0x03,
	0x0a, 0x0c, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x12, 0x3c, 0x0a, 0x09, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61, 0x74,
	0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x09, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x12, 0x42, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61,
	0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x72, 0x65,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x04, 0x74, 0x72,
	0x65, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x70, 0x69, 0x7a, 0x7a, 0x61,
	0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x54,
	0x72, 0x65, 0x65, 0x52, 0x04, 0x74, 0x72, 0x65, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x07, 0x72,
	0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70,
	0x69, 0x7a, 0x7a, 0x61, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x52, 0x07,
	0x72, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x01, 0x79, 0x1a, 0x3c, 0x0a, 0x06, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x2a, 0xa2, 0x04, 0x0a, 0x11, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x45, 0x42, 0x53, 0x49,
	0x54, 0x45, 0x10, 0x00, 0x12, 0x1b, 0x0a, 0x17, 0x44, 0x49, 0x47, 0x49, 0x54, 0x41, 0x4c, 0x5f,
	0x4f, 0x52, 0x44, 0x45, 0x52, 0x49, 0x4e, 0x47, 0x5f, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x10,
	0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x4d, 0x4f, 0x42, 0x49, 0x4c, 0x45, 0x5f, 0x41, 0x50, 0x50, 0x10,
	0x02, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x41, 0x53, 0x4f, 0x4e, 0x52, 0x59, 0x5f, 0x4f, 0x56, 0x45,
	0x4e, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08, 0x47, 0x41, 0x53, 0x5f, 0x4f, 0x56, 0x45, 0x4e, 0x10,
	0x04, 0x12, 0x0f, 0x0a, 0x0b, 0x48, 0x59, 0x42, 0x52, 0x49, 0x44, 0x5f, 0x4f, 0x56, 0x45, 0x4e,
	0x10, 0x05, 0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x55, 0x52, 0x55, 0x4d, 0x5f, 0x57, 0x48, 0x45, 0x41,
	0x54, 0x10, 0x06, 0x12, 0x15, 0x0a, 0x11, 0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x5f, 0x5a, 0x45,
	0x52, 0x4f, 0x5f, 0x46, 0x4c, 0x4f, 0x55, 0x52, 0x10, 0x07, 0x12, 0x18, 0x0a, 0x14, 0x53, 0x41,
	0x4e, 0x5f, 0x4d, 0x41, 0x52, 0x5a, 0x41, 0x4e, 0x4f, 0x5f, 0x54, 0x4f, 0x4d, 0x41, 0x54, 0x4f,
	0x45, 0x53, 0x10, 0x08, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x43, 0x49, 0x4d, 0x55, 0x4d, 0x5f, 0x42,
	0x41, 0x53, 0x49, 0x4c, 0x49, 0x43, 0x55, 0x4d, 0x10, 0x09, 0x12, 0x10, 0x0a, 0x0c, 0x45, 0x58,
	0x54, 0x52, 0x41, 0x5f, 0x56, 0x49, 0x52, 0x47, 0x49, 0x4e, 0x10, 0x0a, 0x12, 0x0e, 0x0a, 0x0a,
	0x57, 0x48, 0x49, 0x54, 0x45, 0x42, 0x4f, 0x41, 0x52, 0x44, 0x10, 0x0b, 0x12, 0x14, 0x0a, 0x10,
	0x4b, 0x49, 0x54, 0x43, 0x48, 0x45, 0x4e, 0x5f, 0x53, 0x54, 0x52, 0x41, 0x54, 0x45, 0x47, 0x59,
	0x10, 0x0c, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x54, 0x52, 0x45, 0x53, 0x53, 0x5f, 0x48, 0x41, 0x4e,
	0x44, 0x4c, 0x49, 0x4e, 0x47, 0x10, 0x0d, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x4c, 0x41, 0x4d, 0x10,
	0x0e, 0x12, 0x0a, 0x0a, 0x06, 0x48, 0x49, 0x54, 0x5f, 0x49, 0x54, 0x10, 0x0f, 0x12, 0x0e, 0x0a,
	0x0a, 0x47, 0x52, 0x41, 0x4e, 0x44, 0x5f, 0x53, 0x4c, 0x41, 0x4d, 0x10, 0x10, 0x12, 0x0e, 0x0a,
	0x0a, 0x47, 0x4f, 0x44, 0x53, 0x5f, 0x54, 0x4f, 0x55, 0x43, 0x48, 0x10, 0x11, 0x12, 0x0f, 0x0a,
	0x0b, 0x43, 0x4f, 0x4e, 0x53, 0x45, 0x43, 0x55, 0x54, 0x49, 0x56, 0x45, 0x10, 0x12, 0x12, 0x0d,
	0x0a, 0x09, 0x4f, 0x4e, 0x5f, 0x41, 0x5f, 0x52, 0x4f, 0x4c, 0x4c, 0x10, 0x13, 0x12, 0x12, 0x0a,
	0x0e, 0x42, 0x4f, 0x4f, 0x54, 0x53, 0x5f, 0x4f, 0x46, 0x5f, 0x48, 0x41, 0x53, 0x54, 0x45, 0x10,
	0x14, 0x12, 0x0b, 0x0a, 0x07, 0x54, 0x49, 0x50, 0x5f, 0x54, 0x4f, 0x45, 0x10, 0x15, 0x12, 0x11,
	0x0a, 0x0d, 0x53, 0x48, 0x41, 0x44, 0x4f, 0x57, 0x5f, 0x45, 0x58, 0x50, 0x45, 0x52, 0x54, 0x10,
	0x16, 0x12, 0x0f, 0x0a, 0x0b, 0x42, 0x49, 0x47, 0x5f, 0x50, 0x4f, 0x43, 0x4b, 0x45, 0x54, 0x53,
	0x10, 0x17, 0x12, 0x18, 0x0a, 0x14, 0x54, 0x48, 0x49, 0x45, 0x56, 0x45, 0x53, 0x5f, 0x46, 0x41,
	0x56, 0x4f, 0x52, 0x49, 0x54, 0x45, 0x5f, 0x42, 0x41, 0x47, 0x10, 0x18, 0x12, 0x0a, 0x0a, 0x06,
	0x43, 0x4f, 0x46, 0x46, 0x45, 0x45, 0x10, 0x19, 0x12, 0x10, 0x0a, 0x0c, 0x4e, 0x49, 0x47, 0x48,
	0x54, 0x53, 0x5f, 0x57, 0x41, 0x54, 0x43, 0x48, 0x10, 0x1a, 0x12, 0x0f, 0x0a, 0x0b, 0x4c, 0x41,
	0x53, 0x45, 0x52, 0x5f, 0x41, 0x4c, 0x41, 0x52, 0x4d, 0x10, 0x1b, 0x12, 0x0d, 0x0a, 0x09, 0x54,
	0x52, 0x49, 0x50, 0x5f, 0x57, 0x49, 0x52, 0x45, 0x10, 0x1c, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x41,
	0x52, 0x44, 0x49, 0x4f, 0x10, 0x1d, 0x2a, 0x50, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x52, 0x4f, 0x44, 0x55, 0x43,
	0x54, 0x49, 0x4f, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x4d, 0x41, 0x4e, 0x44,
	0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x54, 0x41, 0x50, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12,
	0x0b, 0x0a, 0x07, 0x54, 0x48, 0x49, 0x45, 0x56, 0x45, 0x53, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06,
	0x47, 0x55, 0x41, 0x52, 0x44, 0x53, 0x10, 0x04, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x6e, 0x61, 0x74, 0x74, 0x65, 0x2f, 0x70, 0x69,
	0x7a, 0x7a, 0x61, 0x2d, 0x74, 0x72, 0x69, 0x62, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_research_proto_rawDescOnce sync.Once
	file_research_proto_rawDescData = file_research_proto_rawDesc
)

func file_research_proto_rawDescGZIP() []byte {
	file_research_proto_rawDescOnce.Do(func() {
		file_research_proto_rawDescData = protoimpl.X.CompressGZIP(file_research_proto_rawDescData)
	})
	return file_research_proto_rawDescData
}

var file_research_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_research_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_research_proto_goTypes = []interface{}{
	(ResearchDiscovery)(0),      // 0: pizzatribes.ResearchDiscovery
	(ResearchTree)(0),           // 1: pizzatribes.ResearchTree
	(*ResearchInfo)(nil),        // 2: pizzatribes.ResearchInfo
	(*ResearchInfo_Reward)(nil), // 3: pizzatribes.ResearchInfo.Reward
}
var file_research_proto_depIdxs = []int32{
	0, // 0: pizzatribes.ResearchInfo.discovery:type_name -> pizzatribes.ResearchDiscovery
	0, // 1: pizzatribes.ResearchInfo.requirements:type_name -> pizzatribes.ResearchDiscovery
	1, // 2: pizzatribes.ResearchInfo.tree:type_name -> pizzatribes.ResearchTree
	3, // 3: pizzatribes.ResearchInfo.rewards:type_name -> pizzatribes.ResearchInfo.Reward
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_research_proto_init() }
func file_research_proto_init() {
	if File_research_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_research_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResearchInfo); i {
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
		file_research_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResearchInfo_Reward); i {
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
			RawDescriptor: file_research_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_research_proto_goTypes,
		DependencyIndexes: file_research_proto_depIdxs,
		EnumInfos:         file_research_proto_enumTypes,
		MessageInfos:      file_research_proto_msgTypes,
	}.Build()
	File_research_proto = out.File
	file_research_proto_rawDesc = nil
	file_research_proto_goTypes = nil
	file_research_proto_depIdxs = nil
}