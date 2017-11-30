// Code generated by protoc-gen-go. DO NOT EDIT.
// source: schema.proto

/*
Package message is a generated protocol buffer package.

It is generated from these files:
	schema.proto

It has these top-level messages:
	Fraction
	Photo
	PhotoSet
	ComparisonEntry
	Comparison
	ComparisonSet
*/
package message

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Fraction struct {
	Numerator   int32 `protobuf:"varint,1,opt,name=numerator" json:"numerator,omitempty"`
	Denominator int32 `protobuf:"varint,2,opt,name=denominator" json:"denominator,omitempty"`
}

func (m *Fraction) Reset()                    { *m = Fraction{} }
func (m *Fraction) String() string            { return proto.CompactTextString(m) }
func (*Fraction) ProtoMessage()               {}
func (*Fraction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

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

type Photo struct {
	Key        string                 `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Path       string                 `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	NumBytes   int64                  `protobuf:"varint,3,opt,name=num_bytes,json=numBytes" json:"num_bytes,omitempty"`
	UserHide   bool                   `protobuf:"varint,5,opt,name=user_hide,json=userHide" json:"user_hide,omitempty"`
	Md5Sum     []byte                 `protobuf:"bytes,6,opt,name=md5sum,proto3" json:"md5sum,omitempty"`
	Sha256Sum  []byte                 `protobuf:"bytes,7,opt,name=sha256sum,proto3" json:"sha256sum,omitempty"`
	Properties *Photo_PhotoProperties `protobuf:"bytes,4,opt,name=properties" json:"properties,omitempty"`
}

func (m *Photo) Reset()                    { *m = Photo{} }
func (m *Photo) String() string            { return proto.CompactTextString(m) }
func (*Photo) ProtoMessage()               {}
func (*Photo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Photo) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Photo) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Photo) GetNumBytes() int64 {
	if m != nil {
		return m.NumBytes
	}
	return 0
}

func (m *Photo) GetUserHide() bool {
	if m != nil {
		return m.UserHide
	}
	return false
}

func (m *Photo) GetMd5Sum() []byte {
	if m != nil {
		return m.Md5Sum
	}
	return nil
}

func (m *Photo) GetSha256Sum() []byte {
	if m != nil {
		return m.Sha256Sum
	}
	return nil
}

func (m *Photo) GetProperties() *Photo_PhotoProperties {
	if m != nil {
		return m.Properties
	}
	return nil
}

type Photo_PhotoProperties struct {
	EpochInSeconds int64     `protobuf:"varint,1,opt,name=epoch_in_seconds,json=epochInSeconds" json:"epoch_in_seconds,omitempty"`
	Width          int32     `protobuf:"varint,2,opt,name=width" json:"width,omitempty"`
	Height         int32     `protobuf:"varint,3,opt,name=height" json:"height,omitempty"`
	Make           string    `protobuf:"bytes,4,opt,name=make" json:"make,omitempty"`
	Model          string    `protobuf:"bytes,5,opt,name=model" json:"model,omitempty"`
	Aperture       *Fraction `protobuf:"bytes,6,opt,name=aperture" json:"aperture,omitempty"`
	ExposureTime   *Fraction `protobuf:"bytes,7,opt,name=exposure_time,json=exposureTime" json:"exposure_time,omitempty"`
	FocalLength    *Fraction `protobuf:"bytes,8,opt,name=focal_length,json=focalLength" json:"focal_length,omitempty"`
	Iso            int32     `protobuf:"varint,9,opt,name=iso" json:"iso,omitempty"`
}

func (m *Photo_PhotoProperties) Reset()                    { *m = Photo_PhotoProperties{} }
func (m *Photo_PhotoProperties) String() string            { return proto.CompactTextString(m) }
func (*Photo_PhotoProperties) ProtoMessage()               {}
func (*Photo_PhotoProperties) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *Photo_PhotoProperties) GetEpochInSeconds() int64 {
	if m != nil {
		return m.EpochInSeconds
	}
	return 0
}

func (m *Photo_PhotoProperties) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *Photo_PhotoProperties) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Photo_PhotoProperties) GetMake() string {
	if m != nil {
		return m.Make
	}
	return ""
}

func (m *Photo_PhotoProperties) GetModel() string {
	if m != nil {
		return m.Model
	}
	return ""
}

func (m *Photo_PhotoProperties) GetAperture() *Fraction {
	if m != nil {
		return m.Aperture
	}
	return nil
}

func (m *Photo_PhotoProperties) GetExposureTime() *Fraction {
	if m != nil {
		return m.ExposureTime
	}
	return nil
}

func (m *Photo_PhotoProperties) GetFocalLength() *Fraction {
	if m != nil {
		return m.FocalLength
	}
	return nil
}

func (m *Photo_PhotoProperties) GetIso() int32 {
	if m != nil {
		return m.Iso
	}
	return 0
}

type PhotoSet struct {
	Photo []*Photo `protobuf:"bytes,1,rep,name=photo" json:"photo,omitempty"`
}

