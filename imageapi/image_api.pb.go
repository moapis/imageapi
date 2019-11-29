// Code generated by protoc-gen-go. DO NOT EDIT.
// source: image_api.proto

package moapis_file_api

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ImageDimensions struct {
	Width                uint32   `protobuf:"varint,1,opt,name=width,proto3" json:"width,omitempty"`
	Height               uint32   `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	Alpha                uint32   `protobuf:"varint,3,opt,name=alpha,proto3" json:"alpha,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImageDimensions) Reset()         { *m = ImageDimensions{} }
func (m *ImageDimensions) String() string { return proto.CompactTextString(m) }
func (*ImageDimensions) ProtoMessage()    {}
func (*ImageDimensions) Descriptor() ([]byte, []int) {
	return fileDescriptor_5487c576977b4a3c, []int{0}
}

func (m *ImageDimensions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageDimensions.Unmarshal(m, b)
}
func (m *ImageDimensions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageDimensions.Marshal(b, m, deterministic)
}
func (m *ImageDimensions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageDimensions.Merge(m, src)
}
func (m *ImageDimensions) XXX_Size() int {
	return xxx_messageInfo_ImageDimensions.Size(m)
}
func (m *ImageDimensions) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageDimensions.DiscardUnknown(m)
}

var xxx_messageInfo_ImageDimensions proto.InternalMessageInfo

func (m *ImageDimensions) GetWidth() uint32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *ImageDimensions) GetHeight() uint32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *ImageDimensions) GetAlpha() uint32 {
	if m != nil {
		return m.Alpha
	}
	return 0
}

