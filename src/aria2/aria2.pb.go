// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.3
// source: aria2.proto

package aria2

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

type DownloadType int32

const (
	DownloadType_HTTP     DownloadType = 0
	DownloadType_MAGLINK  DownloadType = 1
	DownloadType_TORRENT  DownloadType = 2
	DownloadType_METALINK DownloadType = 3
)

// Enum value maps for DownloadType.
var (
	DownloadType_name = map[int32]string{
		0: "HTTP",
		1: "MAGLINK",
		2: "TORRENT",
		3: "METALINK",
	}
	DownloadType_value = map[string]int32{
		"HTTP":     0,
		"MAGLINK":  1,
		"TORRENT":  2,
		"METALINK": 3,
	}
)

func (x DownloadType) Enum() *DownloadType {
	p := new(DownloadType)
	*p = x
	return p
}

func (x DownloadType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DownloadType) Descriptor() protoreflect.EnumDescriptor {
	return file_aria2_proto_enumTypes[0].Descriptor()
}

func (DownloadType) Type() protoreflect.EnumType {
	return &file_aria2_proto_enumTypes[0]
}

func (x DownloadType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DownloadType.Descriptor instead.
func (DownloadType) EnumDescriptor() ([]byte, []int) {
	return file_aria2_proto_rawDescGZIP(), []int{0}
}

type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code     int32                `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"` //0 for Success, 1 for Received
	Msg      string               `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	FileInfo map[string]*FileInfo `protobuf:"bytes,3,rep,name=FileInfo,proto3" json:"FileInfo,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aria2_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_aria2_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_aria2_proto_rawDescGZIP(), []int{0}
}

func (x *Result) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Result) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *Result) GetFileInfo() map[string]*FileInfo {
	if x != nil {
		return x.FileInfo
	}
	return nil
}

type FileInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GID        string `protobuf:"bytes,1,opt,name=GID,proto3" json:"GID,omitempty"`
	IsFinished bool   `protobuf:"varint,2,opt,name=isFinished,proto3" json:"isFinished,omitempty"`
}

func (x *FileInfo) Reset() {
	*x = FileInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aria2_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileInfo) ProtoMessage() {}

func (x *FileInfo) ProtoReflect() protoreflect.Message {
	mi := &file_aria2_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileInfo.ProtoReflect.Descriptor instead.
func (*FileInfo) Descriptor() ([]byte, []int) {
	return file_aria2_proto_rawDescGZIP(), []int{1}
}

func (x *FileInfo) GetGID() string {
	if x != nil {
		return x.GID
	}
	return ""
}

func (x *FileInfo) GetIsFinished() bool {
	if x != nil {
		return x.IsFinished
	}
	return false
}

type Param struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DownloadInfoList []*DownloadInfo `protobuf:"bytes,1,rep,name=DownloadInfoList,proto3" json:"DownloadInfoList,omitempty"`
	GIDList          []string        `protobuf:"bytes,2,rep,name=GIDList,proto3" json:"GIDList,omitempty"`
}

func (x *Param) Reset() {
	*x = Param{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aria2_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Param) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Param) ProtoMessage() {}

func (x *Param) ProtoReflect() protoreflect.Message {
	mi := &file_aria2_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Param.ProtoReflect.Descriptor instead.
func (*Param) Descriptor() ([]byte, []int) {
	return file_aria2_proto_rawDescGZIP(), []int{2}
}

func (x *Param) GetDownloadInfoList() []*DownloadInfo {
	if x != nil {
		return x.DownloadInfoList
	}
	return nil
}

func (x *Param) GetGIDList() []string {
	if x != nil {
		return x.GIDList
	}
	return nil
}

type DownloadInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DownloadType   DownloadType    `protobuf:"varint,1,opt,name=DownloadType,proto3,enum=aria2.DownloadType" json:"DownloadType,omitempty"`
	URL            string          `protobuf:"bytes,2,opt,name=URL,proto3" json:"URL,omitempty"`
	Destination    string          `protobuf:"bytes,3,opt,name=Destination,proto3" json:"Destination,omitempty"`
	FileName       string          `protobuf:"bytes,4,opt,name=FileName,proto3" json:"FileName,omitempty"`
	DownloadOption *DownloadOption `protobuf:"bytes,5,opt,name=DownloadOption,proto3" json:"DownloadOption,omitempty"`
	Token          string          `protobuf:"bytes,6,opt,name=Token,proto3" json:"Token,omitempty"`
}

