// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/qscheduler/qslib/protos/config.proto

package protos

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
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

// SchedulerConfig represents configuration information about the behavior of
// accounts for this quota scheduler pool.
type SchedulerConfig struct {
	// Configuration for a given account, keyed by account id.
	AccountConfigs map[string]*AccountConfig `protobuf:"bytes,1,rep,name=account_configs,json=accountConfigs,proto3" json:"account_configs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// If set, scheduler will never preempt running tasks.
	DisablePreemption    bool     `protobuf:"varint,2,opt,name=disable_preemption,json=disablePreemption,proto3" json:"disable_preemption,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SchedulerConfig) Reset()         { *m = SchedulerConfig{} }
func (m *SchedulerConfig) String() string { return proto.CompactTextString(m) }
func (*SchedulerConfig) ProtoMessage()    {}
func (*SchedulerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_9585cdf0027c1e37, []int{0}
}

func (m *SchedulerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SchedulerConfig.Unmarshal(m, b)
}
func (m *SchedulerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SchedulerConfig.Marshal(b, m, deterministic)
}
func (m *SchedulerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SchedulerConfig.Merge(m, src)
}
func (m *SchedulerConfig) XXX_Size() int {
	return xxx_messageInfo_SchedulerConfig.Size(m)
}
func (m *SchedulerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_SchedulerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_SchedulerConfig proto.InternalMessageInfo

func (m *SchedulerConfig) GetAccountConfigs() map[string]*AccountConfig {
	if m != nil {
		return m.AccountConfigs
	}
	return nil
}

func (m *SchedulerConfig) GetDisablePreemption() bool {
	if m != nil {
		return m.DisablePreemption
	}
	return false
}

// AccountConfig represents per-quota-account configuration information, such
// as the recharge parameters. This does not represent anything about the
// current state of an account.
type AccountConfig struct {
	// ChargeRate is the rates (per second) at which per-priority accounts grow.
	//
	// Conceptually this is the time-averaged number of workers that this account
	// may use, at each priority level.
	ChargeRate []float32 `protobuf:"fixed32,1,rep,packed,name=charge_rate,json=chargeRate,proto3" json:"charge_rate,omitempty"`
	// MaxChargeSeconds is the maximum amount of time over which this account can
	// accumulate quota before hitting its cap.
	//
	// Conceptually this sets the time window over which the time averaged
	// utilization by this account is measured. Very bursty clients will need to
	// use a wider window, whereas very consistent clients will use a narrow one.
	MaxChargeSeconds float32 `protobuf:"fixed32,2,opt,name=max_charge_seconds,json=maxChargeSeconds,proto3" json:"max_charge_seconds,omitempty"`
	// MaxFanout is the maximum number of concurrent paid jobs that this account
	// will pay for (0 = no limit).
	//
	// Additional jobs beyond this may run if there is idle capacity, but they
	// will run in the FreeBucket priority level.
	MaxFanout            int32    `protobuf:"varint,3,opt,name=max_fanout,json=maxFanout,proto3" json:"max_fanout,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountConfig) Reset()         { *m = AccountConfig{} }
func (m *AccountConfig) String() string { return proto.CompactTextString(m) }
func (*AccountConfig) ProtoMessage()    {}
func (*AccountConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_9585cdf0027c1e37, []int{1}
}

func (m *AccountConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountConfig.Unmarshal(m, b)
}
func (m *AccountConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountConfig.Marshal(b, m, deterministic)
}
func (m *AccountConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountConfig.Merge(m, src)
}
func (m *AccountConfig) XXX_Size() int {
	return xxx_messageInfo_AccountConfig.Size(m)
}
func (m *AccountConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountConfig.DiscardUnknown(m)
}

var xxx_messageInfo_AccountConfig proto.InternalMessageInfo

func (m *AccountConfig) GetChargeRate() []float32 {
	if m != nil {
		return m.ChargeRate
	}
	return nil
}

func (m *AccountConfig) GetMaxChargeSeconds() float32 {
	if m != nil {
		return m.MaxChargeSeconds
	}
	return 0
}

func (m *AccountConfig) GetMaxFanout() int32 {
	if m != nil {
		return m.MaxFanout
	}
	return 0
}

func init() {
	proto.RegisterType((*SchedulerConfig)(nil), "protos.SchedulerConfig")
	proto.RegisterMapType((map[string]*AccountConfig)(nil), "protos.SchedulerConfig.AccountConfigsEntry")
	proto.RegisterType((*AccountConfig)(nil), "protos.AccountConfig")
}

func init() {
	proto.RegisterFile("infra/qscheduler/qslib/protos/config.proto", fileDescriptor_9585cdf0027c1e37)
}

var fileDescriptor_9585cdf0027c1e37 = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4f, 0xc2, 0x30,
	0x14, 0xc6, 0xd3, 0x0d, 0x0c, 0x3c, 0xa2, 0x60, 0x8d, 0xc9, 0x62, 0x62, 0x5c, 0x38, 0x2d, 0xa2,
	0x2c, 0xc1, 0x8b, 0xf1, 0x66, 0x88, 0x1e, 0x3c, 0x99, 0xe2, 0xc1, 0xdb, 0xf2, 0x28, 0x05, 0x16,
	0xb7, 0x16, 0xda, 0xce, 0xc0, 0xc1, 0x7f, 0xdb, 0xb3, 0xa1, 0x05, 0x93, 0x19, 0x6f, 0x7b, 0xbf,
	0xef, 0xf7, 0xbd, 0xb5, 0x85, 0xeb, 0x5c, 0xce, 0x35, 0xa6, 0x6b, 0xc3, 0x97, 0x62, 0x56, 0x15,
	0x42, 0xa7, 0x6b, 0x53, 0xe4, 0xd3, 0x74, 0xa5, 0x95, 0x55, 0x26, 0xe5, 0x4a, 0xce, 0xf3, 0xc5,
	0xd0, 0x4d, 0xf4, 0xc8, 0xc3, 0xfe, 0x37, 0x81, 0xee, 0xe4, 0xe0, 0x8f, 0x9d, 0x41, 0xdf, 0xa0,
	0x8b, 0x9c, 0xab, 0x4a, 0xda, 0xcc, 0x77, 0x4c, 0x44, 0xe2, 0x30, 0xe9, 0x8c, 0x06, 0xbe, 0x6c,
	0x86, 0x7f, 0x1a, 0xc3, 0x47, 0xaf, 0xfb, 0xc9, 0x3c, 0x49, 0xab, 0xb7, 0xec, 0x04, 0x6b, 0x90,
	0xde, 0x02, 0x9d, 0xe5, 0x06, 0xa7, 0x85, 0xc8, 0x56, 0x5a, 0x88, 0x72, 0x65, 0x73, 0x25, 0xa3,
	0x20, 0x26, 0x49, 0x8b, 0x9d, 0xee, 0x93, 0xd7, 0xdf, 0xe0, 0xe2, 0x1d, 0xce, 0xfe, 0xd9, 0x4a,
	0x7b, 0x10, 0x7e, 0x88, 0x6d, 0x44, 0x62, 0x92, 0xb4, 0xd9, 0xee, 0x93, 0x0e, 0xa0, 0xf9, 0x89,
	0x45, 0x25, 0xdc, 0xaa, 0xce, 0xe8, 0xfc, 0x70, 0xc6, 0x5a, 0x9b, 0x79, 0xe7, 0x21, 0xb8, 0x27,
	0x2f, 0x8d, 0x56, 0xd8, 0x6b, 0xf4, 0xbf, 0xe0, 0xb8, 0x66, 0xd0, 0x2b, 0xe8, 0xf0, 0x25, 0xea,
	0x85, 0xc8, 0x34, 0x5a, 0xe1, 0x6e, 0x1c, 0x30, 0xf0, 0x88, 0xa1, 0x15, 0xf4, 0x06, 0x68, 0x89,
	0x9b, 0x6c, 0x2f, 0x19, 0xc1, 0x95, 0x9c, 0x19, 0xf7, 0xd7, 0x80, 0xf5, 0x4a, 0xdc, 0x8c, 0x5d,
	0x30, 0xf1, 0x9c, 0x5e, 0x02, 0xec, 0xec, 0x39, 0x4a, 0x55, 0xd9, 0x28, 0x8c, 0x49, 0xd2, 0x64,
	0xed, 0x12, 0x37, 0xcf, 0x0e, 0x4c, 0xfd, 0xfb, 0xdf, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0x5d,
	0x48, 0xeb, 0x35, 0xb4, 0x01, 0x00, 0x00,
}