type NewImageRequest struct {
	Image                [][]byte         `protobuf:"bytes,1,rep,name=image,proto3" json:"image,omitempty"`
	Dimensions           *ImageDimensions `protobuf:"bytes,2,opt,name=dimensions,proto3" json:"dimensions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *NewImageRequest) Reset()         { *m = NewImageRequest{} }
func (m *NewImageRequest) String() string { return proto.CompactTextString(m) }
func (*NewImageRequest) ProtoMessage()    {}
func (*NewImageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5487c576977b4a3c, []int{1}
}

func (m *NewImageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewImageRequest.Unmarshal(m, b)
}
func (m *NewImageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewImageRequest.Marshal(b, m, deterministic)
}
func (m *NewImageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewImageRequest.Merge(m, src)
}
func (m *NewImageRequest) XXX_Size() int {
	return xxx_messageInfo_NewImageRequest.Size(m)
}
func (m *NewImageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NewImageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NewImageRequest proto.InternalMessageInfo

func (m *NewImageRequest) GetImage() [][]byte {
	if m != nil {
		return m.Image
	}
	return nil
}

func (m *NewImageRequest) GetDimensions() *ImageDimensions {
	if m != nil {
		return m.Dimensions
	}
	return nil
}

type NewImageResponse struct {
	Link                 []string                  `protobuf:"bytes,1,rep,name=link,proto3" json:"link,omitempty"`
	Structure            []*NewImageResponseStruct `protobuf:"bytes,2,rep,name=Structure,proto3" json:"Structure,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *NewImageResponse) Reset()         { *m = NewImageResponse{} }
func (m *NewImageResponse) String() string { return proto.CompactTextString(m) }
func (*NewImageResponse) ProtoMessage()    {}
func (*NewImageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5487c576977b4a3c, []int{2}
}

func (m *NewImageResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewImageResponse.Unmarshal(m, b)
}
func (m *NewImageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewImageResponse.Marshal(b, m, deterministic)
}
func (m *NewImageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewImageResponse.Merge(m, src)
}
func (m *NewImageResponse) XXX_Size() int {
	return xxx_messageInfo_NewImageResponse.Size(m)
}
func (m *NewImageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NewImageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NewImageResponse proto.InternalMessageInfo

func (m *NewImageResponse) GetLink() []string {
	if m != nil {
		return m.Link
	}
	return nil
}

func (m *NewImageResponse) GetStructure() []*NewImageResponseStruct {
	if m != nil {
		return m.Structure
	}
	return nil
}

type RemoveImageRequest struct {
	Link                 []string `protobuf:"bytes,1,rep,name=link,proto3" json:"link,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveImageRequest) Reset()         { *m = RemoveImageRequest{} }
func (m *RemoveImageRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveImageRequest) ProtoMessage()    {}
func (*RemoveImageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5487c576977b4a3c, []int{3}
}

func (m *RemoveImageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveImageRequest.Unmarshal(m, b)
}
func (m *RemoveImageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveImageRequest.Marshal(b, m, deterministic)
}
func (m *RemoveImageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveImageRequest.Merge(m, src)
}
func (m *RemoveImageRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveImageRequest.Size(m)
}
func (m *RemoveImageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveImageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveImageRequest proto.InternalMessageInfo

func (m *RemoveImageRequest) GetLink() []string {
	if m != nil {
		return m.Link
	}
	return nil
}

type RemoveImageResponse struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveImageResponse) Reset()         { *m = RemoveImageResponse{} }
func (m *RemoveImageResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveImageResponse) ProtoMessage()    {}
func (*RemoveImageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5487c576977b4a3c, []int{4}
}

func (m *RemoveImageResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveImageResponse.Unmarshal(m, b)
}
func (m *RemoveImageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveImageResponse.Marshal(b, m, deterministic)
}
func (m *RemoveImageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveImageResponse.Merge(m, src)
}
func (m *RemoveImageResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveImageResponse.Size(m)
}
func (m *RemoveImageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveImageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveImageResponse proto.InternalMessageInfo

func (m *RemoveImageResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type NewImageResponseStruct struct {
	OriginalLink         string   `protobuf:"bytes,1,opt,name=originalLink,proto3" json:"originalLink,omitempty"`
	ResizedLink          string   `protobuf:"bytes,2,opt,name=resizedLink,proto3" json:"resizedLink,omitempty"`
	OriginalID           uint32   `protobuf:"varint,3,opt,name=originalID,proto3" json:"originalID,omitempty"`
	ResizedID            uint32   `protobuf:"varint,4,opt,name=resizedID,proto3" json:"resizedID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewImageResponseStruct) Reset()         { *m = NewImageResponseStruct{} }
func (m *NewImageResponseStruct) String() string { return proto.CompactTextString(m) }
func (*NewImageResponseStruct) ProtoMessage()    {}
func (*NewImageResponseStruct) Descriptor() ([]byte, []int) {
	return fileDescriptor_5487c576977b4a3c, []int{5}
}

func (m *NewImageResponseStruct) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewImageResponseStruct.Unmarshal(m, b)
}
func (m *NewImageResponseStruct) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewImageResponseStruct.Marshal(b, m, deterministic)
}
func (m *NewImageResponseStruct) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewImageResponseStruct.Merge(m, src)
}
func (m *NewImageResponseStruct) XXX_Size() int {
	return xxx_messageInfo_NewImageResponseStruct.Size(m)
}
func (m *NewImageResponseStruct) XXX_DiscardUnknown() {
	xxx_messageInfo_NewImageResponseStruct.DiscardUnknown(m)
}

var xxx_messageInfo_NewImageResponseStruct proto.InternalMessageInfo

func (m *NewImageResponseStruct) GetOriginalLink() string {
	if m != nil {
		return m.OriginalLink
	}
	return ""
}

func (m *NewImageResponseStruct) GetResizedLink() string {
	if m != nil {
		return m.ResizedLink
	}
	return ""
}

func (m *NewImageResponseStruct) GetOriginalID() uint32 {
	if m != nil {
		return m.OriginalID
	}
	return 0
}

func (m *NewImageResponseStruct) GetResizedID() uint32 {
	if m != nil {
		return m.ResizedID
	}
	return 0
}