func (x *DownloadInfo) Reset() {
	*x = DownloadInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aria2_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadInfo) ProtoMessage() {}

func (x *DownloadInfo) ProtoReflect() protoreflect.Message {
	mi := &file_aria2_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadInfo.ProtoReflect.Descriptor instead.
func (*DownloadInfo) Descriptor() ([]byte, []int) {
	return file_aria2_proto_rawDescGZIP(), []int{3}
}

func (x *DownloadInfo) GetDownloadType() DownloadType {
	if x != nil {
		return x.DownloadType
	}
	return DownloadType_HTTP
}

func (x *DownloadInfo) GetURL() string {
	if x != nil {
		return x.URL
	}
	return ""
}

func (x *DownloadInfo) GetDestination() string {
	if x != nil {
		return x.Destination
	}
	return ""
}

func (x *DownloadInfo) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *DownloadInfo) GetDownloadOption() *DownloadOption {
	if x != nil {
		return x.DownloadOption
	}
	return nil
}

func (x *DownloadInfo) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type DownloadOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WithHeader map[string]string `protobuf:"bytes,1,rep,name=WithHeader,proto3" json:"WithHeader,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *DownloadOption) Reset() {
	*x = DownloadOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_aria2_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadOption) ProtoMessage() {}

func (x *DownloadOption) ProtoReflect() protoreflect.Message {
	mi := &file_aria2_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadOption.ProtoReflect.Descriptor instead.
func (*DownloadOption) Descriptor() ([]byte, []int) {
	return file_aria2_proto_rawDescGZIP(), []int{4}
}

func (x *DownloadOption) GetWithHeader() map[string]string {
	if x != nil {
		return x.WithHeader
	}
	return nil
}

var File_aria2_proto protoreflect.FileDescriptor

var file_aria2_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61,
	0x72, 0x69, 0x61, 0x32, 0x22, 0xb5, 0x01, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x4d, 0x73, 0x67, 0x12, 0x37, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x4c,
	0x0a, 0x0d, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x25, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x3c, 0x0a, 0x08,
	0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x47, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x47, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x73,
	0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a,
	0x69, 0x73, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x22, 0x62, 0x0a, 0x05, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x12, 0x3f, 0x0a, 0x10, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x49,
	0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x61, 0x72, 0x69, 0x61, 0x32, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x10, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6e, 0x66, 0x6f,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x47, 0x49, 0x44, 0x4c, 0x69, 0x73, 0x74, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x47, 0x49, 0x44, 0x4c, 0x69, 0x73, 0x74, 0x22, 0xec,
	0x01, 0x0a, 0x0c, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x37, 0x0a, 0x0c, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e, 0x44, 0x6f,
	0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x44, 0x6f, 0x77, 0x6e,
	0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x52, 0x4c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x52, 0x4c, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08,
	0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x46, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x0e, 0x44, 0x6f, 0x77, 0x6e,
	0x6c, 0x6f, 0x61, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x96, 0x01,
	0x0a, 0x0e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x45, 0x0a, 0x0a, 0x57, 0x69, 0x74, 0x68, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x57, 0x69, 0x74, 0x68,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x57, 0x69, 0x74,
	0x68, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x1a, 0x3d, 0x0a, 0x0f, 0x57, 0x69, 0x74, 0x68, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x40, 0x0a, 0x0c, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f,
	0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x54, 0x54, 0x50, 0x10, 0x00,
	0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x41, 0x47, 0x4c, 0x49, 0x4e, 0x4b, 0x10, 0x01, 0x12, 0x0b, 0x0a,
	0x07, 0x54, 0x4f, 0x52, 0x52, 0x45, 0x4e, 0x54, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x4d, 0x45,
	0x54, 0x41, 0x4c, 0x49, 0x4e, 0x4b, 0x10, 0x03, 0x32, 0x6a, 0x0a, 0x0a, 0x41, 0x72, 0x69, 0x61,
	0x32, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x0d, 0x41, 0x77, 0x61, 0x69, 0x74, 0x44,
	0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x0c, 0x2e, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x1a, 0x0d, 0x2e, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x30, 0x01, 0x12, 0x2c, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x44,
	0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x0c, 0x2e, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x1a, 0x0d, 0x2e, 0x61, 0x72, 0x69, 0x61, 0x32, 0x2e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x61, 0x72, 0x69, 0x61, 0x32, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_aria2_proto_rawDescOnce sync.Once
	file_aria2_proto_rawDescData = file_aria2_proto_rawDesc
)

func file_aria2_proto_rawDescGZIP() []byte {
	file_aria2_proto_rawDescOnce.Do(func() {
		file_aria2_proto_rawDescData = protoimpl.X.CompressGZIP(file_aria2_proto_rawDescData)
	})
	return file_aria2_proto_rawDescData
}

var file_aria2_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_aria2_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_aria2_proto_goTypes = []interface{}{
	(DownloadType)(0),      // 0: aria2.DownloadType
	(*Result)(nil),         // 1: aria2.Result
	(*FileInfo)(nil),       // 2: aria2.FileInfo
	(*Param)(nil),          // 3: aria2.Param
	(*DownloadInfo)(nil),   // 4: aria2.DownloadInfo
	(*DownloadOption)(nil), // 5: aria2.DownloadOption
	nil,                    // 6: aria2.Result.FileInfoEntry
	nil,                    // 7: aria2.DownloadOption.WithHeaderEntry
}
var file_aria2_proto_depIdxs = []int32{
	6, // 0: aria2.Result.FileInfo:type_name -> aria2.Result.FileInfoEntry
	4, // 1: aria2.Param.DownloadInfoList:type_name -> aria2.DownloadInfo
	0, // 2: aria2.DownloadInfo.DownloadType:type_name -> aria2.DownloadType
	5, // 3: aria2.DownloadInfo.DownloadOption:type_name -> aria2.DownloadOption
	7, // 4: aria2.DownloadOption.WithHeader:type_name -> aria2.DownloadOption.WithHeaderEntry
	2, // 5: aria2.Result.FileInfoEntry.value:type_name -> aria2.FileInfo
	3, // 6: aria2.Aria2Agent.AwaitDownload:input_type -> aria2.Param
	3, // 7: aria2.Aria2Agent.CheckDownload:input_type -> aria2.Param
	1, // 8: aria2.Aria2Agent.AwaitDownload:output_type -> aria2.Result
	1, // 9: aria2.Aria2Agent.CheckDownload:output_type -> aria2.Result
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_aria2_proto_init() }
func file_aria2_proto_init() {
	if File_aria2_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_aria2_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
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
		file_aria2_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileInfo); i {
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
		file_aria2_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Param); i {
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
		file_aria2_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadInfo); i {
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
		file_aria2_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadOption); i {
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
			RawDescriptor: file_aria2_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_aria2_proto_goTypes,
		DependencyIndexes: file_aria2_proto_depIdxs,
		EnumInfos:         file_aria2_proto_enumTypes,
		MessageInfos:      file_aria2_proto_msgTypes,
	}.Build()
	File_aria2_proto = out.File
	file_aria2_proto_rawDesc = nil
	file_aria2_proto_goTypes = nil
	file_aria2_proto_depIdxs = nil
}
