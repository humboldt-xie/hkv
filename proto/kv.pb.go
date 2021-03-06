// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kv.proto

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
	SyncServerRequest
	SyncServerReply
	DatasetTop
	ServerInfo
	AddServerRequest
	AddServerReply
	AddDatasetRequest
	AddDatasetReply
	DelDatasetRequest
	DelDatasetReply
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

func (m *Data) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *Data) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Data) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

// The request message containing the user's name.
type SetRequest struct {
	Dataset string `protobuf:"bytes,1,opt,name=Dataset" json:"Dataset,omitempty"`
	Key     []byte `protobuf:"bytes,2,opt,name=Key,proto3" json:"Key,omitempty"`
	Value   []byte `protobuf:"bytes,3,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (m *SetRequest) Reset()                    { *m = SetRequest{} }
func (m *SetRequest) String() string            { return proto.CompactTextString(m) }
func (*SetRequest) ProtoMessage()               {}
func (*SetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SetRequest) GetDataset() string {
	if m != nil {
		return m.Dataset
	}
	return ""
}

func (m *SetRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *SetRequest) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type GetRequest struct {
	Dataset string `protobuf:"bytes,1,opt,name=Dataset" json:"Dataset,omitempty"`
	Key     []byte `protobuf:"bytes,2,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetRequest) GetDataset() string {
	if m != nil {
		return m.Dataset
	}
	return ""
}

func (m *GetRequest) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

// The response message containing the greetings
type SetReply struct {
	Sequence int64 `protobuf:"varint,1,opt,name=Sequence" json:"Sequence,omitempty"`
}

func (m *SetReply) Reset()                    { *m = SetReply{} }
func (m *SetReply) String() string            { return proto.CompactTextString(m) }
func (*SetReply) ProtoMessage()               {}
func (*SetReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SetReply) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

type GetReply struct {
	Data *Data `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
}

func (m *GetReply) Reset()                    { *m = GetReply{} }
func (m *GetReply) String() string            { return proto.CompactTextString(m) }
func (*GetReply) ProtoMessage()               {}
func (*GetReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetReply) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

type Dataset struct {
	Name     string `protobuf:"bytes,1,opt,name=Name" json:"Name,omitempty"`
	Sequence int64  `protobuf:"varint,2,opt,name=Sequence" json:"Sequence,omitempty"`
}

func (m *Dataset) Reset()                    { *m = Dataset{} }
func (m *Dataset) String() string            { return proto.CompactTextString(m) }
func (*Dataset) ProtoMessage()               {}
func (*Dataset) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Dataset) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Dataset) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

type MirrorRequest struct {
	Cmd     string   `protobuf:"bytes,1,opt,name=Cmd" json:"Cmd,omitempty"`
	Addr    string   `protobuf:"bytes,2,opt,name=Addr" json:"Addr,omitempty"`
	Dataset *Dataset `protobuf:"bytes,3,opt,name=Dataset" json:"Dataset,omitempty"`
}

func (m *MirrorRequest) Reset()                    { *m = MirrorRequest{} }
func (m *MirrorRequest) String() string            { return proto.CompactTextString(m) }
func (*MirrorRequest) ProtoMessage()               {}
func (*MirrorRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *MirrorRequest) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *MirrorRequest) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *MirrorRequest) GetDataset() *Dataset {
	if m != nil {
		return m.Dataset
	}
	return nil
}

type MirrorResponse struct {
	Cmd          string   `protobuf:"bytes,1,opt,name=Cmd" json:"Cmd,omitempty"`
	Dataset      string   `protobuf:"bytes,2,opt,name=Dataset" json:"Dataset,omitempty"`
	LastSequence int64    `protobuf:"varint,3,opt,name=LastSequence" json:"LastSequence,omitempty"`
	Data         *Data    `protobuf:"bytes,4,opt,name=Data" json:"Data,omitempty"`
	Args         []string `protobuf:"bytes,5,rep,name=Args" json:"Args,omitempty"`
}

func (m *MirrorResponse) Reset()                    { *m = MirrorResponse{} }
func (m *MirrorResponse) String() string            { return proto.CompactTextString(m) }
func (*MirrorResponse) ProtoMessage()               {}
func (*MirrorResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *MirrorResponse) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *MirrorResponse) GetDataset() string {
	if m != nil {
		return m.Dataset
	}
	return ""
}

func (m *MirrorResponse) GetLastSequence() int64 {
	if m != nil {
		return m.LastSequence
	}
	return 0
}

func (m *MirrorResponse) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *MirrorResponse) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

type SyncServerRequest struct {
	DatasetTopVersion int64            `protobuf:"varint,1,opt,name=DatasetTopVersion" json:"DatasetTopVersion,omitempty"`
	Dataset           map[string]int64 `protobuf:"bytes,2,rep,name=Dataset" json:"Dataset,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
}