func init() {
	proto.RegisterType((*ImageDimensions)(nil), "moapis.file_api.ImageDimensions")
	proto.RegisterType((*NewImageRequest)(nil), "moapis.file_api.NewImageRequest")
	proto.RegisterType((*NewImageResponse)(nil), "moapis.file_api.NewImageResponse")
	proto.RegisterType((*RemoveImageRequest)(nil), "moapis.file_api.RemoveImageRequest")
	proto.RegisterType((*RemoveImageResponse)(nil), "moapis.file_api.RemoveImageResponse")
	proto.RegisterType((*NewImageResponseStruct)(nil), "moapis.file_api.NewImageResponseStruct")
}

func init() { proto.RegisterFile("image_api.proto", fileDescriptor_5487c576977b4a3c) }

var fileDescriptor_5487c576977b4a3c = []byte{
	// 404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xdd, 0x8a, 0xd3, 0x40,
	0x14, 0xa6, 0xed, 0x5a, 0xc8, 0x49, 0x35, 0x72, 0x94, 0x12, 0x8b, 0x48, 0x8c, 0x82, 0xb9, 0x31,
	0x17, 0xf5, 0x05, 0x14, 0xea, 0x45, 0x41, 0x44, 0x66, 0x59, 0xf4, 0x4e, 0xc7, 0xcd, 0xb1, 0x19,
	0xcc, 0x9f, 0x33, 0x93, 0x2e, 0xec, 0xb3, 0xf8, 0x6a, 0xbe, 0x8b, 0xf4, 0x24, 0x9a, 0x9f, 0x2e,
	0xbb, 0x37, 0xbd, 0xeb, 0xf9, 0xfa, 0xfd, 0xe5, 0xcc, 0x24, 0xe0, 0xa9, 0x5c, 0xee, 0xe8, 0xab,
	0xac, 0x54, 0x5c, 0xe9, 0xd2, 0x96, 0xe8, 0xe5, 0xa5, 0xac, 0x94, 0x89, 0x7f, 0xa8, 0x8c, 0xe1,
	0xf0, 0x02, 0xbc, 0xed, 0x81, 0xb3, 0x51, 0x39, 0x15, 0x46, 0x95, 0x85, 0xc1, 0xc7, 0x70, 0xef,
	0x4a, 0x25, 0x36, 0xf5, 0x27, 0xc1, 0x24, 0xba, 0x2f, 0x9a, 0x01, 0x97, 0x30, 0x4f, 0x49, 0xed,
	0x52, 0xeb, 0x4f, 0x19, 0x6e, 0xa7, 0x03, 0x5b, 0x66, 0x55, 0x2a, 0xfd, 0x59, 0xc3, 0xe6, 0x21,
	0x54, 0xe0, 0x7d, 0xa4, 0x2b, 0x76, 0x16, 0xf4, 0xab, 0x26, 0xc3, 0x44, 0x6e, 0xe3, 0x4f, 0x82,
	0x59, 0xb4, 0x10, 0xcd, 0x80, 0x6f, 0x01, 0x92, 0xff, 0xd1, 0x6c, 0xed, 0xae, 0x83, 0x78, 0xd4,
	0x32, 0x1e, 0x55, 0x14, 0x3d, 0x4d, 0x98, 0xc3, 0xc3, 0x2e, 0xca, 0x54, 0x65, 0x61, 0x08, 0x11,
	0xce, 0x32, 0x55, 0xfc, 0xe4, 0x28, 0x47, 0xf0, 0x6f, 0x7c, 0x0f, 0xce, 0xb9, 0xd5, 0xf5, 0xa5,
	0xad, 0x35, 0xf9, 0xd3, 0x60, 0x16, 0xb9, 0xeb, 0x57, 0x47, 0x41, 0x63, 0xa7, 0x46, 0x21, 0x3a,
	0x65, 0x18, 0x01, 0x0a, 0xca, 0xcb, 0x3d, 0x0d, 0x1e, 0xee, 0x86, 0xc0, 0xf0, 0x35, 0x3c, 0x1a,
	0x30, 0xdb, 0x6e, 0x4b, 0x98, 0x1b, 0x2b, 0x6d, 0x6d, 0x78, 0xbf, 0x8e, 0x68, 0xa7, 0xf0, 0xf7,
	0x04, 0x96, 0x37, 0xc7, 0x63, 0x08, 0x8b, 0x52, 0xab, 0x9d, 0x2a, 0x64, 0xf6, 0xa1, 0x49, 0x39,
	0x08, 0x07, 0x18, 0x06, 0xe0, 0x6a, 0x32, 0xea, 0x9a, 0x12, 0xa6, 0x4c, 0x99, 0xd2, 0x87, 0xf0,
	0x19, 0xc0, 0x3f, 0xc5, 0x76, 0xd3, 0x1e, 0x57, 0x0f, 0xc1, 0xa7, 0xe0, 0xb4, 0xf4, 0xed, 0xc6,
	0x3f, 0xe3, 0xbf, 0x3b, 0x60, 0xfd, 0x67, 0x06, 0x0b, 0xee, 0x76, 0x4e, 0x7a, 0xaf, 0x2e, 0x09,
	0x2f, 0xe0, 0x41, 0xaf, 0xae, 0xba, 0x26, 0x0c, 0x6e, 0x59, 0x27, 0xaf, 0x69, 0xf5, 0xfc, 0xce,
	0x85, 0xe3, 0xe7, 0xee, 0x38, 0x3f, 0x69, 0x32, 0xa4, 0xf7, 0x27, 0x32, 0xfe, 0x06, 0x4f, 0x86,
	0x7d, 0xdf, 0x15, 0xc9, 0x69, 0x13, 0x24, 0xac, 0x46, 0x09, 0xb6, 0xf7, 0x5a, 0x9d, 0x24, 0xe2,
	0x0b, 0xb8, 0xbd, 0x3b, 0x85, 0x2f, 0x8e, 0x14, 0xc7, 0x77, 0x73, 0xf5, 0xf2, 0x76, 0x52, 0xe3,
	0xfc, 0x7d, 0xce, 0x1f, 0x88, 0x37, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x81, 0x56, 0x16, 0x75,
	0x33, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ImageServiceClient is the client API for ImageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ImageServiceClient interface {
	NewImageResize(ctx context.Context, in *NewImageRequest, opts ...grpc.CallOption) (*NewImageResponse, error)
	NewImagePreserve(ctx context.Context, in *NewImageRequest, opts ...grpc.CallOption) (*NewImageResponse, error)
	NewImageResizeAndPreserve(ctx context.Context, in *NewImageRequest, opts ...grpc.CallOption) (*NewImageResponse, error)
	NewImageResizeAtDimensions(ctx context.Context, in *NewImageRequest, opts ...grpc.CallOption) (*NewImageResponse, error)
	RemoveImage(ctx context.Context, in *RemoveImageRequest, opts ...grpc.CallOption) (*RemoveImageResponse, error)
}

type imageServiceClient struct {
	cc *grpc.ClientConn
}

func NewImageServiceClient(cc *grpc.ClientConn) ImageServiceClient {
	return &imageServiceClient{cc}
}

func (c *imageServiceClient) NewImageResize(ctx context.Context, in *NewImageRequest, opts ...grpc.CallOption) (*NewImageResponse, error) {
	out := new(NewImageResponse)
	err := c.cc.Invoke(ctx, "/moapis.file_api.ImageService/NewImageResize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageServiceClient) NewImagePreserve(ctx context.Context, in *NewImageRequest, opts ...grpc.CallOption) (*NewImageResponse, error) {
	out := new(NewImageResponse)
	err := c.cc.Invoke(ctx, "/moapis.file_api.ImageService/NewImagePreserve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageServiceClient) NewImageResizeAndPreserve(ctx context.Context, in *NewImageRequest, opts ...grpc.CallOption) (*NewImageResponse, error) {
	out := new(NewImageResponse)
	err := c.cc.Invoke(ctx, "/moapis.file_api.ImageService/NewImageResizeAndPreserve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageServiceClient) NewImageResizeAtDimensions(ctx context.Context, in *NewImageRequest, opts ...grpc.CallOption) (*NewImageResponse, error) {
	out := new(NewImageResponse)
	err := c.cc.Invoke(ctx, "/moapis.file_api.ImageService/NewImageResizeAtDimensions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageServiceClient) RemoveImage(ctx context.Context, in *RemoveImageRequest, opts ...grpc.CallOption) (*RemoveImageResponse, error) {
	out := new(RemoveImageResponse)
	err := c.cc.Invoke(ctx, "/moapis.file_api.ImageService/RemoveImage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ImageServiceServer is the server API for ImageService service.
type ImageServiceServer interface {
	NewImageResize(context.Context, *NewImageRequest) (*NewImageResponse, error)
	NewImagePreserve(context.Context, *NewImageRequest) (*NewImageResponse, error)
	NewImageResizeAndPreserve(context.Context, *NewImageRequest) (*NewImageResponse, error)
	NewImageResizeAtDimensions(context.Context, *NewImageRequest) (*NewImageResponse, error)
	RemoveImage(context.Context, *RemoveImageRequest) (*RemoveImageResponse, error)
}

// UnimplementedImageServiceServer can be embedded to have forward compatible implementations.
type UnimplementedImageServiceServer struct {
}

func (*UnimplementedImageServiceServer) NewImageResize(ctx context.Context, req *NewImageRequest) (*NewImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewImageResize not implemented")
}
func (*UnimplementedImageServiceServer) NewImagePreserve(ctx context.Context, req *NewImageRequest) (*NewImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewImagePreserve not implemented")
}
func (*UnimplementedImageServiceServer) NewImageResizeAndPreserve(ctx context.Context, req *NewImageRequest) (*NewImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewImageResizeAndPreserve not implemented")
}
func (*UnimplementedImageServiceServer) NewImageResizeAtDimensions(ctx context.Context, req *NewImageRequest) (*NewImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewImageResizeAtDimensions not implemented")
}
func (*UnimplementedImageServiceServer) RemoveImage(ctx context.Context, req *RemoveImageRequest) (*RemoveImageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveImage not implemented")
}

func RegisterImageServiceServer(s *grpc.Server, srv ImageServiceServer) {
	s.RegisterService(&_ImageService_serviceDesc, srv)
}

func _ImageService_NewImageResize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageServiceServer).NewImageResize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/moapis.file_api.ImageService/NewImageResize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageServiceServer).NewImageResize(ctx, req.(*NewImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageService_NewImagePreserve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageServiceServer).NewImagePreserve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/moapis.file_api.ImageService/NewImagePreserve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageServiceServer).NewImagePreserve(ctx, req.(*NewImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageService_NewImageResizeAndPreserve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageServiceServer).NewImageResizeAndPreserve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/moapis.file_api.ImageService/NewImageResizeAndPreserve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageServiceServer).NewImageResizeAndPreserve(ctx, req.(*NewImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageService_NewImageResizeAtDimensions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageServiceServer).NewImageResizeAtDimensions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/moapis.file_api.ImageService/NewImageResizeAtDimensions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageServiceServer).NewImageResizeAtDimensions(ctx, req.(*NewImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageService_RemoveImage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveImageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageServiceServer).RemoveImage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/moapis.file_api.ImageService/RemoveImage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageServiceServer).RemoveImage(ctx, req.(*RemoveImageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ImageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "moapis.file_api.ImageService",
	HandlerType: (*ImageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewImageResize",
			Handler:    _ImageService_NewImageResize_Handler,
		},
		{
			MethodName: "NewImagePreserve",
			Handler:    _ImageService_NewImagePreserve_Handler,
		},
		{
			MethodName: "NewImageResizeAndPreserve",
			Handler:    _ImageService_NewImageResizeAndPreserve_Handler,
		},
		{
			MethodName: "NewImageResizeAtDimensions",
			Handler:    _ImageService_NewImageResizeAtDimensions_Handler,
		},
		{
			MethodName: "RemoveImage",
			Handler:    _ImageService_RemoveImage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "image_api.proto",
}
