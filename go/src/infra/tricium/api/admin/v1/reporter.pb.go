// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/tricium/api/admin/v1/reporter.proto

package admin

import prpc "go.chromium.org/luci/grpc/prpc"

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

type ReportResultsRequest struct {
	RunId                int64    `protobuf:"varint,1,opt,name=run_id,json=runId" json:"run_id,omitempty"`
	Analyzer             string   `protobuf:"bytes,2,opt,name=analyzer" json:"analyzer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportResultsRequest) Reset()         { *m = ReportResultsRequest{} }
func (m *ReportResultsRequest) String() string { return proto.CompactTextString(m) }
func (*ReportResultsRequest) ProtoMessage()    {}
func (*ReportResultsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_reporter_e472c2b31bcd314c, []int{0}
}
func (m *ReportResultsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportResultsRequest.Unmarshal(m, b)
}
func (m *ReportResultsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportResultsRequest.Marshal(b, m, deterministic)
}
func (dst *ReportResultsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportResultsRequest.Merge(dst, src)
}
func (m *ReportResultsRequest) XXX_Size() int {
	return xxx_messageInfo_ReportResultsRequest.Size(m)
}
func (m *ReportResultsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportResultsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReportResultsRequest proto.InternalMessageInfo

func (m *ReportResultsRequest) GetRunId() int64 {
	if m != nil {
		return m.RunId
	}
	return 0
}

func (m *ReportResultsRequest) GetAnalyzer() string {
	if m != nil {
		return m.Analyzer
	}
	return ""
}

type ReportResultsResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportResultsResponse) Reset()         { *m = ReportResultsResponse{} }
func (m *ReportResultsResponse) String() string { return proto.CompactTextString(m) }
func (*ReportResultsResponse) ProtoMessage()    {}
func (*ReportResultsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_reporter_e472c2b31bcd314c, []int{1}
}
func (m *ReportResultsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportResultsResponse.Unmarshal(m, b)
}
func (m *ReportResultsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportResultsResponse.Marshal(b, m, deterministic)
}
func (dst *ReportResultsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportResultsResponse.Merge(dst, src)
}
func (m *ReportResultsResponse) XXX_Size() int {
	return xxx_messageInfo_ReportResultsResponse.Size(m)
}
func (m *ReportResultsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportResultsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ReportResultsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ReportResultsRequest)(nil), "admin.ReportResultsRequest")
	proto.RegisterType((*ReportResultsResponse)(nil), "admin.ReportResultsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ReporterClient is the client API for Reporter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReporterClient interface {
	// ReportResults reports Tricium results.
	ReportResults(ctx context.Context, in *ReportResultsRequest, opts ...grpc.CallOption) (*ReportResultsResponse, error)
}
type reporterPRPCClient struct {
	client *prpc.Client
}

func NewReporterPRPCClient(client *prpc.Client) ReporterClient {
	return &reporterPRPCClient{client}
}

func (c *reporterPRPCClient) ReportResults(ctx context.Context, in *ReportResultsRequest, opts ...grpc.CallOption) (*ReportResultsResponse, error) {
	out := new(ReportResultsResponse)
	err := c.client.Call(ctx, "admin.Reporter", "ReportResults", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type reporterClient struct {
	cc *grpc.ClientConn
}

func NewReporterClient(cc *grpc.ClientConn) ReporterClient {
	return &reporterClient{cc}
}

func (c *reporterClient) ReportResults(ctx context.Context, in *ReportResultsRequest, opts ...grpc.CallOption) (*ReportResultsResponse, error) {
	out := new(ReportResultsResponse)
	err := c.cc.Invoke(ctx, "/admin.Reporter/ReportResults", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReporterServer is the server API for Reporter service.
type ReporterServer interface {
	// ReportResults reports Tricium results.
	ReportResults(context.Context, *ReportResultsRequest) (*ReportResultsResponse, error)
}

func RegisterReporterServer(s prpc.Registrar, srv ReporterServer) {
	s.RegisterService(&_Reporter_serviceDesc, srv)
}

func _Reporter_ReportResults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReportResultsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReporterServer).ReportResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin.Reporter/ReportResults",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReporterServer).ReportResults(ctx, req.(*ReportResultsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Reporter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "admin.Reporter",
	HandlerType: (*ReporterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReportResults",
			Handler:    _Reporter_ReportResults_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infra/tricium/api/admin/v1/reporter.proto",
}

func init() {
	proto.RegisterFile("infra/tricium/api/admin/v1/reporter.proto", fileDescriptor_reporter_e472c2b31bcd314c)
}

var fileDescriptor_reporter_e472c2b31bcd314c = []byte{
	// 180 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0xb1, 0x0a, 0xc2, 0x30,
	0x10, 0x40, 0xa9, 0xd2, 0x52, 0x03, 0x2e, 0xc1, 0x62, 0xa9, 0x0e, 0xa5, 0x53, 0x5d, 0x1a, 0xd4,
	0xaf, 0xa8, 0x63, 0x06, 0x57, 0x89, 0xf6, 0x84, 0x40, 0x9b, 0xc4, 0x4b, 0x22, 0xe8, 0xd7, 0x0b,
	0x51, 0x04, 0xa5, 0xe3, 0xf1, 0x8e, 0x77, 0xef, 0xc8, 0x46, 0xaa, 0x2b, 0x0a, 0xe6, 0x50, 0x5e,
	0xa4, 0x1f, 0x98, 0x30, 0x92, 0x89, 0x6e, 0x90, 0x8a, 0xdd, 0xb7, 0x0c, 0xc1, 0x68, 0x74, 0x80,
	0x8d, 0x41, 0xed, 0x34, 0x8d, 0x03, 0xa8, 0x5a, 0xb2, 0xe0, 0x01, 0x70, 0xb0, 0xbe, 0x77, 0x96,
	0xc3, 0xcd, 0x83, 0x75, 0x34, 0x23, 0x09, 0x7a, 0x75, 0x92, 0x5d, 0x1e, 0x95, 0x51, 0x3d, 0xe5,
	0x31, 0x7a, 0xd5, 0x76, 0xb4, 0x20, 0xa9, 0x50, 0xa2, 0x7f, 0x3c, 0x01, 0xf3, 0x49, 0x19, 0xd5,
	0x33, 0xfe, 0x9d, 0xab, 0x25, 0xc9, 0xfe, 0x54, 0xd6, 0x68, 0x65, 0x61, 0x77, 0x24, 0x29, 0xff,
	0x1c, 0xa7, 0x07, 0x32, 0xff, 0x59, 0xa2, 0xab, 0x26, 0x84, 0x34, 0x63, 0x15, 0xc5, 0x7a, 0x1c,
	0xbe, 0xbd, 0xe7, 0x24, 0x7c, 0xb2, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0x36, 0x37, 0xc9, 0x5d,
	0xf6, 0x00, 0x00, 0x00,
}