func (m *SyncServerRequest) Reset()                    { *m = SyncServerRequest{} }
func (m *SyncServerRequest) String() string            { return proto.CompactTextString(m) }
func (*SyncServerRequest) ProtoMessage()               {}
func (*SyncServerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *SyncServerRequest) GetDatasetTopVersion() int64 {
	if m != nil {
		return m.DatasetTopVersion
	}
	return 0
}

func (m *SyncServerRequest) GetDataset() map[string]int64 {
	if m != nil {
		return m.Dataset
	}
	return nil
}

type SyncServerReply struct {
	Dataset *DatasetTop            `protobuf:"bytes,1,opt,name=Dataset" json:"Dataset,omitempty"`
	Server  map[string]*ServerInfo `protobuf:"bytes,2,rep,name=Server" json:"Server,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *SyncServerReply) Reset()                    { *m = SyncServerReply{} }
func (m *SyncServerReply) String() string            { return proto.CompactTextString(m) }
func (*SyncServerReply) ProtoMessage()               {}
func (*SyncServerReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *SyncServerReply) GetDataset() *DatasetTop {
	if m != nil {
		return m.Dataset
	}
	return nil
}

func (m *SyncServerReply) GetServer() map[string]*ServerInfo {
	if m != nil {
		return m.Server
	}
	return nil
}

// sche server can change
type DatasetTop struct {
	Version int64             `protobuf:"varint,1,opt,name=Version" json:"Version,omitempty"`
	Table   map[string]string `protobuf:"bytes,2,rep,name=Table" json:"Table,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *DatasetTop) Reset()                    { *m = DatasetTop{} }
func (m *DatasetTop) String() string            { return proto.CompactTextString(m) }
func (*DatasetTop) ProtoMessage()               {}
func (*DatasetTop) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *DatasetTop) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *DatasetTop) GetTable() map[string]string {
	if m != nil {
		return m.Table
	}
	return nil
}