func (m *PhotoSet) Reset()                    { *m = PhotoSet{} }
func (m *PhotoSet) String() string            { return proto.CompactTextString(m) }
func (*PhotoSet) ProtoMessage()               {}
func (*PhotoSet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PhotoSet) GetPhoto() []*Photo {
	if m != nil {
		return m.Photo
	}
	return nil
}

type ComparisonEntry struct {
	Photo1 string `protobuf:"bytes,1,opt,name=photo1" json:"photo1,omitempty"`
	Photo2 string `protobuf:"bytes,2,opt,name=photo2" json:"photo2,omitempty"`
	// Should only be 1 or 2.
	Score int32 `protobuf:"varint,3,opt,name=score" json:"score,omitempty"`
}

func (m *ComparisonEntry) Reset()                    { *m = ComparisonEntry{} }
func (m *ComparisonEntry) String() string            { return proto.CompactTextString(m) }
func (*ComparisonEntry) ProtoMessage()               {}
func (*ComparisonEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ComparisonEntry) GetPhoto1() string {
	if m != nil {
		return m.Photo1
	}
	return ""
}

func (m *ComparisonEntry) GetPhoto2() string {
	if m != nil {
		return m.Photo2
	}
	return ""
}

func (m *ComparisonEntry) GetScore() int32 {
	if m != nil {
		return m.Score
	}
	return 0
}

type Comparison struct {
	Epoch int64              `protobuf:"varint,1,opt,name=epoch" json:"epoch,omitempty"`
	Entry []*ComparisonEntry `protobuf:"bytes,2,rep,name=entry" json:"entry,omitempty"`
}

func (m *Comparison) Reset()                    { *m = Comparison{} }
func (m *Comparison) String() string            { return proto.CompactTextString(m) }
func (*Comparison) ProtoMessage()               {}
func (*Comparison) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Comparison) GetEpoch() int64 {
	if m != nil {
		return m.Epoch
	}
	return 0
}

func (m *Comparison) GetEntry() []*ComparisonEntry {
	if m != nil {
		return m.Entry
	}
	return nil
}

type ComparisonSet struct {
	Comparison []*Comparison `protobuf:"bytes,1,rep,name=comparison" json:"comparison,omitempty"`
}

func (m *ComparisonSet) Reset()                    { *m = ComparisonSet{} }
func (m *ComparisonSet) String() string            { return proto.CompactTextString(m) }
func (*ComparisonSet) ProtoMessage()               {}
func (*ComparisonSet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ComparisonSet) GetComparison() []*Comparison {
	if m != nil {
		return m.Comparison
	}
	return nil
}

func init() {
	proto.RegisterType((*Fraction)(nil), "message.Fraction")
	proto.RegisterType((*Photo)(nil), "message.Photo")
	proto.RegisterType((*Photo_PhotoProperties)(nil), "message.Photo.PhotoProperties")
	proto.RegisterType((*PhotoSet)(nil), "message.PhotoSet")
	proto.RegisterType((*ComparisonEntry)(nil), "message.ComparisonEntry")
	proto.RegisterType((*Comparison)(nil), "message.Comparison")
	proto.RegisterType((*ComparisonSet)(nil), "message.ComparisonSet")
}

