// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.2
// source: pb/tzinfo.proto

package pb

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

type Point struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lng float32 `protobuf:"fixed32,1,opt,name=lng,proto3" json:"lng,omitempty"`
	Lat float32 `protobuf:"fixed32,2,opt,name=lat,proto3" json:"lat,omitempty"`
}

func (x *Point) Reset() {
	*x = Point{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_tzinfo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Point) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Point) ProtoMessage() {}

func (x *Point) ProtoReflect() protoreflect.Message {
	mi := &file_pb_tzinfo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Point.ProtoReflect.Descriptor instead.
func (*Point) Descriptor() ([]byte, []int) {
	return file_pb_tzinfo_proto_rawDescGZIP(), []int{0}
}

func (x *Point) GetLng() float32 {
	if x != nil {
		return x.Lng
	}
	return 0
}

func (x *Point) GetLat() float32 {
	if x != nil {
		return x.Lat
	}
	return 0
}

type Polygon struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Points []*Point   `protobuf:"bytes,1,rep,name=points,proto3" json:"points,omitempty"`
	Holes  []*Polygon `protobuf:"bytes,2,rep,name=holes,proto3" json:"holes,omitempty"`
}

func (x *Polygon) Reset() {
	*x = Polygon{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_tzinfo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Polygon) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Polygon) ProtoMessage() {}

func (x *Polygon) ProtoReflect() protoreflect.Message {
	mi := &file_pb_tzinfo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Polygon.ProtoReflect.Descriptor instead.
func (*Polygon) Descriptor() ([]byte, []int) {
	return file_pb_tzinfo_proto_rawDescGZIP(), []int{1}
}

func (x *Polygon) GetPoints() []*Point {
	if x != nil {
		return x.Points
	}
	return nil
}

func (x *Polygon) GetHoles() []*Polygon {
	if x != nil {
		return x.Holes
	}
	return nil
}

type Timezone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Polygons []*Polygon `protobuf:"bytes,1,rep,name=polygons,proto3" json:"polygons,omitempty"`
	Name     string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Timezone) Reset() {
	*x = Timezone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_tzinfo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Timezone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Timezone) ProtoMessage() {}

func (x *Timezone) ProtoReflect() protoreflect.Message {
	mi := &file_pb_tzinfo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Timezone.ProtoReflect.Descriptor instead.
func (*Timezone) Descriptor() ([]byte, []int) {
	return file_pb_tzinfo_proto_rawDescGZIP(), []int{2}
}

func (x *Timezone) GetPolygons() []*Polygon {
	if x != nil {
		return x.Polygons
	}
	return nil
}

func (x *Timezone) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Timezones struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timezones []*Timezone `protobuf:"bytes,1,rep,name=timezones,proto3" json:"timezones,omitempty"`
}

func (x *Timezones) Reset() {
	*x = Timezones{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_tzinfo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Timezones) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Timezones) ProtoMessage() {}

func (x *Timezones) ProtoReflect() protoreflect.Message {
	mi := &file_pb_tzinfo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Timezones.ProtoReflect.Descriptor instead.
func (*Timezones) Descriptor() ([]byte, []int) {
	return file_pb_tzinfo_proto_rawDescGZIP(), []int{3}
}

func (x *Timezones) GetTimezones() []*Timezone {
	if x != nil {
		return x.Timezones
	}
	return nil
}

var File_pb_tzinfo_proto protoreflect.FileDescriptor

var file_pb_tzinfo_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x62, 0x2f, 0x74, 0x7a, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x2b, 0x0a, 0x05, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x6c, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x6c, 0x6e, 0x67,
	0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x6c,
	0x61, 0x74, 0x22, 0x4f, 0x0a, 0x07, 0x50, 0x6f, 0x6c, 0x79, 0x67, 0x6f, 0x6e, 0x12, 0x21, 0x0a,
	0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e,
	0x70, 0x62, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x12, 0x21, 0x0a, 0x05, 0x68, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x6f, 0x6c, 0x79, 0x67, 0x6f, 0x6e, 0x52, 0x05, 0x68, 0x6f,
	0x6c, 0x65, 0x73, 0x22, 0x47, 0x0a, 0x08, 0x54, 0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x12,
	0x27, 0x0a, 0x08, 0x70, 0x6f, 0x6c, 0x79, 0x67, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x6f, 0x6c, 0x79, 0x67, 0x6f, 0x6e, 0x52, 0x08,
	0x70, 0x6f, 0x6c, 0x79, 0x67, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x37, 0x0a, 0x09,
	0x54, 0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70,
	0x62, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x7a, 0x6f, 0x6e, 0x65, 0x73, 0x42, 0x21, 0x5a, 0x1f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x69, 0x6e, 0x67, 0x73, 0x61, 0x74, 0x75, 0x72, 0x6e, 0x2f, 0x74,
	0x7a, 0x66, 0x2f, 0x70, 0x62, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_tzinfo_proto_rawDescOnce sync.Once
	file_pb_tzinfo_proto_rawDescData = file_pb_tzinfo_proto_rawDesc
)

func file_pb_tzinfo_proto_rawDescGZIP() []byte {
	file_pb_tzinfo_proto_rawDescOnce.Do(func() {
		file_pb_tzinfo_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_tzinfo_proto_rawDescData)
	})
	return file_pb_tzinfo_proto_rawDescData
}

var file_pb_tzinfo_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pb_tzinfo_proto_goTypes = []interface{}{
	(*Point)(nil),     // 0: pb.Point
	(*Polygon)(nil),   // 1: pb.Polygon
	(*Timezone)(nil),  // 2: pb.Timezone
	(*Timezones)(nil), // 3: pb.Timezones
}
var file_pb_tzinfo_proto_depIdxs = []int32{
	0, // 0: pb.Polygon.points:type_name -> pb.Point
	1, // 1: pb.Polygon.holes:type_name -> pb.Polygon
	1, // 2: pb.Timezone.polygons:type_name -> pb.Polygon
	2, // 3: pb.Timezones.timezones:type_name -> pb.Timezone
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pb_tzinfo_proto_init() }
func file_pb_tzinfo_proto_init() {
	if File_pb_tzinfo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_tzinfo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Point); i {
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
		file_pb_tzinfo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Polygon); i {
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
		file_pb_tzinfo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Timezone); i {
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
		file_pb_tzinfo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Timezones); i {
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
			RawDescriptor: file_pb_tzinfo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_tzinfo_proto_goTypes,
		DependencyIndexes: file_pb_tzinfo_proto_depIdxs,
		MessageInfos:      file_pb_tzinfo_proto_msgTypes,
	}.Build()
	File_pb_tzinfo_proto = out.File
	file_pb_tzinfo_proto_rawDesc = nil
	file_pb_tzinfo_proto_goTypes = nil
	file_pb_tzinfo_proto_depIdxs = nil
}