type ServerInfo struct {
	Id      string            `protobuf:"bytes,1,opt,name=Id" json:"Id,omitempty"`
	Version int64             `protobuf:"varint,2,opt,name=Version" json:"Version,omitempty"`
	Addr    string            `protobuf:"bytes,3,opt,name=Addr" json:"Addr,omitempty"`
	Dataset map[string]string `protobuf:"bytes,4,rep,name=Dataset" json:"Dataset,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *ServerInfo) Reset()                    { *m = ServerInfo{} }
func (m *ServerInfo) String() string            { return proto.CompactTextString(m) }
func (*ServerInfo) ProtoMessage()               {}
func (*ServerInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *ServerInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ServerInfo) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ServerInfo) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *ServerInfo) GetDataset() map[string]string {
	if m != nil {
		return m.Dataset
	}
	return nil
}

type AddServerRequest struct {
}

func (m *AddServerRequest) Reset()                    { *m = AddServerRequest{} }
func (m *AddServerRequest) String() string            { return proto.CompactTextString(m) }
func (*AddServerRequest) ProtoMessage()               {}
func (*AddServerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type AddServerReply struct {
}

func (m *AddServerReply) Reset()                    { *m = AddServerReply{} }
func (m *AddServerReply) String() string            { return proto.CompactTextString(m) }
func (*AddServerReply) ProtoMessage()               {}
func (*AddServerReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

type AddDatasetRequest struct {
	ServerId string `protobuf:"bytes,1,opt,name=ServerId" json:"ServerId,omitempty"`
	Dataset  string `protobuf:"bytes,2,opt,name=Dataset" json:"Dataset,omitempty"`
}

func (m *AddDatasetRequest) Reset()                    { *m = AddDatasetRequest{} }
func (m *AddDatasetRequest) String() string            { return proto.CompactTextString(m) }
func (*AddDatasetRequest) ProtoMessage()               {}
func (*AddDatasetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *AddDatasetRequest) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *AddDatasetRequest) GetDataset() string {
	if m != nil {
		return m.Dataset
	}
	return ""
}

type AddDatasetReply struct {
}

func (m *AddDatasetReply) Reset()                    { *m = AddDatasetReply{} }
func (m *AddDatasetReply) String() string            { return proto.CompactTextString(m) }
func (*AddDatasetReply) ProtoMessage()               {}
func (*AddDatasetReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type DelDatasetRequest struct {
	ServerId string `protobuf:"bytes,1,opt,name=ServerId" json:"ServerId,omitempty"`
	Dataset  string `protobuf:"bytes,2,opt,name=Dataset" json:"Dataset,omitempty"`
}

func (m *DelDatasetRequest) Reset()                    { *m = DelDatasetRequest{} }
func (m *DelDatasetRequest) String() string            { return proto.CompactTextString(m) }
func (*DelDatasetRequest) ProtoMessage()               {}
func (*DelDatasetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *DelDatasetRequest) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
}

func (m *DelDatasetRequest) GetDataset() string {
	if m != nil {
		return m.Dataset
	}
	return ""
}

type DelDatasetReply struct {
	ServerId string `protobuf:"bytes,1,opt,name=ServerId" json:"ServerId,omitempty"`
}

func (m *DelDatasetReply) Reset()                    { *m = DelDatasetReply{} }
func (m *DelDatasetReply) String() string            { return proto.CompactTextString(m) }
func (*DelDatasetReply) ProtoMessage()               {}
func (*DelDatasetReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *DelDatasetReply) GetServerId() string {
	if m != nil {
		return m.ServerId
	}
	return ""
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
	proto.RegisterType((*SyncServerRequest)(nil), "kv.SyncServerRequest")
	proto.RegisterType((*SyncServerReply)(nil), "kv.SyncServerReply")
	proto.RegisterType((*DatasetTop)(nil), "kv.DatasetTop")
	proto.RegisterType((*ServerInfo)(nil), "kv.ServerInfo")
	proto.RegisterType((*AddServerRequest)(nil), "kv.AddServerRequest")
	proto.RegisterType((*AddServerReply)(nil), "kv.AddServerReply")
	proto.RegisterType((*AddDatasetRequest)(nil), "kv.AddDatasetRequest")
	proto.RegisterType((*AddDatasetReply)(nil), "kv.AddDatasetReply")
	proto.RegisterType((*DelDatasetRequest)(nil), "kv.DelDatasetRequest")
	proto.RegisterType((*DelDatasetReply)(nil), "kv.DelDatasetReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

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
	Metadata: "kv.proto",
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
	Metadata: "kv.proto",
}

// Client API for Cluster service

type ClusterClient interface {
	AddServer(ctx context.Context, in *AddServerRequest, opts ...grpc.CallOption) (*AddServerReply, error)
	AddDataset(ctx context.Context, in *AddDatasetRequest, opts ...grpc.CallOption) (*AddDatasetReply, error)
	DelDataset(ctx context.Context, in *DelDatasetRequest, opts ...grpc.CallOption) (*DelDatasetReply, error)
	SyncServer(ctx context.Context, in *SyncServerRequest, opts ...grpc.CallOption) (*SyncServerReply, error)
}

type clusterClient struct {
	cc *grpc.ClientConn
}

func NewClusterClient(cc *grpc.ClientConn) ClusterClient {
	return &clusterClient{cc}
}

func (c *clusterClient) AddServer(ctx context.Context, in *AddServerRequest, opts ...grpc.CallOption) (*AddServerReply, error) {
	out := new(AddServerReply)
	err := grpc.Invoke(ctx, "/kv.Cluster/AddServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) AddDataset(ctx context.Context, in *AddDatasetRequest, opts ...grpc.CallOption) (*AddDatasetReply, error) {
	out := new(AddDatasetReply)
	err := grpc.Invoke(ctx, "/kv.Cluster/AddDataset", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) DelDataset(ctx context.Context, in *DelDatasetRequest, opts ...grpc.CallOption) (*DelDatasetReply, error) {
	out := new(DelDatasetReply)
	err := grpc.Invoke(ctx, "/kv.Cluster/DelDataset", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clusterClient) SyncServer(ctx context.Context, in *SyncServerRequest, opts ...grpc.CallOption) (*SyncServerReply, error) {
	out := new(SyncServerReply)
	err := grpc.Invoke(ctx, "/kv.Cluster/SyncServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Cluster service

type ClusterServer interface {
	AddServer(context.Context, *AddServerRequest) (*AddServerReply, error)
	AddDataset(context.Context, *AddDatasetRequest) (*AddDatasetReply, error)
	DelDataset(context.Context, *DelDatasetRequest) (*DelDatasetReply, error)
	SyncServer(context.Context, *SyncServerRequest) (*SyncServerReply, error)
}

func RegisterClusterServer(s *grpc.Server, srv ClusterServer) {
	s.RegisterService(&_Cluster_serviceDesc, srv)
}

func _Cluster_AddServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).AddServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv.Cluster/AddServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).AddServer(ctx, req.(*AddServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_AddDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).AddDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv.Cluster/AddDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).AddDataset(ctx, req.(*AddDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_DelDataset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelDatasetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).DelDataset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv.Cluster/DelDataset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).DelDataset(ctx, req.(*DelDatasetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cluster_SyncServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterServer).SyncServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv.Cluster/SyncServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterServer).SyncServer(ctx, req.(*SyncServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Cluster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kv.Cluster",
	HandlerType: (*ClusterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddServer",
			Handler:    _Cluster_AddServer_Handler,
		},
		{
			MethodName: "AddDataset",
			Handler:    _Cluster_AddDataset_Handler,
		},
		{
			MethodName: "DelDataset",
			Handler:    _Cluster_DelDataset_Handler,
		},
		{
			MethodName: "SyncServer",
			Handler:    _Cluster_SyncServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kv.proto",
}

func init() { proto.RegisterFile("kv.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 705 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x8e, 0xed, 0xb4, 0x4d, 0x26, 0x21, 0x4d, 0x96, 0x22, 0x05, 0x83, 0x44, 0xb5, 0xa2, 0x28,
	0x07, 0x08, 0x28, 0xa8, 0x6a, 0xa9, 0x90, 0x50, 0xd5, 0xa2, 0x28, 0x14, 0x7a, 0x70, 0xaa, 0x9e,
	0xb8, 0xb8, 0xf5, 0x52, 0x55, 0x71, 0x63, 0x63, 0x6f, 0x23, 0xf9, 0x31, 0xe0, 0x75, 0x10, 0xe2,
	0xad, 0x38, 0xa3, 0xfd, 0xf5, 0xae, 0x53, 0x7e, 0x84, 0x7a, 0xdb, 0x1d, 0xcf, 0x7c, 0xf3, 0xcd,
	0x37, 0xb3, 0x63, 0x68, 0xcc, 0x16, 0xc3, 0x34, 0x4b, 0x68, 0x82, 0xdc, 0xd9, 0x02, 0xbf, 0x83,
	0xfa, 0x61, 0x48, 0x43, 0xe4, 0x43, 0x63, 0x4a, 0x3e, 0x5f, 0x93, 0xf9, 0x39, 0xe9, 0x3b, 0x9b,
	0xce, 0xc0, 0x0b, 0xf4, 0x1d, 0x75, 0xc1, 0x3b, 0x22, 0x45, 0xdf, 0xdd, 0x74, 0x06, 0xed, 0x80,
	0x1d, 0xd1, 0x06, 0xac, 0x9c, 0x86, 0xf1, 0x35, 0xe9, 0x7b, 0xdc, 0x26, 0x2e, 0xf8, 0x18, 0x60,
	0x4a, 0x68, 0xc0, 0xc2, 0x72, 0x8a, 0xfa, 0xb0, 0xc6, 0x90, 0x73, 0x42, 0x39, 0x60, 0x33, 0x50,
	0xd7, 0x7f, 0xc6, 0xdb, 0x05, 0x18, 0xff, 0x17, 0x1e, 0x7e, 0xc2, 0xaa, 0xa1, 0x01, 0x49, 0xe3,
	0xe2, 0x4f, 0x95, 0xe1, 0x01, 0x34, 0xc6, 0xca, 0xef, 0x21, 0xd4, 0xa3, 0x90, 0x86, 0xdc, 0xa7,
	0x35, 0x6a, 0x0c, 0x67, 0x8b, 0x21, 0x4b, 0x10, 0x70, 0x2b, 0x7e, 0xa5, 0xb3, 0x23, 0x04, 0xf5,
	0xe3, 0xf0, 0x8a, 0x48, 0x16, 0xfc, 0x6c, 0x25, 0x71, 0x2b, 0x49, 0x3e, 0xc2, 0x9d, 0x0f, 0x97,
	0x59, 0x96, 0x64, 0xaa, 0x92, 0x2e, 0x78, 0x07, 0x57, 0x91, 0x8c, 0x67, 0x47, 0x06, 0xb9, 0x1f,
	0x45, 0x19, 0x0f, 0x6d, 0x06, 0xfc, 0x8c, 0xb6, 0xca, 0x7a, 0x3d, 0x4e, 0xa9, 0xa5, 0x28, 0xe5,
	0x84, 0xea, 0xe2, 0xf1, 0x57, 0x07, 0x3a, 0x0a, 0x3e, 0x4f, 0x93, 0x79, 0x4e, 0x6e, 0xc0, 0x37,
	0xb4, 0x73, 0x6d, 0xed, 0x30, 0xb4, 0xdf, 0x87, 0x39, 0xd5, 0xe4, 0x3d, 0x4e, 0xde, 0xb2, 0x31,
	0x65, 0x98, 0x7b, 0xbf, 0x5e, 0x55, 0x86, 0x4f, 0x0e, 0xe3, 0x9e, 0x5d, 0xe4, 0xfd, 0x95, 0x4d,
	0x8f, 0x73, 0xcf, 0x2e, 0x72, 0xfc, 0xcd, 0x81, 0xde, 0xb4, 0x98, 0x9f, 0x4f, 0x49, 0xb6, 0x20,
	0xba, 0xee, 0xa7, 0xd0, 0x93, 0x69, 0x4f, 0x92, 0xf4, 0x94, 0x64, 0xf9, 0x65, 0x32, 0x97, 0x2d,
	0x59, 0xfe, 0x80, 0x5e, 0x9b, 0x9c, 0xbd, 0x41, 0x6b, 0x84, 0x59, 0xe2, 0x25, 0x54, 0xa5, 0xc8,
	0xdb, 0x39, 0xcd, 0x0a, 0x5d, 0x97, 0xbf, 0x07, 0x6d, 0xf3, 0x03, 0xd3, 0x64, 0x46, 0x0a, 0xa5,
	0xc9, 0x4c, 0xcc, 0xdc, 0x82, 0xcf, 0x9c, 0xe8, 0x97, 0xb8, 0xec, 0xb9, 0xbb, 0x0e, 0xfe, 0xe1,
	0xc0, 0xba, 0x99, 0x87, 0x4d, 0xc7, 0xc0, 0x9e, 0xbe, 0xd6, 0xa8, 0x63, 0x74, 0xe3, 0x24, 0x49,
	0x4b, 0x45, 0x77, 0x60, 0x55, 0x04, 0x4a, 0xda, 0x8f, 0xaa, 0xb4, 0xd3, 0xb8, 0x18, 0x8a, 0xb3,
	0xe0, 0x2c, 0xdd, 0xfd, 0x09, 0xb4, 0x0c, 0xf3, 0x0d, 0x8c, 0x1f, 0x9b, 0x8c, 0x25, 0x03, 0x11,
	0x31, 0x99, 0x7f, 0x4a, 0xcc, 0x0a, 0xbe, 0x38, 0x00, 0x25, 0x37, 0xd6, 0x7e, 0x5b, 0x6e, 0x75,
	0x45, 0xcf, 0x61, 0xe5, 0x24, 0x3c, 0x8b, 0x89, 0xe4, 0x7a, 0xdf, 0x2e, 0x6a, 0xc8, 0xbf, 0x09,
	0x96, 0xc2, 0xcf, 0xdf, 0x05, 0x28, 0x8d, 0x7f, 0x53, 0xb5, 0x69, 0x72, 0xfa, 0xee, 0xb0, 0xf5,
	0xa0, 0xd8, 0xa2, 0x0e, 0xb8, 0x13, 0x35, 0xa3, 0xee, 0x24, 0x32, 0x39, 0xba, 0x36, 0x47, 0xf5,
	0x38, 0x3c, 0xe3, 0x71, 0x6c, 0x97, 0xed, 0xa8, 0x73, 0xe6, 0x0f, 0x6c, 0x31, 0x6e, 0x67, 0x2a,
	0x2c, 0xfe, 0x08, 0xba, 0xfb, 0x51, 0x64, 0xcd, 0x1e, 0xee, 0x42, 0xc7, 0xb0, 0xa5, 0x71, 0x81,
	0x27, 0xd0, 0xdb, 0x8f, 0x22, 0xf5, 0x4a, 0xe5, 0xe0, 0xf3, 0xed, 0xc0, 0xa9, 0xa9, 0x8a, 0xf5,
	0xfd, 0xf7, 0x4f, 0x13, 0xf7, 0x60, 0xdd, 0x84, 0x92, 0xe8, 0x87, 0x24, 0xbe, 0x15, 0xf4, 0x67,
	0xb0, 0x6e, 0x42, 0xe9, 0x4d, 0x79, 0x33, 0xd0, 0x28, 0x00, 0xf7, 0x68, 0x81, 0xb6, 0xc0, 0x9b,
	0x12, 0x8a, 0xe4, 0xe4, 0x29, 0x06, 0x7e, 0x5b, 0xdf, 0x19, 0xc9, 0x1a, 0x73, 0x1b, 0x2b, 0xb7,
	0x71, 0xc5, 0x6d, 0xac, 0xdd, 0x46, 0x6f, 0x60, 0x55, 0x6c, 0x2e, 0xb4, 0xad, 0x4f, 0x3d, 0xe6,
	0x63, 0xad, 0x4b, 0x1f, 0x99, 0x26, 0xb1, 0xe2, 0x70, 0x6d, 0xe0, 0xbc, 0x70, 0x46, 0x3f, 0x1d,
	0x58, 0x3b, 0x88, 0xaf, 0x73, 0x4a, 0x32, 0xb4, 0x03, 0x4d, 0xdd, 0x0a, 0xb4, 0xc1, 0x42, 0xaa,
	0xdd, 0x12, 0x40, 0x95, 0x7e, 0xd5, 0xd0, 0x1e, 0x40, 0x29, 0x33, 0xba, 0x27, 0x7d, 0x6c, 0x8d,
	0xfd, 0xbb, 0x55, 0xb3, 0x8e, 0x2d, 0x45, 0x14, 0xb1, 0x4b, 0xfd, 0x11, 0xb1, 0x15, 0xad, 0x45,
	0x6c, 0xb9, 0x15, 0x44, 0xec, 0xd2, 0x72, 0x13, 0xb1, 0x95, 0xe5, 0x81, 0x6b, 0x67, 0xab, 0xfc,
	0x07, 0xfe, 0xf2, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf7, 0x65, 0x31, 0xf4, 0xcc, 0x07, 0x00,
	0x00,
}