func init() { proto.RegisterFile("schema.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 512 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0x4d, 0x6f, 0xd3, 0x4c,
	0x10, 0x96, 0x9b, 0x38, 0xb5, 0x27, 0xe9, 0xc7, 0xbb, 0x2f, 0x42, 0x16, 0x20, 0x64, 0x59, 0x1c,
	0x7c, 0x21, 0x82, 0x94, 0xf6, 0xc8, 0x81, 0x2f, 0x01, 0xe2, 0x50, 0x6d, 0x91, 0x38, 0x5a, 0xae,
	0x3d, 0xc4, 0xab, 0x66, 0xbd, 0xd6, 0xee, 0x5a, 0x90, 0xbf, 0xc3, 0x4f, 0xe1, 0x97, 0xa1, 0x1d,
	0x6f, 0x9d, 0x50, 0x95, 0x4b, 0x34, 0xcf, 0xc7, 0x8e, 0x67, 0x9e, 0xdd, 0xc0, 0xc2, 0x54, 0x0d,
	0xca, 0x72, 0xd9, 0x69, 0x65, 0x15, 0x3b, 0x94, 0x68, 0x4c, 0xb9, 0xc6, 0xec, 0x33, 0x44, 0x1f,
	0x74, 0x59, 0x59, 0xa1, 0x5a, 0xf6, 0x04, 0xe2, 0xb6, 0x97, 0xa8, 0x4b, 0xab, 0x74, 0x12, 0xa4,
	0x41, 0x1e, 0xf2, 0x1d, 0xc1, 0x52, 0x98, 0xd7, 0xd8, 0x2a, 0x29, 0x5a, 0xd2, 0x0f, 0x48, 0xdf,
	0xa7, 0xb2, 0x5f, 0x53, 0x08, 0x2f, 0x1b, 0xd7, 0xfe, 0x14, 0x26, 0x37, 0xb8, 0xa5, 0x1e, 0x31,
	0x77, 0x25, 0x63, 0x30, 0xed, 0x4a, 0xdb, 0xd0, 0xb1, 0x98, 0x53, 0xcd, 0x1e, 0xd3, 0xf7, 0x8a,
	0xeb, 0xad, 0x45, 0x93, 0x4c, 0xd2, 0x20, 0x9f, 0xf0, 0xa8, 0xed, 0xe5, 0x1b, 0x87, 0x9d, 0xd8,
	0x1b, 0xd4, 0x45, 0x23, 0x6a, 0x4c, 0xc2, 0x34, 0xc8, 0x23, 0x1e, 0x39, 0xe2, 0xa3, 0xa8, 0x91,
	0x3d, 0x84, 0x99, 0xac, 0xcf, 0x4d, 0x2f, 0x93, 0x59, 0x1a, 0xe4, 0x0b, 0xee, 0x91, 0xdb, 0xc0,
	0x34, 0xe5, 0xea, 0xfc, 0xc2, 0x49, 0x87, 0x24, 0xed, 0x08, 0xf6, 0x1a, 0xa0, 0xd3, 0xaa, 0x43,
	0x6d, 0x05, 0x9a, 0x64, 0x9a, 0x06, 0xf9, 0x7c, 0xf5, 0x74, 0xe9, 0x93, 0x58, 0xd2, 0xe4, 0xc3,
	0xef, 0xe5, 0xe8, 0xe2, 0x7b, 0x27, 0x1e, 0xfd, 0x3e, 0x80, 0x93, 0x3b, 0x3a, 0xcb, 0xe1, 0x14,
	0x3b, 0x55, 0x35, 0x85, 0x68, 0x0b, 0x83, 0x95, 0x6a, 0x6b, 0x43, 0x6b, 0x4f, 0xf8, 0x31, 0xf1,
	0x9f, 0xda, 0xab, 0x81, 0x65, 0x0f, 0x20, 0xfc, 0x21, 0x6a, 0x1f, 0x41, 0xc8, 0x07, 0xe0, 0x36,
	0x69, 0x50, 0xac, 0x1b, 0x4b, 0x01, 0x84, 0xdc, 0x23, 0x97, 0x97, 0x2c, 0x6f, 0x90, 0xa6, 0x8c,
	0x39, 0xd5, 0xae, 0x83, 0x54, 0x35, 0x6e, 0x28, 0x8e, 0x98, 0x0f, 0x80, 0x3d, 0x87, 0xa8, 0x74,
	0xd3, 0xf4, 0x1a, 0x29, 0x8d, 0xf9, 0xea, 0xbf, 0x71, 0xa7, 0xdb, 0xab, 0xe5, 0xa3, 0x85, 0x5d,
	0xc0, 0x11, 0xfe, 0xec, 0x94, 0xe9, 0x35, 0x16, 0x56, 0x48, 0xa4, 0x98, 0xee, 0x3d, 0xb3, 0xb8,
	0xf5, 0x7d, 0x15, 0x12, 0xd9, 0x2b, 0x58, 0x7c, 0x57, 0x55, 0xb9, 0x29, 0x36, 0xd8, 0xae, 0x6d,
	0x93, 0x44, 0xff, 0x3a, 0x36, 0x27, 0xdb, 0x17, 0x72, 0xb9, 0x87, 0x20, 0x8c, 0x4a, 0x62, 0xda,
	0xcd, 0x95, 0xd9, 0x0b, 0x88, 0x28, 0xc3, 0x2b, 0xb4, 0xec, 0x19, 0x84, 0x9d, 0xab, 0x93, 0x20,
	0x9d, 0xe4, 0xf3, 0xd5, 0xf1, 0xdf, 0x77, 0xc1, 0x07, 0x31, 0xfb, 0x06, 0x27, 0x6f, 0x95, 0xec,
	0x4a, 0x2d, 0x8c, 0x6a, 0xdf, 0xb7, 0x56, 0x6f, 0x5d, 0x6a, 0xa4, 0xbd, 0xf4, 0x4f, 0xcc, 0xa3,
	0x91, 0x5f, 0xf9, 0x77, 0xe6, 0x91, 0x4b, 0xce, 0x54, 0x4a, 0xa3, 0x0f, 0x79, 0x00, 0x19, 0x07,
	0xd8, 0x35, 0x76, 0x1e, 0xba, 0x31, 0x7f, 0x7d, 0x03, 0x60, 0x4b, 0x08, 0xd1, 0x7d, 0x32, 0x39,
	0xa0, 0x11, 0x93, 0x71, 0xc4, 0x3b, 0x23, 0xf1, 0xc1, 0x96, 0xbd, 0x83, 0xa3, 0x9d, 0xe2, 0x76,
	0x3c, 0x03, 0xa8, 0x46, 0xc2, 0x2f, 0xfa, 0xff, 0x3d, 0x5d, 0xf8, 0x9e, 0xed, 0x7a, 0x46, 0xff,
	0xd2, 0xb3, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1e, 0x00, 0x57, 0x5b, 0xb5, 0x03, 0x00, 0x00,
}
