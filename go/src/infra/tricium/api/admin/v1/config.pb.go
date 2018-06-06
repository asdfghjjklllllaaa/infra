// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/tricium/api/admin/v1/config.proto

package admin

import prpc "go.chromium.org/luci/grpc/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v1 "infra/tricium/api/v1"

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

type ValidateRequest struct {
	// The project configuration to validate.
	ProjectConfig *v1.ProjectConfig `protobuf:"bytes,1,opt,name=project_config,json=projectConfig" json:"project_config,omitempty"`
	// The service config to use (optional).
	//
	// If not provided, the default service config will be used.
	ServiceConfig        *v1.ServiceConfig `protobuf:"bytes,2,opt,name=service_config,json=serviceConfig" json:"service_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ValidateRequest) Reset()         { *m = ValidateRequest{} }
func (m *ValidateRequest) String() string { return proto.CompactTextString(m) }
func (*ValidateRequest) ProtoMessage()    {}
func (*ValidateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_043929b00c61711a, []int{0}
}
func (m *ValidateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidateRequest.Unmarshal(m, b)
}
func (m *ValidateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidateRequest.Marshal(b, m, deterministic)
}
func (dst *ValidateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidateRequest.Merge(dst, src)
}
func (m *ValidateRequest) XXX_Size() int {
	return xxx_messageInfo_ValidateRequest.Size(m)
}
func (m *ValidateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ValidateRequest proto.InternalMessageInfo

func (m *ValidateRequest) GetProjectConfig() *v1.ProjectConfig {
	if m != nil {
		return m.ProjectConfig
	}
	return nil
}

func (m *ValidateRequest) GetServiceConfig() *v1.ServiceConfig {
	if m != nil {
		return m.ServiceConfig
	}
	return nil
}

// TODO(emso): Return structured errors for invalid configs.
type ValidateResponse struct {
	// The config used for validation.
	//
	// This is the resulting config after flattening and merging the provided
	// project and service config.
	ValidatedConfig      *v1.ProjectConfig `protobuf:"bytes,1,opt,name=validated_config,json=validatedConfig" json:"validated_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ValidateResponse) Reset()         { *m = ValidateResponse{} }
func (m *ValidateResponse) String() string { return proto.CompactTextString(m) }
func (*ValidateResponse) ProtoMessage()    {}
func (*ValidateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_043929b00c61711a, []int{1}
}
func (m *ValidateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidateResponse.Unmarshal(m, b)
}
func (m *ValidateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidateResponse.Marshal(b, m, deterministic)
}
func (dst *ValidateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidateResponse.Merge(dst, src)
}
func (m *ValidateResponse) XXX_Size() int {
	return xxx_messageInfo_ValidateResponse.Size(m)
}
func (m *ValidateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ValidateResponse proto.InternalMessageInfo

func (m *ValidateResponse) GetValidatedConfig() *v1.ProjectConfig {
	if m != nil {
		return m.ValidatedConfig
	}
	return nil
}

type GenerateWorkflowRequest struct {
	// The project to generate a workflow config for.
	//
	// The project name used must be known to Tricium.
	Project string `protobuf:"bytes,1,opt,name=project" json:"project,omitempty"`
	// The paths to generate the workflow config.
	//
	// This list of file metadata includes file paths which are used to
	// decide which workers to include in the workflow.
	Files                []*v1.Data_File `protobuf:"bytes,2,rep,name=files" json:"files,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *GenerateWorkflowRequest) Reset()         { *m = GenerateWorkflowRequest{} }
func (m *GenerateWorkflowRequest) String() string { return proto.CompactTextString(m) }
func (*GenerateWorkflowRequest) ProtoMessage()    {}
func (*GenerateWorkflowRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_043929b00c61711a, []int{2}
}
func (m *GenerateWorkflowRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateWorkflowRequest.Unmarshal(m, b)
}
func (m *GenerateWorkflowRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateWorkflowRequest.Marshal(b, m, deterministic)
}
func (dst *GenerateWorkflowRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateWorkflowRequest.Merge(dst, src)
}
func (m *GenerateWorkflowRequest) XXX_Size() int {
	return xxx_messageInfo_GenerateWorkflowRequest.Size(m)
}
func (m *GenerateWorkflowRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateWorkflowRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateWorkflowRequest proto.InternalMessageInfo

func (m *GenerateWorkflowRequest) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *GenerateWorkflowRequest) GetFiles() []*v1.Data_File {
	if m != nil {
		return m.Files
	}
	return nil
}

type GenerateWorkflowResponse struct {
	// The generated workflow.
	Workflow             *Workflow `protobuf:"bytes,1,opt,name=workflow" json:"workflow,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GenerateWorkflowResponse) Reset()         { *m = GenerateWorkflowResponse{} }
func (m *GenerateWorkflowResponse) String() string { return proto.CompactTextString(m) }
func (*GenerateWorkflowResponse) ProtoMessage()    {}
func (*GenerateWorkflowResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_043929b00c61711a, []int{3}
}
func (m *GenerateWorkflowResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateWorkflowResponse.Unmarshal(m, b)
}
func (m *GenerateWorkflowResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateWorkflowResponse.Marshal(b, m, deterministic)
}
func (dst *GenerateWorkflowResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateWorkflowResponse.Merge(dst, src)
}
func (m *GenerateWorkflowResponse) XXX_Size() int {
	return xxx_messageInfo_GenerateWorkflowResponse.Size(m)
}
func (m *GenerateWorkflowResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateWorkflowResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateWorkflowResponse proto.InternalMessageInfo

func (m *GenerateWorkflowResponse) GetWorkflow() *Workflow {
	if m != nil {
		return m.Workflow
	}
	return nil
}

func init() {
	proto.RegisterType((*ValidateRequest)(nil), "admin.ValidateRequest")
	proto.RegisterType((*ValidateResponse)(nil), "admin.ValidateResponse")
	proto.RegisterType((*GenerateWorkflowRequest)(nil), "admin.GenerateWorkflowRequest")
	proto.RegisterType((*GenerateWorkflowResponse)(nil), "admin.GenerateWorkflowResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConfigClient is the client API for Config service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConfigClient interface {
	// Validates a Tricium config.
	//
	// The config to validate is specified in the request.
	// TODO(emso): Make this RPC public to let users validate configs when they
	// want, or via luci-config.
	Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
	// Generates a workflow config from a Tricium config.
	//
	// The Tricium config to generate for is specified by the request.
	GenerateWorkflow(ctx context.Context, in *GenerateWorkflowRequest, opts ...grpc.CallOption) (*GenerateWorkflowResponse, error)
}
type configPRPCClient struct {
	client *prpc.Client
}

func NewConfigPRPCClient(client *prpc.Client) ConfigClient {
	return &configPRPCClient{client}
}

func (c *configPRPCClient) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.client.Call(ctx, "admin.Config", "Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configPRPCClient) GenerateWorkflow(ctx context.Context, in *GenerateWorkflowRequest, opts ...grpc.CallOption) (*GenerateWorkflowResponse, error) {
	out := new(GenerateWorkflowResponse)
	err := c.client.Call(ctx, "admin.Config", "GenerateWorkflow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type configClient struct {
	cc *grpc.ClientConn
}

func NewConfigClient(cc *grpc.ClientConn) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, "/admin.Config/Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) GenerateWorkflow(ctx context.Context, in *GenerateWorkflowRequest, opts ...grpc.CallOption) (*GenerateWorkflowResponse, error) {
	out := new(GenerateWorkflowResponse)
	err := c.cc.Invoke(ctx, "/admin.Config/GenerateWorkflow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServer is the server API for Config service.
type ConfigServer interface {
	// Validates a Tricium config.
	//
	// The config to validate is specified in the request.
	// TODO(emso): Make this RPC public to let users validate configs when they
	// want, or via luci-config.
	Validate(context.Context, *ValidateRequest) (*ValidateResponse, error)
	// Generates a workflow config from a Tricium config.
	//
	// The Tricium config to generate for is specified by the request.
	GenerateWorkflow(context.Context, *GenerateWorkflowRequest) (*GenerateWorkflowResponse, error)
}

func RegisterConfigServer(s prpc.Registrar, srv ConfigServer) {
	s.RegisterService(&_Config_serviceDesc, srv)
}

func _Config_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin.Config/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).Validate(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_GenerateWorkflow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateWorkflowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).GenerateWorkflow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin.Config/GenerateWorkflow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).GenerateWorkflow(ctx, req.(*GenerateWorkflowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Config_serviceDesc = grpc.ServiceDesc{
	ServiceName: "admin.Config",
	HandlerType: (*ConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Validate",
			Handler:    _Config_Validate_Handler,
		},
		{
			MethodName: "GenerateWorkflow",
			Handler:    _Config_GenerateWorkflow_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infra/tricium/api/admin/v1/config.proto",
}

func init() {
	proto.RegisterFile("infra/tricium/api/admin/v1/config.proto", fileDescriptor_config_043929b00c61711a)
}

var fileDescriptor_config_043929b00c61711a = []byte{
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x4e, 0xc2, 0x40,
	0x10, 0x4d, 0x31, 0x20, 0x0e, 0x51, 0xc8, 0x1e, 0xa0, 0xe9, 0x41, 0x90, 0x8b, 0x18, 0x93, 0x36,
	0xe2, 0xd1, 0x78, 0x30, 0x1a, 0xb9, 0x9a, 0x12, 0xf5, 0x64, 0xcc, 0xda, 0x4e, 0xcd, 0x6a, 0xe9,
	0xd6, 0xdd, 0x05, 0x3e, 0xc3, 0xbb, 0x5f, 0x6b, 0xe8, 0xee, 0x56, 0x52, 0xc4, 0x78, 0x9c, 0x99,
	0xf7, 0xde, 0xbe, 0x79, 0xb3, 0x70, 0xcc, 0xb2, 0x44, 0xd0, 0x40, 0x09, 0x16, 0xb1, 0xf9, 0x2c,
	0xa0, 0x39, 0x0b, 0x68, 0x3c, 0x63, 0x59, 0xb0, 0x38, 0x0b, 0x22, 0x9e, 0x25, 0xec, 0xd5, 0xcf,
	0x05, 0x57, 0x9c, 0xd4, 0x8b, 0xb6, 0x77, 0xf2, 0x07, 0x7e, 0xc9, 0xc5, 0x7b, 0x92, 0xf2, 0xa5,
	0x66, 0x78, 0x47, 0x9b, 0xd0, 0x8a, 0xa8, 0xd7, 0xff, 0x15, 0x12, 0x53, 0x45, 0x35, 0x60, 0xf8,
	0xe9, 0x40, 0xfb, 0x81, 0xa6, 0x2c, 0xa6, 0x0a, 0x43, 0xfc, 0x98, 0xa3, 0x54, 0xe4, 0x12, 0x0e,
	0x72, 0xc1, 0xdf, 0x30, 0x52, 0xcf, 0x5a, 0xcc, 0x75, 0x06, 0xce, 0xa8, 0x35, 0xee, 0xfa, 0x46,
	0xc7, 0xbf, 0xd3, 0xe3, 0xeb, 0x62, 0x1a, 0xee, 0xe7, 0xeb, 0xe5, 0x8a, 0x2e, 0x51, 0x2c, 0x58,
	0x84, 0x96, 0x5e, 0xab, 0xd0, 0xa7, 0x7a, 0x6c, 0xe9, 0x72, 0xbd, 0x1c, 0xde, 0x43, 0xe7, 0xc7,
	0x90, 0xcc, 0x79, 0x26, 0x91, 0x5c, 0x41, 0x67, 0x61, 0x7a, 0xf1, 0xff, 0x3c, 0xb5, 0x4b, 0xbc,
	0x91, 0x7d, 0x82, 0xde, 0x04, 0x33, 0x14, 0x54, 0xe1, 0xa3, 0x89, 0xd1, 0xee, 0xeb, 0xc2, 0xae,
	0xd9, 0xa0, 0x10, 0xdd, 0x0b, 0x6d, 0x49, 0x46, 0x50, 0x4f, 0x58, 0x8a, 0xd2, 0xad, 0x0d, 0x76,
	0x46, 0xad, 0x31, 0x29, 0x1f, 0xbb, 0x59, 0x25, 0x78, 0xcb, 0x52, 0x0c, 0x35, 0x60, 0x38, 0x01,
	0x77, 0x53, 0xde, 0xb8, 0x3f, 0x85, 0xa6, 0xbd, 0x9c, 0x71, 0xdd, 0xf6, 0x8b, 0x9b, 0xfa, 0x25,
	0xb4, 0x04, 0x8c, 0xbf, 0x1c, 0x68, 0x98, 0x20, 0x2f, 0xa0, 0x69, 0x93, 0x20, 0x5d, 0xc3, 0xa8,
	0xdc, 0xca, 0xeb, 0x6d, 0xf4, 0xcd, 0xa3, 0x53, 0xe8, 0x54, 0x0d, 0x91, 0x43, 0x03, 0xde, 0x12,
	0x84, 0xd7, 0xdf, 0x3a, 0xd7, 0xa2, 0x2f, 0x8d, 0xe2, 0xd3, 0x9c, 0x7f, 0x07, 0x00, 0x00, 0xff,
	0xff, 0x82, 0x7c, 0xd3, 0xb1, 0xd5, 0x02, 0x00, 0x00,
}
