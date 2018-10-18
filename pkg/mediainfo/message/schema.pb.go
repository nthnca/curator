// Code generated by protoc-gen-go. DO NOT EDIT.
// source: schema.proto

package message

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FileType int32

const (
	FileType_JPG FileType = 0
	FileType_RAW FileType = 1
)

var FileType_name = map[int32]string{
	0: "JPG",
	1: "RAW",
}

var FileType_value = map[string]int32{
	"JPG": 0,
	"RAW": 1,
}

func (x FileType) String() string {
	return proto.EnumName(FileType_name, int32(x))
}

func (FileType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1c5fb4d8cc22d66a, []int{0}
}

// This represents a single image. It may contain a set of different files for
// example (.jpg, .raw, .mp4, etc) that have different resolutions or image
// quality for the same image.
type Media struct {
	// This is the "key" that you can use to refer to this image.
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// This is the timestamp of when this Media object was last modified. When
	// you modify a media object, don't modify the current one, but just add a
	// new one with a new timestamp.
	TimestampSeconds int64 `protobuf:"varint,2,opt,name=timestamp_seconds,json=timestampSeconds,proto3" json:"timestamp_seconds,omitempty"`
	// Metadata for the image.
	Photo *PhotoInfo `protobuf:"bytes,3,opt,name=photo,proto3" json:"photo,omitempty"`
	// The list of files that represent this image.
	File []*FileInfo `protobuf:"bytes,4,rep,name=file,proto3" json:"file,omitempty"`
	// Normalized name - This is constructed from the image metadata.
	Name string `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	// Allows us to tag images so we can easily sort and organize.
	Tags                 []string `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Media) Reset()         { *m = Media{} }
func (m *Media) String() string { return proto.CompactTextString(m) }
func (*Media) ProtoMessage()    {}
func (*Media) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c5fb4d8cc22d66a, []int{0}
}

func (m *Media) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Media.Unmarshal(m, b)
}
func (m *Media) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Media.Marshal(b, m, deterministic)
}
func (m *Media) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Media.Merge(m, src)
}
func (m *Media) XXX_Size() int {
	return xxx_messageInfo_Media.Size(m)
}
func (m *Media) XXX_DiscardUnknown() {
	xxx_messageInfo_Media.DiscardUnknown(m)
}

var xxx_messageInfo_Media proto.InternalMessageInfo

func (m *Media) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Media) GetTimestampSeconds() int64 {
	if m != nil {
		return m.TimestampSeconds
	}
	return 0
}

func (m *Media) GetPhoto() *PhotoInfo {
	if m != nil {
		return m.Photo
	}
	return nil
}

func (m *Media) GetFile() []*FileInfo {
	if m != nil {
		return m.File
	}
	return nil
}

func (m *Media) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Media) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

