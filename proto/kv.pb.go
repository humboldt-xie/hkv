// Code generated by protoc-gen-go.
// source: kv.proto
// DO NOT EDIT!

/*
Package kv is a generated protocol buffer package.

It is generated from these files:
	kv.proto

It has these top-level messages:
	Data
	SetRequest
	GetRequest
	SetReply
	GetReply
	Dataset
	MirrorRequest
	MirrorResponse
*/
package kv

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
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Data struct {
	Sequence int64  `protobuf:"varint,1,opt,name=Sequence" json:"Sequence,omitempty"`
	Key      []byte `protobuf:"bytes,2,opt,name=Key,proto3" json:"Key,omitempty"`
	Value    []byte `protobuf:"bytes,3,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (m *Data) Reset()                    { *m = Data{} }
func (m *Data) String() string            { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()               {}
func (*Data) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The request message containing the user's name.
type SetRequest struct {
	Key   string `protobuf:"bytes,1,opt,name=Key" json:"Key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=Value" json:"Value,omitempty"`
}

func (m *SetRequest) Reset()                    { *m = SetRequest{} }
func (m *SetRequest) String() string            { return proto.CompactTextString(m) }
func (*SetRequest) ProtoMessage()               {}
func (*SetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type GetRequest struct {
	Key string `protobuf:"bytes,1,opt,name=Key" json:"Key,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// The response message containing the greetings
type SetReply struct {
	Result string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
}

func (m *SetReply) Reset()                    { *m = SetReply{} }
func (m *SetReply) String() string            { return proto.CompactTextString(m) }
func (*SetReply) ProtoMessage()               {}
func (*SetReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type GetReply struct {
	Value string `protobuf:"bytes,1,opt,name=Value" json:"Value,omitempty"`
}

func (m *GetReply) Reset()                    { *m = GetReply{} }
func (m *GetReply) String() string            { return proto.CompactTextString(m) }
func (*GetReply) ProtoMessage()               {}
func (*GetReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type Dataset struct {
	Name     string `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Sequence int64  `protobuf:"varint,2,opt,name=Sequence" json:"Sequence,omitempty"`
}

func (m *Dataset) Reset()                    { *m = Dataset{} }
func (m *Dataset) String() string            { return proto.CompactTextString(m) }
func (*Dataset) ProtoMessage()               {}
func (*Dataset) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type MirrorRequest struct {
	Cmd     string   `protobuf:"bytes,1,opt,name=Cmd" json:"Cmd,omitempty"`
	Addr    string   `protobuf:"bytes,2,opt,name=Addr" json:"Addr,omitempty"`
	Dataset *Dataset `protobuf:"bytes,3,opt,name=Dataset" json:"Dataset,omitempty"`
}

func (m *MirrorRequest) Reset()                    { *m = MirrorRequest{} }
func (m *MirrorRequest) String() string            { return proto.CompactTextString(m) }
func (*MirrorRequest) ProtoMessage()               {}
func (*MirrorRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *MirrorRequest) GetDataset() *Dataset {
	if m != nil {
		return m.Dataset
	}
	return nil
}

type MirrorResponse struct {
	Cmd     string   `protobuf:"bytes,1,opt,name=Cmd" json:"Cmd,omitempty"`
	Dataset string   `protobuf:"bytes,2,opt,name=Dataset" json:"Dataset,omitempty"`
	Args    []string `protobuf:"bytes,3,rep,name=Args" json:"Args,omitempty"`
	Data    *Data    `protobuf:"bytes,4,opt,name=Data" json:"Data,omitempty"`
}

func (m *MirrorResponse) Reset()                    { *m = MirrorResponse{} }
func (m *MirrorResponse) String() string            { return proto.CompactTextString(m) }
func (*MirrorResponse) ProtoMessage()               {}
func (*MirrorResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *MirrorResponse) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Data)(nil), "kv.Data")
	proto.RegisterType((*SetRequest)(nil), "kv.SetRequest")
	proto.RegisterType((*GetRequest)(nil), "kv.GetRequest")
	proto.RegisterType((*SetReply)(nil), "kv.SetReply")
	proto.RegisterType((*GetReply)(nil), "kv.GetReply")
	proto.RegisterType((*Dataset)(nil), "kv.Dataset")
	proto.RegisterType((*MirrorRequest)(nil), "kv.MirrorRequest")
	proto.RegisterType((*MirrorResponse)(nil), "kv.MirrorResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Kv service

type KvClient interface {
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetReply, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
}

type kvClient struct {
	cc *grpc.ClientConn
}

func NewKvClient(cc *grpc.ClientConn) KvClient {
	return &kvClient{cc}
}

func (c *kvClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetReply, error) {
	out := new(SetReply)
	err := grpc.Invoke(ctx, "/kv.Kv/Set", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kvClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := grpc.Invoke(ctx, "/kv.Kv/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Kv service

type KvServer interface {
	Set(context.Context, *SetRequest) (*SetReply, error)
	Get(context.Context, *GetRequest) (*GetReply, error)
}

func RegisterKvServer(s *grpc.Server, srv KvServer) {
	s.RegisterService(&_Kv_serviceDesc, srv)
}

func _Kv_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv.Kv/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kv_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv.Kv/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Kv_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kv.Kv",
	HandlerType: (*KvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _Kv_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Kv_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

// Client API for Mirror service

type MirrorClient interface {
	Mirror(ctx context.Context, opts ...grpc.CallOption) (Mirror_MirrorClient, error)
}

type mirrorClient struct {
	cc *grpc.ClientConn
}

func NewMirrorClient(cc *grpc.ClientConn) MirrorClient {
	return &mirrorClient{cc}
}

func (c *mirrorClient) Mirror(ctx context.Context, opts ...grpc.CallOption) (Mirror_MirrorClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Mirror_serviceDesc.Streams[0], c.cc, "/kv.Mirror/Mirror", opts...)
	if err != nil {
		return nil, err
	}
	x := &mirrorMirrorClient{stream}
	return x, nil
}

type Mirror_MirrorClient interface {
	Send(*MirrorRequest) error
	Recv() (*MirrorResponse, error)
	grpc.ClientStream
}

type mirrorMirrorClient struct {
	grpc.ClientStream
}

func (x *mirrorMirrorClient) Send(m *MirrorRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mirrorMirrorClient) Recv() (*MirrorResponse, error) {
	m := new(MirrorResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Mirror service

type MirrorServer interface {
	Mirror(Mirror_MirrorServer) error
}

func RegisterMirrorServer(s *grpc.Server, srv MirrorServer) {
	s.RegisterService(&_Mirror_serviceDesc, srv)
}

func _Mirror_Mirror_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MirrorServer).Mirror(&mirrorMirrorServer{stream})
}

type Mirror_MirrorServer interface {
	Send(*MirrorResponse) error
	Recv() (*MirrorRequest, error)
	grpc.ServerStream
}

type mirrorMirrorServer struct {
	grpc.ServerStream
}

func (x *mirrorMirrorServer) Send(m *MirrorResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mirrorMirrorServer) Recv() (*MirrorRequest, error) {
	m := new(MirrorRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Mirror_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kv.Mirror",
	HandlerType: (*MirrorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Mirror",
			Handler:       _Mirror_Mirror_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("kv.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 338 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x52, 0xdf, 0x4b, 0xc3, 0x30,
	0x10, 0x5e, 0xdb, 0x39, 0xeb, 0x6d, 0x0e, 0x3d, 0x44, 0xca, 0x10, 0x19, 0x01, 0x61, 0x4f, 0x43,
	0xa6, 0x3e, 0xf8, 0x24, 0x43, 0xa1, 0xe0, 0xd0, 0x87, 0x0c, 0x7c, 0xf2, 0xa5, 0xba, 0x20, 0xb2,
	0x1f, 0xad, 0x49, 0x5a, 0xd8, 0x7f, 0x6f, 0x92, 0x36, 0x59, 0x3b, 0xc4, 0xb7, 0xbb, 0xef, 0xbe,
	0xbb, 0xef, 0xbe, 0x4b, 0x20, 0x5c, 0x16, 0xe3, 0x8c, 0xa7, 0x32, 0x45, 0x7f, 0x59, 0x90, 0x67,
	0x68, 0x3f, 0x25, 0x32, 0xc1, 0x01, 0x84, 0x73, 0xf6, 0x93, 0xb3, 0xcd, 0x27, 0x8b, 0xbc, 0xa1,
	0x37, 0x0a, 0xa8, 0xcb, 0xf1, 0x04, 0x82, 0x19, 0xdb, 0x46, 0xbe, 0x82, 0x7b, 0x54, 0x87, 0x78,
	0x06, 0x07, 0x6f, 0xc9, 0x2a, 0x67, 0x51, 0x60, 0xb0, 0x32, 0x21, 0xb7, 0x00, 0x73, 0x26, 0xa9,
	0x6e, 0x13, 0xd2, 0x76, 0xe9, 0x61, 0x47, 0x7b, 0x5d, 0xbe, 0xc1, 0xaa, 0xae, 0x4b, 0x80, 0xf8,
	0x9f, 0x2e, 0x42, 0xf4, 0x66, 0xaa, 0x9e, 0xad, 0xb6, 0x78, 0x0e, 0x1d, 0xce, 0x44, 0xbe, 0x92,
	0x15, 0xa1, 0xca, 0xc8, 0x10, 0xc2, 0xd8, 0x72, 0x9c, 0x8a, 0x57, 0x57, 0xb9, 0x87, 0x43, 0xed,
	0x53, 0x30, 0x89, 0x08, 0xed, 0xd7, 0x64, 0x6d, 0xeb, 0x26, 0x6e, 0xd8, 0xf7, 0x9b, 0xf6, 0xc9,
	0x3b, 0x1c, 0xbf, 0x7c, 0x73, 0x9e, 0xf2, 0xda, 0x8e, 0x8f, 0xeb, 0x85, 0xdd, 0x51, 0x85, 0x7a,
	0xe4, 0x74, 0xb1, 0xe0, 0x95, 0x31, 0x13, 0xe3, 0x95, 0x53, 0x34, 0x57, 0xea, 0x4e, 0xba, 0x63,
	0x75, 0xf9, 0x0a, 0xa2, 0xb6, 0x46, 0x36, 0xd0, 0xb7, 0xd3, 0x45, 0x96, 0x6e, 0x04, 0xfb, 0x63,
	0x7c, 0xb4, 0x1b, 0x55, 0x2a, 0xd4, 0xbd, 0x4c, 0xf9, 0x97, 0x50, 0x0a, 0x81, 0x11, 0x56, 0x31,
	0x5e, 0x94, 0x4f, 0x1a, 0xb5, 0x8d, 0x6a, 0x68, 0x55, 0xa9, 0x41, 0x27, 0x14, 0xfc, 0x59, 0xa1,
	0x96, 0x0b, 0xd4, 0x51, 0xb1, 0xaf, 0x8b, 0xbb, 0x37, 0x1b, 0xf4, 0x5c, 0xae, 0x2e, 0x49, 0x5a,
	0x9a, 0x16, 0x5b, 0x5a, 0xbc, 0x47, 0x8b, 0x1d, 0x6d, 0xf2, 0x00, 0x9d, 0xd2, 0x03, 0xde, 0xb9,
	0xe8, 0x54, 0x73, 0x1a, 0x77, 0x1b, 0x60, 0x1d, 0x2a, 0xcd, 0x92, 0xd6, 0xc8, 0xbb, 0xf6, 0x3e,
	0x3a, 0xe6, 0x43, 0xde, 0xfc, 0x06, 0x00, 0x00, 0xff, 0xff, 0x0e, 0xe9, 0x32, 0x8e, 0x9c, 0x02,
	0x00, 0x00,
}