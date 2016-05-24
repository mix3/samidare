// Code generated by protoc-gen-go.
// source: samidare.proto
// DO NOT EDIT!

/*
Package samidare is a generated protocol buffer package.

It is generated from these files:
	samidare.proto

It has these top-level messages:
	WatchRequest
	WatchResponse
*/
package samidare

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type WatchRequest_Type int32

const (
	WatchRequest_INCREMENT WatchRequest_Type = 0
	WatchRequest_DECREMENT WatchRequest_Type = 1
)

var WatchRequest_Type_name = map[int32]string{
	0: "INCREMENT",
	1: "DECREMENT",
}
var WatchRequest_Type_value = map[string]int32{
	"INCREMENT": 0,
	"DECREMENT": 1,
}

func (x WatchRequest_Type) String() string {
	return proto.EnumName(WatchRequest_Type_name, int32(x))
}
func (WatchRequest_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type WatchResponse_Type int32

const (
	WatchResponse_DESIRED WatchResponse_Type = 0
	WatchResponse_CURRENT WatchResponse_Type = 1
)

var WatchResponse_Type_name = map[int32]string{
	0: "DESIRED",
	1: "CURRENT",
}
var WatchResponse_Type_value = map[string]int32{
	"DESIRED": 0,
	"CURRENT": 1,
}

func (x WatchResponse_Type) String() string {
	return proto.EnumName(WatchResponse_Type_name, int32(x))
}
func (WatchResponse_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type WatchRequest struct {
	Type WatchRequest_Type `protobuf:"varint,1,opt,name=type,enum=samidare.WatchRequest_Type" json:"type,omitempty"`
}

func (m *WatchRequest) Reset()                    { *m = WatchRequest{} }
func (m *WatchRequest) String() string            { return proto.CompactTextString(m) }
func (*WatchRequest) ProtoMessage()               {}
func (*WatchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type WatchResponse struct {
	Type WatchResponse_Type `protobuf:"varint,1,opt,name=type,enum=samidare.WatchResponse_Type" json:"type,omitempty"`
	Num  int64              `protobuf:"varint,2,opt,name=num" json:"num,omitempty"`
}

func (m *WatchResponse) Reset()                    { *m = WatchResponse{} }
func (m *WatchResponse) String() string            { return proto.CompactTextString(m) }
func (*WatchResponse) ProtoMessage()               {}
func (*WatchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*WatchRequest)(nil), "samidare.WatchRequest")
	proto.RegisterType((*WatchResponse)(nil), "samidare.WatchResponse")
	proto.RegisterEnum("samidare.WatchRequest_Type", WatchRequest_Type_name, WatchRequest_Type_value)
	proto.RegisterEnum("samidare.WatchResponse_Type", WatchResponse_Type_name, WatchResponse_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for Samidare service

type SamidareClient interface {
	Watch(ctx context.Context, opts ...grpc.CallOption) (Samidare_WatchClient, error)
}

type samidareClient struct {
	cc *grpc.ClientConn
}

func NewSamidareClient(cc *grpc.ClientConn) SamidareClient {
	return &samidareClient{cc}
}

func (c *samidareClient) Watch(ctx context.Context, opts ...grpc.CallOption) (Samidare_WatchClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Samidare_serviceDesc.Streams[0], c.cc, "/samidare.Samidare/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &samidareWatchClient{stream}
	return x, nil
}

type Samidare_WatchClient interface {
	Send(*WatchRequest) error
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type samidareWatchClient struct {
	grpc.ClientStream
}

func (x *samidareWatchClient) Send(m *WatchRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *samidareWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Samidare service

type SamidareServer interface {
	Watch(Samidare_WatchServer) error
}

func RegisterSamidareServer(s *grpc.Server, srv SamidareServer) {
	s.RegisterService(&_Samidare_serviceDesc, srv)
}

func _Samidare_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SamidareServer).Watch(&samidareWatchServer{stream})
}

type Samidare_WatchServer interface {
	Send(*WatchResponse) error
	Recv() (*WatchRequest, error)
	grpc.ServerStream
}

type samidareWatchServer struct {
	grpc.ServerStream
}

func (x *samidareWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *samidareWatchServer) Recv() (*WatchRequest, error) {
	m := new(WatchRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Samidare_serviceDesc = grpc.ServiceDesc{
	ServiceName: "samidare.Samidare",
	HandlerType: (*SamidareServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _Samidare_Watch_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
}

var fileDescriptor0 = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x4e, 0xcc, 0xcd,
	0x4c, 0x49, 0x2c, 0x4a, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x80, 0xf1, 0x95, 0xe2,
	0xb9, 0x78, 0xc2, 0x13, 0x4b, 0x92, 0x33, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x34,
	0xb9, 0x58, 0x4a, 0x2a, 0x0b, 0x52, 0x25, 0x18, 0x15, 0x18, 0x35, 0xf8, 0x8c, 0xa4, 0xf5, 0xe0,
	0x1a, 0x91, 0x55, 0xe9, 0x85, 0x00, 0x95, 0x28, 0xa9, 0x70, 0xb1, 0x80, 0x68, 0x21, 0x5e, 0x2e,
	0x4e, 0x4f, 0x3f, 0xe7, 0x20, 0x57, 0x5f, 0x57, 0xbf, 0x10, 0x01, 0x06, 0x10, 0xd7, 0xc5, 0x15,
	0xc6, 0x65, 0x54, 0xca, 0xe2, 0xe2, 0x85, 0x6a, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x15, 0xd2,
	0x42, 0xb1, 0x41, 0x06, 0xc3, 0x06, 0x88, 0x32, 0xb0, 0x15, 0x42, 0xdc, 0x5c, 0xcc, 0x79, 0xa5,
	0xb9, 0x12, 0x4c, 0x40, 0xa5, 0xcc, 0x4a, 0x0a, 0x50, 0xfb, 0xb8, 0xb9, 0xd8, 0x5d, 0x5c, 0x83,
	0x3d, 0x83, 0x5c, 0x5d, 0x80, 0xb6, 0x01, 0x39, 0xce, 0xa1, 0x41, 0x41, 0x60, 0xbb, 0x8c, 0xbc,
	0xb8, 0x38, 0x82, 0xa1, 0xa6, 0x09, 0xd9, 0x71, 0xb1, 0x82, 0x0d, 0x14, 0x12, 0xc3, 0xee, 0x07,
	0x29, 0x71, 0x1c, 0x36, 0x2b, 0x31, 0x68, 0x30, 0x1a, 0x30, 0x26, 0xb1, 0x81, 0x43, 0xca, 0x18,
	0x10, 0x00, 0x00, 0xff, 0xff, 0x8d, 0xe3, 0x3f, 0x70, 0x3b, 0x01, 0x00, 0x00,
}