// Information about the original media capture, this information shouldn't
// change.
type PhotoInfo struct {
	TimestampSeconds     int64     `protobuf:"varint,1,opt,name=timestamp_seconds,json=timestampSeconds,proto3" json:"timestamp_seconds,omitempty"`
	Datetime             string    `protobuf:"bytes,10,opt,name=datetime,proto3" json:"datetime,omitempty"`
	Make                 string    `protobuf:"bytes,2,opt,name=make,proto3" json:"make,omitempty"`
	Model                string    `protobuf:"bytes,3,opt,name=model,proto3" json:"model,omitempty"`
	Aperture             *Fraction `protobuf:"bytes,4,opt,name=aperture,proto3" json:"aperture,omitempty"`
	ExposureTime         *Fraction `protobuf:"bytes,5,opt,name=exposure_time,json=exposureTime,proto3" json:"exposure_time,omitempty"`
	FocalLength          *Fraction `protobuf:"bytes,6,opt,name=focal_length,json=focalLength,proto3" json:"focal_length,omitempty"`
	Iso                  int32     `protobuf:"varint,7,opt,name=iso,proto3" json:"iso,omitempty"`
	Width                int32     `protobuf:"varint,8,opt,name=width,proto3" json:"width,omitempty"`
	Height               int32     `protobuf:"varint,9,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *PhotoInfo) Reset()         { *m = PhotoInfo{} }
func (m *PhotoInfo) String() string { return proto.CompactTextString(m) }
func (*PhotoInfo) ProtoMessage()    {}
func (*PhotoInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c5fb4d8cc22d66a, []int{1}
}

func (m *PhotoInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PhotoInfo.Unmarshal(m, b)
}
func (m *PhotoInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PhotoInfo.Marshal(b, m, deterministic)
}
func (m *PhotoInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PhotoInfo.Merge(m, src)
}
func (m *PhotoInfo) XXX_Size() int {
	return xxx_messageInfo_PhotoInfo.Size(m)
}
func (m *PhotoInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PhotoInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PhotoInfo proto.InternalMessageInfo

func (m *PhotoInfo) GetTimestampSeconds() int64 {
	if m != nil {
		return m.TimestampSeconds
	}
	return 0
}

func (m *PhotoInfo) GetDatetime() string {
	if m != nil {
		return m.Datetime
	}
	return ""
}

func (m *PhotoInfo) GetMake() string {
	if m != nil {
		return m.Make
	}
	return ""
}

func (m *PhotoInfo) GetModel() string {
	if m != nil {
		return m.Model
	}
	return ""
}

func (m *PhotoInfo) GetAperture() *Fraction {
	if m != nil {
		return m.Aperture
	}
	return nil
}

func (m *PhotoInfo) GetExposureTime() *Fraction {
	if m != nil {
		return m.ExposureTime
	}
	return nil
}

func (m *PhotoInfo) GetFocalLength() *Fraction {
	if m != nil {
		return m.FocalLength
	}
	return nil
}

func (m *PhotoInfo) GetIso() int32 {
	if m != nil {
		return m.Iso
	}
	return 0
}

func (m *PhotoInfo) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *PhotoInfo) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

type FileInfo struct {
	Filename             string   `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	Type                 FileType `protobuf:"varint,2,opt,name=type,proto3,enum=message.FileType" json:"type,omitempty"`
	Md5Sum               []byte   `protobuf:"bytes,3,opt,name=md5sum,proto3" json:"md5sum,omitempty"`
	Sha256Sum            []byte   `protobuf:"bytes,4,opt,name=sha256sum,proto3" json:"sha256sum,omitempty"`
	SizeInBytes          int64    `protobuf:"varint,5,opt,name=size_in_bytes,json=sizeInBytes,proto3" json:"size_in_bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileInfo) Reset()         { *m = FileInfo{} }
func (m *FileInfo) String() string { return proto.CompactTextString(m) }
func (*FileInfo) ProtoMessage()    {}
func (*FileInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c5fb4d8cc22d66a, []int{2}
}

func (m *FileInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileInfo.Unmarshal(m, b)
}
func (m *FileInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileInfo.Marshal(b, m, deterministic)
}
func (m *FileInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileInfo.Merge(m, src)
}
func (m *FileInfo) XXX_Size() int {
	return xxx_messageInfo_FileInfo.Size(m)
}
func (m *FileInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_FileInfo.DiscardUnknown(m)
}

var xxx_messageInfo_FileInfo proto.InternalMessageInfo

func (m *FileInfo) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

func (m *FileInfo) GetType() FileType {
	if m != nil {
		return m.Type
	}
	return FileType_JPG
}

func (m *FileInfo) GetMd5Sum() []byte {
	if m != nil {
		return m.Md5Sum
	}
	return nil
}

func (m *FileInfo) GetSha256Sum() []byte {
	if m != nil {
		return m.Sha256Sum
	}
	return nil
}

func (m *FileInfo) GetSizeInBytes() int64 {
	if m != nil {
		return m.SizeInBytes
	}
	return 0
}

type Fraction struct {
	Numerator            int32    `protobuf:"varint,1,opt,name=numerator,proto3" json:"numerator,omitempty"`
	Denominator          int32    `protobuf:"varint,2,opt,name=denominator,proto3" json:"denominator,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Fraction) Reset()         { *m = Fraction{} }
func (m *Fraction) String() string { return proto.CompactTextString(m) }
func (*Fraction) ProtoMessage()    {}
func (*Fraction) Descriptor() ([]byte, []int) {
	return fileDescriptor_1c5fb4d8cc22d66a, []int{3}
}

func (m *Fraction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Fraction.Unmarshal(m, b)
}
func (m *Fraction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Fraction.Marshal(b, m, deterministic)
}
func (m *Fraction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Fraction.Merge(m, src)
}
func (m *Fraction) XXX_Size() int {
	return xxx_messageInfo_Fraction.Size(m)
}
func (m *Fraction) XXX_DiscardUnknown() {
	xxx_messageInfo_Fraction.DiscardUnknown(m)
}

var xxx_messageInfo_Fraction proto.InternalMessageInfo

func (m *Fraction) GetNumerator() int32 {
	if m != nil {
		return m.Numerator
	}
	return 0
}

func (m *Fraction) GetDenominator() int32 {
	if m != nil {
		return m.Denominator
	}
	return 0
}

func init() {
	proto.RegisterEnum("message.FileType", FileType_name, FileType_value)
	proto.RegisterType((*Media)(nil), "message.Media")
	proto.RegisterType((*PhotoInfo)(nil), "message.PhotoInfo")
	proto.RegisterType((*FileInfo)(nil), "message.FileInfo")
	proto.RegisterType((*Fraction)(nil), "message.Fraction")
}

func init() { proto.RegisterFile("schema.proto", fileDescriptor_1c5fb4d8cc22d66a) }

var fileDescriptor_1c5fb4d8cc22d66a = []byte{
	// 480 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x93, 0xdf, 0x6a, 0x13, 0x41,
	0x14, 0xc6, 0x9d, 0x6e, 0xb6, 0xc9, 0x9e, 0xa4, 0x92, 0x0e, 0x22, 0x83, 0xf4, 0x62, 0x09, 0x08,
	0x8b, 0x62, 0x2f, 0xa2, 0xed, 0xbd, 0x5e, 0x28, 0x2d, 0x0a, 0x65, 0x2c, 0x78, 0x19, 0xa6, 0xd9,
	0x93, 0xec, 0xd0, 0xcc, 0xce, 0xb2, 0x33, 0x41, 0xe3, 0x0b, 0xf9, 0x18, 0x3e, 0x88, 0x2f, 0x23,
	0xe7, 0x6c, 0x92, 0xaa, 0x98, 0xbb, 0x73, 0x7e, 0xdf, 0x37, 0x7b, 0xfe, 0xb1, 0x30, 0x0a, 0xf3,
	0x0a, 0x9d, 0x39, 0x6f, 0x5a, 0x1f, 0xbd, 0xec, 0x3b, 0x0c, 0xc1, 0x2c, 0x71, 0xf2, 0x53, 0x40,
	0xfa, 0x09, 0x4b, 0x6b, 0xe4, 0x18, 0x92, 0x7b, 0xdc, 0x28, 0x91, 0x8b, 0x62, 0xa4, 0x29, 0x94,
	0x2f, 0xe1, 0x34, 0x5a, 0x87, 0x21, 0x1a, 0xd7, 0xcc, 0x02, 0xce, 0x7d, 0x5d, 0x06, 0x75, 0x94,
	0x8b, 0x22, 0xd1, 0xe3, 0xbd, 0xf0, 0xb9, 0xe3, 0xb2, 0x80, 0xb4, 0xa9, 0x7c, 0xf4, 0x2a, 0xc9,
	0x45, 0x31, 0x9c, 0xca, 0xf3, 0x6d, 0x85, 0xf3, 0x1b, 0xa2, 0x57, 0xf5, 0xc2, 0xeb, 0xce, 0x20,
	0x9f, 0x43, 0x6f, 0x61, 0x57, 0xa8, 0x7a, 0x79, 0x52, 0x0c, 0xa7, 0xa7, 0x7b, 0xe3, 0x7b, 0xbb,
	0x42, 0xf6, 0xb1, 0x2c, 0x25, 0xf4, 0x6a, 0xe3, 0x50, 0xa5, 0xb9, 0x28, 0x32, 0xcd, 0x31, 0xb1,
	0x68, 0x96, 0x41, 0xf5, 0xf3, 0x84, 0x18, 0xc5, 0x93, 0x5f, 0x47, 0x90, 0xed, 0x6b, 0xfc, 0xbf,
	0x67, 0x71, 0xa0, 0xe7, 0x67, 0x30, 0x28, 0x4d, 0x44, 0xe2, 0x0a, 0xb8, 0xcc, 0x3e, 0xa7, 0x52,
	0xce, 0xdc, 0x23, 0xcf, 0x9b, 0x69, 0x8e, 0xe5, 0x13, 0x48, 0x9d, 0x2f, 0x71, 0xc5, 0x33, 0x66,
	0xba, 0x4b, 0xe4, 0x2b, 0x18, 0x98, 0x06, 0xdb, 0xb8, 0x6e, 0x69, 0x26, 0xf1, 0xf7, 0x4c, 0xad,
	0x99, 0x47, 0xeb, 0x6b, 0xbd, 0xb7, 0xc8, 0x4b, 0x38, 0xc1, 0x6f, 0x8d, 0x0f, 0xeb, 0x16, 0x67,
	0x5c, 0x39, 0x3d, 0xf4, 0x66, 0xb4, 0xf3, 0xdd, 0x52, 0x43, 0x6f, 0x60, 0xb4, 0xf0, 0x73, 0xb3,
	0x9a, 0xad, 0xb0, 0x5e, 0xc6, 0x4a, 0x1d, 0x1f, 0x7a, 0x36, 0x64, 0xdb, 0x47, 0x76, 0xd1, 0x55,
	0x6d, 0xf0, 0xaa, 0x9f, 0x8b, 0x22, 0xd5, 0x14, 0xd2, 0x10, 0x5f, 0x6d, 0x19, 0x2b, 0x35, 0x60,
	0xd6, 0x25, 0xf2, 0x29, 0x1c, 0x57, 0x68, 0x97, 0x55, 0x54, 0x19, 0xe3, 0x6d, 0x36, 0xf9, 0x21,
	0x60, 0xb0, 0x3b, 0x0c, 0xed, 0x8b, 0x4e, 0xc3, 0x67, 0x11, 0xdd, 0xbe, 0x76, 0x39, 0x5d, 0x35,
	0x6e, 0x9a, 0x6e, 0x5f, 0x8f, 0xff, 0xb9, 0xea, 0xed, 0xa6, 0x41, 0xcd, 0x32, 0xd5, 0x71, 0xe5,
	0x45, 0x58, 0x3b, 0xde, 0xe1, 0x48, 0x6f, 0x33, 0x79, 0x06, 0x59, 0xa8, 0xcc, 0xf4, 0xe2, 0x92,
	0xa4, 0x1e, 0x4b, 0x0f, 0x40, 0x4e, 0xe0, 0x24, 0xd8, 0xef, 0x38, 0xb3, 0xf5, 0xec, 0x6e, 0x13,
	0x31, 0xf0, 0xce, 0x12, 0x3d, 0x24, 0x78, 0x55, 0xbf, 0x23, 0x34, 0xb9, 0x86, 0xc1, 0x6e, 0x05,
	0xf4, 0xb5, 0x7a, 0xed, 0xb0, 0x35, 0xd1, 0xb7, 0xdc, 0x69, 0xaa, 0x1f, 0x80, 0xcc, 0x61, 0x58,
	0x62, 0xed, 0x9d, 0xad, 0x59, 0x3f, 0x62, 0xfd, 0x4f, 0xf4, 0xe2, 0xac, 0x1b, 0x9a, 0xfa, 0x96,
	0x7d, 0x48, 0xae, 0x6f, 0x3e, 0x8c, 0x1f, 0x51, 0xa0, 0xdf, 0x7e, 0x19, 0x8b, 0xbb, 0x63, 0xfe,
	0x87, 0x5e, 0xff, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x1a, 0xf6, 0xdf, 0xda, 0x53, 0x03, 0x00, 0x00,
}
