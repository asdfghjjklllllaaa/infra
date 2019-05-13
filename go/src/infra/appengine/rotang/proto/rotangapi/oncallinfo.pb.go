// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/appengine/rotang/proto/rotangapi/oncallinfo.proto

package rotangapi

import prpc "go.chromium.org/luci/grpc/prpc"

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// OncallRequest contains the name of the rotation of interest.
type OncallRequest struct {
	// name is a required field containing the name of the rotation
	// of interest.
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OncallRequest) Reset()         { *m = OncallRequest{} }
func (m *OncallRequest) String() string { return proto.CompactTextString(m) }
func (*OncallRequest) ProtoMessage()    {}
func (*OncallRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{0}
}

func (m *OncallRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OncallRequest.Unmarshal(m, b)
}
func (m *OncallRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OncallRequest.Marshal(b, m, deterministic)
}
func (m *OncallRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OncallRequest.Merge(m, src)
}
func (m *OncallRequest) XXX_Size() int {
	return xxx_messageInfo_OncallRequest.Size(m)
}
func (m *OncallRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OncallRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OncallRequest proto.InternalMessageInfo

func (m *OncallRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// OncallResponse contains the current on-callers for a rotation.
type OncallResponse struct {
	// shift contains the current shift entry for the specified rotation.
	Shift *ShiftEntry `protobuf:"bytes,1,opt,name=shift,proto3" json:"shift,omitempty"`
	// tz_consider indicates if the rotation generator considers the TimeZone
	// of members when scheduling memmbers.
	// A rotation using a Generator considering timezones will generatate a shift
	// with one on-caller per TZ of their members.
	// Eg. if a rotation have members:
	//
	// Australia/Sydney:
	//  syd1@oncall.com
	//  syd2@oncall.com
	// EST
	//  est1@oncall.com
	//  est2@oncall.com
	// US/Pacific
	//  mtv1@oncall.com
	//  mtv2@oncall.com
	// UTC
	//  utc1@oncall.com
	//  utc2@oncall.com
	//
	// A rotation using a tz_consider generator would generate a shift with one
	// one member from each TZ.
	//   syd1@oncall.com , est1@oncall.com , mtv1@oncall.com and utc1@oncall.com
	TzConsider           bool     `protobuf:"varint,2,opt,name=tz_consider,json=tzConsider,proto3" json:"tz_consider,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OncallResponse) Reset()         { *m = OncallResponse{} }
func (m *OncallResponse) String() string { return proto.CompactTextString(m) }
func (*OncallResponse) ProtoMessage()    {}
func (*OncallResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{1}
}

func (m *OncallResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OncallResponse.Unmarshal(m, b)
}
func (m *OncallResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OncallResponse.Marshal(b, m, deterministic)
}
func (m *OncallResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OncallResponse.Merge(m, src)
}
func (m *OncallResponse) XXX_Size() int {
	return xxx_messageInfo_OncallResponse.Size(m)
}
func (m *OncallResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OncallResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OncallResponse proto.InternalMessageInfo

func (m *OncallResponse) GetShift() *ShiftEntry {
	if m != nil {
		return m.Shift
	}
	return nil
}

func (m *OncallResponse) GetTzConsider() bool {
	if m != nil {
		return m.TzConsider
	}
	return false
}

// ListRotationsRequest is used to get a list of all configured rotations.
type ListRotationsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRotationsRequest) Reset()         { *m = ListRotationsRequest{} }
func (m *ListRotationsRequest) String() string { return proto.CompactTextString(m) }
func (*ListRotationsRequest) ProtoMessage()    {}
func (*ListRotationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{2}
}

func (m *ListRotationsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRotationsRequest.Unmarshal(m, b)
}
func (m *ListRotationsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRotationsRequest.Marshal(b, m, deterministic)
}
func (m *ListRotationsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRotationsRequest.Merge(m, src)
}
func (m *ListRotationsRequest) XXX_Size() int {
	return xxx_messageInfo_ListRotationsRequest.Size(m)
}
func (m *ListRotationsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRotationsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRotationsRequest proto.InternalMessageInfo

// ListRotationsResponse contains the configured rotations.
type ListRotationsResponse struct {
	Rotations            []*Rotation `protobuf:"bytes,1,rep,name=rotations,proto3" json:"rotations,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ListRotationsResponse) Reset()         { *m = ListRotationsResponse{} }
func (m *ListRotationsResponse) String() string { return proto.CompactTextString(m) }
func (*ListRotationsResponse) ProtoMessage()    {}
func (*ListRotationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{3}
}

func (m *ListRotationsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRotationsResponse.Unmarshal(m, b)
}
func (m *ListRotationsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRotationsResponse.Marshal(b, m, deterministic)
}
func (m *ListRotationsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRotationsResponse.Merge(m, src)
}
func (m *ListRotationsResponse) XXX_Size() int {
	return xxx_messageInfo_ListRotationsResponse.Size(m)
}
func (m *ListRotationsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRotationsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListRotationsResponse proto.InternalMessageInfo

func (m *ListRotationsResponse) GetRotations() []*Rotation {
	if m != nil {
		return m.Rotations
	}
	return nil
}

// ShiftsRequest defines shifts of interest to fetch.
type ShiftsRequest struct {
	// name is a required field containing the rotation name of interest.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// start contains the beginning of the time range of interest.
	// Leaving this field empty will fetch shifts from the beginning of time.
	Start *timestamp.Timestamp `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
	// end contains the end of the time range of interest.
	// Leaving this field empty will fetch shift to the end of time.
	End                  *timestamp.Timestamp `protobuf:"bytes,3,opt,name=end,proto3" json:"end,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ShiftsRequest) Reset()         { *m = ShiftsRequest{} }
func (m *ShiftsRequest) String() string { return proto.CompactTextString(m) }
func (*ShiftsRequest) ProtoMessage()    {}
func (*ShiftsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{4}
}

func (m *ShiftsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShiftsRequest.Unmarshal(m, b)
}
func (m *ShiftsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShiftsRequest.Marshal(b, m, deterministic)
}
func (m *ShiftsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShiftsRequest.Merge(m, src)
}
func (m *ShiftsRequest) XXX_Size() int {
	return xxx_messageInfo_ShiftsRequest.Size(m)
}
func (m *ShiftsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShiftsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShiftsRequest proto.InternalMessageInfo

func (m *ShiftsRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShiftsRequest) GetStart() *timestamp.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *ShiftsRequest) GetEnd() *timestamp.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

// ShiftsResponse contains the shifts requested by ShiftsRequest.
type ShiftsResponse struct {
	Shifts               []*ShiftEntry `protobuf:"bytes,1,rep,name=shifts,proto3" json:"shifts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ShiftsResponse) Reset()         { *m = ShiftsResponse{} }
func (m *ShiftsResponse) String() string { return proto.CompactTextString(m) }
func (*ShiftsResponse) ProtoMessage()    {}
func (*ShiftsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{5}
}

func (m *ShiftsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShiftsResponse.Unmarshal(m, b)
}
func (m *ShiftsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShiftsResponse.Marshal(b, m, deterministic)
}
func (m *ShiftsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShiftsResponse.Merge(m, src)
}
func (m *ShiftsResponse) XXX_Size() int {
	return xxx_messageInfo_ShiftsResponse.Size(m)
}
func (m *ShiftsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ShiftsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ShiftsResponse proto.InternalMessageInfo

func (m *ShiftsResponse) GetShifts() []*ShiftEntry {
	if m != nil {
		return m.Shifts
	}
	return nil
}

// Shift defines a shift configuration.
// RotaNG supports rotations with multiple shifts.
// Eg. follow-the-sun configurations:
// Shift:
//   SYD: 8hours
//   MTV: 8hours
//   EY:  8hours.
type ShiftConfiguration struct {
	Name                 string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Duration             *duration.Duration `protobuf:"bytes,2,opt,name=duration,proto3" json:"duration,omitempty"`
	Entries              []*ShiftEntry      `protobuf:"bytes,3,rep,name=entries,proto3" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ShiftConfiguration) Reset()         { *m = ShiftConfiguration{} }
func (m *ShiftConfiguration) String() string { return proto.CompactTextString(m) }
func (*ShiftConfiguration) ProtoMessage()    {}
func (*ShiftConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{6}
}

func (m *ShiftConfiguration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShiftConfiguration.Unmarshal(m, b)
}
func (m *ShiftConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShiftConfiguration.Marshal(b, m, deterministic)
}
func (m *ShiftConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShiftConfiguration.Merge(m, src)
}
func (m *ShiftConfiguration) XXX_Size() int {
	return xxx_messageInfo_ShiftConfiguration.Size(m)
}
func (m *ShiftConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_ShiftConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_ShiftConfiguration proto.InternalMessageInfo

func (m *ShiftConfiguration) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShiftConfiguration) GetDuration() *duration.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

func (m *ShiftConfiguration) GetEntries() []*ShiftEntry {
	if m != nil {
		return m.Entries
	}
	return nil
}

// ShiftEntry defines a single shift.
type ShiftEntry struct {
	// name is the Shift configuration this entry belongs to.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// oncallers is the email addresses of the members on-call
	// for this shift.
	Oncallers []*OnCaller `protobuf:"bytes,2,rep,name=oncallers,proto3" json:"oncallers,omitempty"`
	// start is the start of this shift.
	Start *timestamp.Timestamp `protobuf:"bytes,3,opt,name=start,proto3" json:"start,omitempty"`
	// end is the end of the shift.
	End *timestamp.Timestamp `protobuf:"bytes,4,opt,name=end,proto3" json:"end,omitempty"`
	// comment contains an optional comment about this shift.
	// Eg. Information about a shift swap.
	Comment string `protobuf:"bytes,5,opt,name=comment,proto3" json:"comment,omitempty"`
	// event_id contains the Google Calendar Event ID of this shift.
	// Only enabled rotations will have event_id's.
	EventId              string   `protobuf:"bytes,6,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShiftEntry) Reset()         { *m = ShiftEntry{} }
func (m *ShiftEntry) String() string { return proto.CompactTextString(m) }
func (*ShiftEntry) ProtoMessage()    {}
func (*ShiftEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{7}
}

func (m *ShiftEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShiftEntry.Unmarshal(m, b)
}
func (m *ShiftEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShiftEntry.Marshal(b, m, deterministic)
}
func (m *ShiftEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShiftEntry.Merge(m, src)
}
func (m *ShiftEntry) XXX_Size() int {
	return xxx_messageInfo_ShiftEntry.Size(m)
}
func (m *ShiftEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_ShiftEntry.DiscardUnknown(m)
}

var xxx_messageInfo_ShiftEntry proto.InternalMessageInfo

func (m *ShiftEntry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ShiftEntry) GetOncallers() []*OnCaller {
	if m != nil {
		return m.Oncallers
	}
	return nil
}

func (m *ShiftEntry) GetStart() *timestamp.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *ShiftEntry) GetEnd() *timestamp.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

func (m *ShiftEntry) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func (m *ShiftEntry) GetEventId() string {
	if m != nil {
		return m.EventId
	}
	return ""
}

// OnCaller contains one member on-call for a shift.
type OnCaller struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Tz                   string   `protobuf:"bytes,3,opt,name=tz,proto3" json:"tz,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OnCaller) Reset()         { *m = OnCaller{} }
func (m *OnCaller) String() string { return proto.CompactTextString(m) }
func (*OnCaller) ProtoMessage()    {}
func (*OnCaller) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{8}
}

func (m *OnCaller) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OnCaller.Unmarshal(m, b)
}
func (m *OnCaller) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OnCaller.Marshal(b, m, deterministic)
}
func (m *OnCaller) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OnCaller.Merge(m, src)
}
func (m *OnCaller) XXX_Size() int {
	return xxx_messageInfo_OnCaller.Size(m)
}
func (m *OnCaller) XXX_DiscardUnknown() {
	xxx_messageInfo_OnCaller.DiscardUnknown(m)
}

var xxx_messageInfo_OnCaller proto.InternalMessageInfo

func (m *OnCaller) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *OnCaller) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *OnCaller) GetTz() string {
	if m != nil {
		return m.Tz
	}
	return ""
}

// Rotation contains the Rotation information returned by the List
// RPC call.
type Rotation struct {
	// name contain the name of the rotation.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// enabled signifies if the rotation is enabled.
	// A rotation needs to be enabled to create/modify Google calendar
	// events and sending e-mails.
	Enabled              bool     `protobuf:"varint,2,opt,name=enabled,proto3" json:"enabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Rotation) Reset()         { *m = Rotation{} }
func (m *Rotation) String() string { return proto.CompactTextString(m) }
func (*Rotation) ProtoMessage()    {}
func (*Rotation) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b351329bbdaa59d, []int{9}
}

func (m *Rotation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Rotation.Unmarshal(m, b)
}
func (m *Rotation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Rotation.Marshal(b, m, deterministic)
}
func (m *Rotation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Rotation.Merge(m, src)
}
func (m *Rotation) XXX_Size() int {
	return xxx_messageInfo_Rotation.Size(m)
}
func (m *Rotation) XXX_DiscardUnknown() {
	xxx_messageInfo_Rotation.DiscardUnknown(m)
}

var xxx_messageInfo_Rotation proto.InternalMessageInfo

func (m *Rotation) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Rotation) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func init() {
	proto.RegisterType((*OncallRequest)(nil), "rotangapi.OncallRequest")
	proto.RegisterType((*OncallResponse)(nil), "rotangapi.OncallResponse")
	proto.RegisterType((*ListRotationsRequest)(nil), "rotangapi.ListRotationsRequest")
	proto.RegisterType((*ListRotationsResponse)(nil), "rotangapi.ListRotationsResponse")
	proto.RegisterType((*ShiftsRequest)(nil), "rotangapi.ShiftsRequest")
	proto.RegisterType((*ShiftsResponse)(nil), "rotangapi.ShiftsResponse")
	proto.RegisterType((*ShiftConfiguration)(nil), "rotangapi.ShiftConfiguration")
	proto.RegisterType((*ShiftEntry)(nil), "rotangapi.ShiftEntry")
	proto.RegisterType((*OnCaller)(nil), "rotangapi.OnCaller")
	proto.RegisterType((*Rotation)(nil), "rotangapi.Rotation")
}

func init() {
	proto.RegisterFile("infra/appengine/rotang/proto/rotangapi/oncallinfo.proto", fileDescriptor_7b351329bbdaa59d)
}

var fileDescriptor_7b351329bbdaa59d = []byte{
	// 544 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xcf, 0x6f, 0x12, 0x41,
	0x14, 0xc7, 0xbb, 0x50, 0x28, 0x3c, 0x02, 0x87, 0xb1, 0x35, 0x03, 0x07, 0x21, 0xe3, 0x85, 0x44,
	0xdd, 0xb5, 0x18, 0xa3, 0x37, 0x0e, 0xd4, 0x43, 0x8d, 0x49, 0x93, 0xb5, 0x67, 0x9b, 0x81, 0x9d,
	0xc5, 0x49, 0xd8, 0x99, 0x75, 0x67, 0x30, 0x91, 0xab, 0xff, 0x80, 0x7f, 0xa4, 0x07, 0xff, 0x0c,
	0xc3, 0xfc, 0x00, 0xb6, 0xdd, 0x56, 0x7b, 0xe3, 0xbd, 0xf7, 0xd9, 0xe1, 0xfb, 0x7d, 0x3f, 0xe0,
	0x1d, 0x17, 0x69, 0x41, 0x23, 0x9a, 0xe7, 0x4c, 0x2c, 0xb9, 0x60, 0x51, 0x21, 0x35, 0x15, 0xcb,
	0x28, 0x2f, 0xa4, 0x96, 0x2e, 0xa0, 0x39, 0x8f, 0xa4, 0x58, 0xd0, 0xd5, 0x8a, 0x8b, 0x54, 0x86,
	0xa6, 0x84, 0xda, 0xbb, 0xda, 0x60, 0xb8, 0x94, 0x72, 0xb9, 0x62, 0xf6, 0x9b, 0xf9, 0x3a, 0x8d,
	0x34, 0xcf, 0x98, 0xd2, 0x34, 0xcb, 0x2d, 0x3b, 0x78, 0x76, 0x1b, 0x48, 0xd6, 0x05, 0xd5, 0x5c,
	0x0a, 0x5b, 0x27, 0xcf, 0xa1, 0x7b, 0x65, 0xde, 0x8f, 0xd9, 0xb7, 0x35, 0x53, 0x1a, 0x21, 0x38,
	0x16, 0x34, 0x63, 0x38, 0x18, 0x05, 0xe3, 0x76, 0x6c, 0x7e, 0x93, 0x2f, 0xd0, 0xf3, 0x90, 0xca,
	0xa5, 0x50, 0x0c, 0xbd, 0x80, 0x86, 0xfa, 0xca, 0x53, 0x6d, 0xb0, 0xce, 0xe4, 0x2c, 0xdc, 0x49,
	0x0a, 0x3f, 0x6f, 0xf3, 0x1f, 0x84, 0x2e, 0x7e, 0xc4, 0x96, 0x41, 0x43, 0xe8, 0xe8, 0xcd, 0xcd,
	0x42, 0x0a, 0xc5, 0x13, 0x56, 0xe0, 0xda, 0x28, 0x18, 0xb7, 0x62, 0xd0, 0x9b, 0x99, 0xcb, 0x90,
	0xa7, 0x70, 0xfa, 0x89, 0x2b, 0x1d, 0x4b, 0x6d, 0xa4, 0x29, 0xa7, 0x85, 0x7c, 0x84, 0xb3, 0x5b,
	0x79, 0xf7, 0xf7, 0xe7, 0x60, 0x7a, 0x60, 0x92, 0x38, 0x18, 0xd5, 0xc7, 0x9d, 0xc9, 0x93, 0x03,
	0x09, 0xfe, 0x83, 0x78, 0x4f, 0x91, 0x9f, 0x01, 0x74, 0x8d, 0x34, 0xf5, 0x80, 0x53, 0xf4, 0x1a,
	0x1a, 0x4a, 0xd3, 0x42, 0x1b, 0x91, 0x9d, 0xc9, 0x20, 0xb4, 0xed, 0x0b, 0x7d, 0xfb, 0xc2, 0x6b,
	0xdf, 0xdf, 0xd8, 0x82, 0xe8, 0x25, 0xd4, 0x99, 0x48, 0x70, 0xfd, 0x9f, 0xfc, 0x16, 0x23, 0x53,
	0xe8, 0x79, 0x11, 0xce, 0xca, 0x2b, 0x68, 0x9a, 0x2e, 0x79, 0x1f, 0xf7, 0xb4, 0xd2, 0x41, 0xe4,
	0x57, 0x00, 0xc8, 0xa4, 0x67, 0x52, 0xa4, 0x7c, 0xe9, 0x86, 0x59, 0xe9, 0xe5, 0x2d, 0xb4, 0xfc,
	0xb0, 0x9d, 0x9d, 0xfe, 0x1d, 0x79, 0x17, 0x0e, 0x88, 0x77, 0x28, 0x8a, 0xe0, 0x84, 0x09, 0x5d,
	0x70, 0xa6, 0x70, 0xfd, 0x21, 0x45, 0x9e, 0x22, 0x7f, 0x02, 0x80, 0x7d, 0xbe, 0x52, 0xca, 0x39,
	0xb4, 0xed, 0x16, 0xb3, 0x42, 0xe1, 0xda, 0x9d, 0x79, 0x5d, 0x89, 0x99, 0xa9, 0xc5, 0x7b, 0x6a,
	0x3f, 0x89, 0xfa, 0x23, 0x27, 0x71, 0xfc, 0x5f, 0x93, 0x40, 0x18, 0x4e, 0x16, 0x32, 0xcb, 0x98,
	0xd0, 0xb8, 0x61, 0x94, 0xfa, 0x10, 0xf5, 0xa1, 0xc5, 0xbe, 0x33, 0xa1, 0x6f, 0x78, 0x82, 0x9b,
	0xb6, 0x64, 0xe2, 0xcb, 0x84, 0x5c, 0x40, 0xcb, 0x6b, 0x45, 0xa7, 0xd0, 0x60, 0x19, 0xe5, 0x2b,
	0x67, 0xd4, 0x06, 0x3b, 0xf7, 0xb5, 0x03, 0xf7, 0x3d, 0xa8, 0xe9, 0x8d, 0xf1, 0xd1, 0x8e, 0x6b,
	0x7a, 0x43, 0xde, 0x43, 0xcb, 0x6f, 0x68, 0x65, 0xb7, 0xf0, 0x76, 0x02, 0x74, 0xbe, 0x62, 0x89,
	0xbb, 0x15, 0x1f, 0x4e, 0x7e, 0x07, 0x00, 0xf6, 0x12, 0x2f, 0x45, 0x2a, 0xd1, 0x14, 0x9a, 0x36,
	0x42, 0xb8, 0xd4, 0xcd, 0x83, 0x7b, 0x1e, 0xf4, 0x2b, 0x2a, 0x76, 0xf5, 0xc8, 0x11, 0xba, 0x86,
	0x6e, 0xe9, 0xc0, 0xd0, 0xf0, 0x80, 0xae, 0x3a, 0xc9, 0xc1, 0xe8, 0x7e, 0x60, 0xf7, 0xea, 0x14,
	0x9a, 0x76, 0xc9, 0x4b, 0xb2, 0x4a, 0xc7, 0x57, 0x92, 0x55, 0xbe, 0x08, 0x72, 0x34, 0x6f, 0x9a,
	0xa1, 0xbd, 0xf9, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x93, 0x0d, 0xf7, 0x22, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OncallInfoClient is the client API for OncallInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OncallInfoClient interface {
	// Oncall can be used to get the current oncaller(s) for the
	// provided rotation.
	Oncall(ctx context.Context, in *OncallRequest, opts ...grpc.CallOption) (*OncallResponse, error)
	// ListRotations lists all configured rotations.
	ListRotations(ctx context.Context, in *ListRotationsRequest, opts ...grpc.CallOption) (*ListRotationsResponse, error)
	// Shifts returns ShiftEntries for a specific time range.
	// Leaving times empty returns all shifts.
	Shifts(ctx context.Context, in *ShiftsRequest, opts ...grpc.CallOption) (*ShiftsResponse, error)
}
type oncallInfoPRPCClient struct {
	client *prpc.Client
}

func NewOncallInfoPRPCClient(client *prpc.Client) OncallInfoClient {
	return &oncallInfoPRPCClient{client}
}

func (c *oncallInfoPRPCClient) Oncall(ctx context.Context, in *OncallRequest, opts ...grpc.CallOption) (*OncallResponse, error) {
	out := new(OncallResponse)
	err := c.client.Call(ctx, "rotangapi.OncallInfo", "Oncall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oncallInfoPRPCClient) ListRotations(ctx context.Context, in *ListRotationsRequest, opts ...grpc.CallOption) (*ListRotationsResponse, error) {
	out := new(ListRotationsResponse)
	err := c.client.Call(ctx, "rotangapi.OncallInfo", "ListRotations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oncallInfoPRPCClient) Shifts(ctx context.Context, in *ShiftsRequest, opts ...grpc.CallOption) (*ShiftsResponse, error) {
	out := new(ShiftsResponse)
	err := c.client.Call(ctx, "rotangapi.OncallInfo", "Shifts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type oncallInfoClient struct {
	cc *grpc.ClientConn
}

func NewOncallInfoClient(cc *grpc.ClientConn) OncallInfoClient {
	return &oncallInfoClient{cc}
}

func (c *oncallInfoClient) Oncall(ctx context.Context, in *OncallRequest, opts ...grpc.CallOption) (*OncallResponse, error) {
	out := new(OncallResponse)
	err := c.cc.Invoke(ctx, "/rotangapi.OncallInfo/Oncall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oncallInfoClient) ListRotations(ctx context.Context, in *ListRotationsRequest, opts ...grpc.CallOption) (*ListRotationsResponse, error) {
	out := new(ListRotationsResponse)
	err := c.cc.Invoke(ctx, "/rotangapi.OncallInfo/ListRotations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oncallInfoClient) Shifts(ctx context.Context, in *ShiftsRequest, opts ...grpc.CallOption) (*ShiftsResponse, error) {
	out := new(ShiftsResponse)
	err := c.cc.Invoke(ctx, "/rotangapi.OncallInfo/Shifts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OncallInfoServer is the server API for OncallInfo service.
type OncallInfoServer interface {
	// Oncall can be used to get the current oncaller(s) for the
	// provided rotation.
	Oncall(context.Context, *OncallRequest) (*OncallResponse, error)
	// ListRotations lists all configured rotations.
	ListRotations(context.Context, *ListRotationsRequest) (*ListRotationsResponse, error)
	// Shifts returns ShiftEntries for a specific time range.
	// Leaving times empty returns all shifts.
	Shifts(context.Context, *ShiftsRequest) (*ShiftsResponse, error)
}

// UnimplementedOncallInfoServer can be embedded to have forward compatible implementations.
type UnimplementedOncallInfoServer struct {
}

func (*UnimplementedOncallInfoServer) Oncall(ctx context.Context, req *OncallRequest) (*OncallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Oncall not implemented")
}
func (*UnimplementedOncallInfoServer) ListRotations(ctx context.Context, req *ListRotationsRequest) (*ListRotationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRotations not implemented")
}
func (*UnimplementedOncallInfoServer) Shifts(ctx context.Context, req *ShiftsRequest) (*ShiftsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shifts not implemented")
}

func RegisterOncallInfoServer(s prpc.Registrar, srv OncallInfoServer) {
	s.RegisterService(&_OncallInfo_serviceDesc, srv)
}

func _OncallInfo_Oncall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OncallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OncallInfoServer).Oncall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rotangapi.OncallInfo/Oncall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OncallInfoServer).Oncall(ctx, req.(*OncallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OncallInfo_ListRotations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRotationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OncallInfoServer).ListRotations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rotangapi.OncallInfo/ListRotations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OncallInfoServer).ListRotations(ctx, req.(*ListRotationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OncallInfo_Shifts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShiftsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OncallInfoServer).Shifts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rotangapi.OncallInfo/Shifts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OncallInfoServer).Shifts(ctx, req.(*ShiftsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _OncallInfo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rotangapi.OncallInfo",
	HandlerType: (*OncallInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Oncall",
			Handler:    _OncallInfo_Oncall_Handler,
		},
		{
			MethodName: "ListRotations",
			Handler:    _OncallInfo_ListRotations_Handler,
		},
		{
			MethodName: "Shifts",
			Handler:    _OncallInfo_Shifts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infra/appengine/rotang/proto/rotangapi/oncallinfo.proto",
}
