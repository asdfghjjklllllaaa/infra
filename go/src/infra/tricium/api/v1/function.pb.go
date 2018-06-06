// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/tricium/api/v1/function.proto

package tricium

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Tricium functions; isolators and analyzers.
type Function_Type int32

const (
	Function_NONE     Function_Type = 0
	Function_ISOLATOR Function_Type = 1
	Function_ANALYZER Function_Type = 2
)

var Function_Type_name = map[int32]string{
	0: "NONE",
	1: "ISOLATOR",
	2: "ANALYZER",
}
var Function_Type_value = map[string]int32{
	"NONE":     0,
	"ISOLATOR": 1,
	"ANALYZER": 2,
}

func (x Function_Type) String() string {
	return proto.EnumName(Function_Type_name, int32(x))
}
func (Function_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0, 0} }

// Tricium function.
//
// There are two types of functions; isolators and analyzers.
// All functions have one input (needs) and one output (provides).
// For analyzer functions, the output must be of type Data.Results.
type Function struct {
	// The type of this function.
	//
	// This field is required.
	Type Function_Type `protobuf:"varint,1,opt,name=type,enum=tricium.Function_Type" json:"type,omitempty"`
	// The name of the function.
	//
	// This name is used for selection, customization and reporting of
	// progress/results. The name must be unique among Tricium function
	// within a Tricium instance.
	// This field is required.
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	// Data needed by this function.
	//
	// This field is required.
	Needs Data_Type `protobuf:"varint,3,opt,name=needs,enum=tricium.Data_Type" json:"needs,omitempty"`
	// Data provided by this function.
	//
	// This field is required.
	Provides Data_Type `protobuf:"varint,4,opt,name=provides,enum=tricium.Data_Type" json:"provides,omitempty"`
	// Path filters for this function.
	//
	// Applicable when this function is an analyzer. Defined as a glob.
	// The path filters only apply to the last part of the path.
	PathFilters []string `protobuf:"bytes,5,rep,name=path_filters,json=pathFilters" json:"path_filters,omitempty"`
	// Email address of the owner of this function.
	//
	// This field is required.
	Owner string `protobuf:"bytes,6,opt,name=owner" json:"owner,omitempty"`
	// Monorail bug component for bug filing.
	//
	// This field is required.
	MonorailComponent string `protobuf:"bytes,7,opt,name=monorail_component,json=monorailComponent" json:"monorail_component,omitempty"`
	// Function configuration options.
	//
	// These options enable projects to configure how a function behaves without
	// customization via a new implementation. For instance, an analyzer function
	// may expose the list of possible checks as configuration options.
	ConfigDefs []*ConfigDef `protobuf:"bytes,8,rep,name=config_defs,json=configDefs" json:"config_defs,omitempty"`
	// Function implementations.
	//
	// A function may run on many platforms and the implementation per platform
	// may differ. When possible, an implementation may be shared between several
	// platforms.
	Impls []*Impl `protobuf:"bytes,9,rep,name=impls" json:"impls,omitempty"`
}

func (m *Function) Reset()                    { *m = Function{} }
func (m *Function) String() string            { return proto.CompactTextString(m) }
func (*Function) ProtoMessage()               {}
func (*Function) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *Function) GetType() Function_Type {
	if m != nil {
		return m.Type
	}
	return Function_NONE
}

func (m *Function) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Function) GetNeeds() Data_Type {
	if m != nil {
		return m.Needs
	}
	return Data_NONE
}

func (m *Function) GetProvides() Data_Type {
	if m != nil {
		return m.Provides
	}
	return Data_NONE
}

func (m *Function) GetPathFilters() []string {
	if m != nil {
		return m.PathFilters
	}
	return nil
}

func (m *Function) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *Function) GetMonorailComponent() string {
	if m != nil {
		return m.MonorailComponent
	}
	return ""
}

func (m *Function) GetConfigDefs() []*ConfigDef {
	if m != nil {
		return m.ConfigDefs
	}
	return nil
}

func (m *Function) GetImpls() []*Impl {
	if m != nil {
		return m.Impls
	}
	return nil
}

// Definition of a function configuration.
//
// An analyzer may expose flags as configuration options, e.g. ClangTidy
// is configured with a 'checks' flag.
type ConfigDef struct {
	// Name of configuration option.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Default value for the config, e.g., checks="all".
	Default string `protobuf:"bytes,2,opt,name=default" json:"default,omitempty"`
}

func (m *ConfigDef) Reset()                    { *m = ConfigDef{} }
func (m *ConfigDef) String() string            { return proto.CompactTextString(m) }
func (*ConfigDef) ProtoMessage()               {}
func (*ConfigDef) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *ConfigDef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ConfigDef) GetDefault() string {
	if m != nil {
		return m.Default
	}
	return ""
}

// Function implementation.
//
// Implementation can be either command-based or recipe-based.
// NB! Recipe-based implementations are not supported yet.
//
// If platform-specific data is needed or provided, the specific platform
// details should be provided in the implementation. Note that the runtime
// platform of the implementation may be different then the platform(s)
// used to refine the data-dependency.
type Impl struct {
	// Data-dependency details specific to this implementation.
	//
	// For instance, if the needed data must be parameterized with a
	// specific platform then the 'needs_for_platform' field should be set to
	// that platform. Likewise for any provided data type that must be
	// parameterized with a specific platform, this should be indicated with
	// the 'provides_for_platform' field. Either if these fields can be left
	// out for implementations of functions not needing or providing
	// platform-specific data.
	NeedsForPlatform    Platform_Name `protobuf:"varint,1,opt,name=needs_for_platform,json=needsForPlatform,enum=tricium.Platform_Name" json:"needs_for_platform,omitempty"`
	ProvidesForPlatform Platform_Name `protobuf:"varint,2,opt,name=provides_for_platform,json=providesForPlatform,enum=tricium.Platform_Name" json:"provides_for_platform,omitempty"`
	// The platform to run this implementation on.
	//
	// This may be different from the platforms used to refine data-dependencies,
	// as long as the data consumed/produced follows the specification.
	RuntimePlatform Platform_Name `protobuf:"varint,3,opt,name=runtime_platform,json=runtimePlatform,enum=tricium.Platform_Name" json:"runtime_platform,omitempty"`
	// CIPD packages needed by this implementation.
	CipdPackages []*CipdPackage `protobuf:"bytes,4,rep,name=cipd_packages,json=cipdPackages" json:"cipd_packages,omitempty"`
	// Types that are valid to be assigned to Impl:
	//	*Impl_Recipe
	//	*Impl_Cmd
	Impl isImpl_Impl `protobuf_oneof:"impl"`
	// Deadline for execution of corresponding workers.
	//
	// Note that this deadline includes the launch of a swarming task for the
	// corresponding worker, and collection of results from that worker.
	// Deadline should be given in seconds.
	Deadline int32 `protobuf:"varint,7,opt,name=deadline" json:"deadline,omitempty"`
}

func (m *Impl) Reset()                    { *m = Impl{} }
func (m *Impl) String() string            { return proto.CompactTextString(m) }
func (*Impl) ProtoMessage()               {}
func (*Impl) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

type isImpl_Impl interface {
	isImpl_Impl()
}

type Impl_Recipe struct {
	Recipe *Recipe `protobuf:"bytes,5,opt,name=recipe,oneof"`
}
type Impl_Cmd struct {
	Cmd *Cmd `protobuf:"bytes,6,opt,name=cmd,oneof"`
}

func (*Impl_Recipe) isImpl_Impl() {}
func (*Impl_Cmd) isImpl_Impl()    {}

func (m *Impl) GetImpl() isImpl_Impl {
	if m != nil {
		return m.Impl
	}
	return nil
}

func (m *Impl) GetNeedsForPlatform() Platform_Name {
	if m != nil {
		return m.NeedsForPlatform
	}
	return Platform_ANY
}

func (m *Impl) GetProvidesForPlatform() Platform_Name {
	if m != nil {
		return m.ProvidesForPlatform
	}
	return Platform_ANY
}

func (m *Impl) GetRuntimePlatform() Platform_Name {
	if m != nil {
		return m.RuntimePlatform
	}
	return Platform_ANY
}

func (m *Impl) GetCipdPackages() []*CipdPackage {
	if m != nil {
		return m.CipdPackages
	}
	return nil
}

func (m *Impl) GetRecipe() *Recipe {
	if x, ok := m.GetImpl().(*Impl_Recipe); ok {
		return x.Recipe
	}
	return nil
}

func (m *Impl) GetCmd() *Cmd {
	if x, ok := m.GetImpl().(*Impl_Cmd); ok {
		return x.Cmd
	}
	return nil
}

func (m *Impl) GetDeadline() int32 {
	if m != nil {
		return m.Deadline
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Impl) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Impl_OneofMarshaler, _Impl_OneofUnmarshaler, _Impl_OneofSizer, []interface{}{
		(*Impl_Recipe)(nil),
		(*Impl_Cmd)(nil),
	}
}

func _Impl_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Impl)
	// impl
	switch x := m.Impl.(type) {
	case *Impl_Recipe:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Recipe); err != nil {
			return err
		}
	case *Impl_Cmd:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Cmd); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Impl.Impl has unexpected type %T", x)
	}
	return nil
}

func _Impl_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Impl)
	switch tag {
	case 5: // impl.recipe
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Recipe)
		err := b.DecodeMessage(msg)
		m.Impl = &Impl_Recipe{msg}
		return true, err
	case 6: // impl.cmd
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Cmd)
		err := b.DecodeMessage(msg)
		m.Impl = &Impl_Cmd{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Impl_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Impl)
	// impl
	switch x := m.Impl.(type) {
	case *Impl_Recipe:
		s := proto.Size(x.Recipe)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Impl_Cmd:
		s := proto.Size(x.Cmd)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Specification of a recipe for a recipe-based analyzer.
type Recipe struct {
	// Recipe CIPD package.
	CipdPackage string `protobuf:"bytes,1,opt,name=cipd_package,json=cipdPackage" json:"cipd_package,omitempty"`
	// CIPD package version.
	CipdVersion string `protobuf:"bytes,2,opt,name=cipd_version,json=cipdVersion" json:"cipd_version,omitempty"`
	// Extra recipe properties to add, as a JSON mapping of keys to values.
	Properties string `protobuf:"bytes,3,opt,name=properties" json:"properties,omitempty"`
}

func (m *Recipe) Reset()                    { *m = Recipe{} }
func (m *Recipe) String() string            { return proto.CompactTextString(m) }
func (*Recipe) ProtoMessage()               {}
func (*Recipe) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *Recipe) GetCipdPackage() string {
	if m != nil {
		return m.CipdPackage
	}
	return ""
}

func (m *Recipe) GetCipdVersion() string {
	if m != nil {
		return m.CipdVersion
	}
	return ""
}

func (m *Recipe) GetProperties() string {
	if m != nil {
		return m.Properties
	}
	return ""
}

// Specification of a command.
type Cmd struct {
	// Executable binary.
	Exec string `protobuf:"bytes,1,opt,name=exec" json:"exec,omitempty"`
	// Arguments in order.
	Args []string `protobuf:"bytes,2,rep,name=args" json:"args,omitempty"`
}

func (m *Cmd) Reset()                    { *m = Cmd{} }
func (m *Cmd) String() string            { return proto.CompactTextString(m) }
func (*Cmd) ProtoMessage()               {}
func (*Cmd) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *Cmd) GetExec() string {
	if m != nil {
		return m.Exec
	}
	return ""
}

func (m *Cmd) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

// Specification of a CIPD package that is installed as a dependency of an
// analyzer.
type CipdPackage struct {
	// CIPD package name.
	PackageName string `protobuf:"bytes,1,opt,name=package_name,json=packageName" json:"package_name,omitempty"`
	// Relative path from the working directory where the package shall be
	// installed. Cannot be empty or start with a slash.
	Path string `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	// CIPD package version.
	Version string `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
}

func (m *CipdPackage) Reset()                    { *m = CipdPackage{} }
func (m *CipdPackage) String() string            { return proto.CompactTextString(m) }
func (*CipdPackage) ProtoMessage()               {}
func (*CipdPackage) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

func (m *CipdPackage) GetPackageName() string {
	if m != nil {
		return m.PackageName
	}
	return ""
}

func (m *CipdPackage) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *CipdPackage) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func init() {
	proto.RegisterType((*Function)(nil), "tricium.Function")
	proto.RegisterType((*ConfigDef)(nil), "tricium.ConfigDef")
	proto.RegisterType((*Impl)(nil), "tricium.Impl")
	proto.RegisterType((*Recipe)(nil), "tricium.Recipe")
	proto.RegisterType((*Cmd)(nil), "tricium.Cmd")
	proto.RegisterType((*CipdPackage)(nil), "tricium.CipdPackage")
	proto.RegisterEnum("tricium.Function_Type", Function_Type_name, Function_Type_value)
}

func init() { proto.RegisterFile("infra/tricium/api/v1/function.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 623 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0x5f, 0x6f, 0xd3, 0x3e,
	0x14, 0x5d, 0x9b, 0xb4, 0x6b, 0x6f, 0xba, 0xdf, 0xfa, 0x33, 0x03, 0x45, 0x7b, 0x80, 0x92, 0xbd,
	0x04, 0xc4, 0x5a, 0xb1, 0x3d, 0xed, 0xb1, 0x74, 0x9b, 0x36, 0x34, 0x75, 0x93, 0x99, 0x90, 0xe0,
	0x81, 0xc8, 0x24, 0xce, 0xb0, 0x48, 0x6c, 0xcb, 0x49, 0x07, 0xfb, 0x04, 0x7c, 0x01, 0x3e, 0x30,
	0xb2, 0x13, 0xa7, 0x99, 0x18, 0x7b, 0xbb, 0x7f, 0xce, 0x3d, 0xbe, 0xf1, 0x39, 0x0e, 0xec, 0x31,
	0x9e, 0x2a, 0x32, 0x2b, 0x15, 0x8b, 0xd9, 0x2a, 0x9f, 0x11, 0xc9, 0x66, 0xb7, 0x6f, 0x67, 0xe9,
	0x8a, 0xc7, 0x25, 0x13, 0x7c, 0x2a, 0x95, 0x28, 0x05, 0xda, 0xac, 0xdb, 0xbb, 0x2f, 0x1e, 0x44,
	0x27, 0xa4, 0x24, 0x15, 0x72, 0xf7, 0x61, 0x3a, 0x99, 0x91, 0x32, 0x15, 0x2a, 0xaf, 0x40, 0xc1,
	0x6f, 0x07, 0x06, 0xa7, 0xf5, 0x09, 0xe8, 0x35, 0xb8, 0xe5, 0x9d, 0xa4, 0x7e, 0x67, 0xd2, 0x09,
	0xff, 0x3b, 0x78, 0x36, 0xad, 0x47, 0xa7, 0x16, 0x30, 0xbd, 0xbe, 0x93, 0x14, 0x1b, 0x0c, 0x42,
	0xe0, 0x72, 0x92, 0x53, 0xbf, 0x3b, 0xe9, 0x84, 0x43, 0x6c, 0x62, 0x14, 0x42, 0x8f, 0x53, 0x9a,
	0x14, 0xbe, 0x63, 0x08, 0x50, 0x43, 0x70, 0xac, 0xb7, 0x32, 0xc3, 0x15, 0x00, 0x4d, 0x61, 0x20,
	0x95, 0xb8, 0x65, 0x09, 0x2d, 0x7c, 0xf7, 0x9f, 0xe0, 0x06, 0x83, 0x5e, 0xc2, 0x48, 0x92, 0xf2,
	0x5b, 0x94, 0xb2, 0xac, 0xa4, 0xaa, 0xf0, 0x7b, 0x13, 0x27, 0x1c, 0x62, 0x4f, 0xd7, 0x4e, 0xab,
	0x12, 0xda, 0x81, 0x9e, 0xf8, 0xc1, 0xa9, 0xf2, 0xfb, 0x66, 0xa3, 0x2a, 0x41, 0xfb, 0x80, 0x72,
	0xc1, 0x85, 0x22, 0x2c, 0x8b, 0x62, 0x91, 0x4b, 0xc1, 0x29, 0x2f, 0xfd, 0x4d, 0x03, 0xf9, 0xdf,
	0x76, 0x16, 0xb6, 0x81, 0x0e, 0xc1, 0x8b, 0x05, 0x4f, 0xd9, 0x4d, 0x94, 0xd0, 0xb4, 0xf0, 0x07,
	0x13, 0x27, 0xf4, 0x5a, 0xab, 0x2d, 0x4c, 0xef, 0x98, 0xa6, 0x18, 0x62, 0x1b, 0x16, 0x68, 0x0f,
	0x7a, 0x2c, 0x97, 0x59, 0xe1, 0x0f, 0x0d, 0x7c, 0xab, 0x81, 0x9f, 0xe7, 0x32, 0xc3, 0x55, 0x2f,
	0x78, 0x03, 0xae, 0xfe, 0x26, 0x34, 0x00, 0x77, 0x79, 0xb9, 0x3c, 0x19, 0x6f, 0xa0, 0x11, 0x0c,
	0xce, 0x3f, 0x5c, 0x5e, 0xcc, 0xaf, 0x2f, 0xf1, 0xb8, 0xa3, 0xb3, 0xf9, 0x72, 0x7e, 0xf1, 0xe9,
	0xf3, 0x09, 0x1e, 0x77, 0x83, 0x23, 0x18, 0x36, 0x67, 0x35, 0x57, 0xdd, 0x69, 0x5d, 0xb5, 0x0f,
	0x9b, 0x09, 0x4d, 0xc9, 0x2a, 0x2b, 0x6b, 0x05, 0x6c, 0x1a, 0xfc, 0x72, 0xc0, 0xd5, 0x07, 0xa3,
	0x63, 0x40, 0xe6, 0xb2, 0xa3, 0x54, 0xa8, 0xc8, 0xca, 0xfe, 0x97, 0xb6, 0x57, 0xd6, 0x0f, 0x4b,
	0x92, 0x53, 0x3c, 0x36, 0x13, 0xa7, 0x42, 0xd9, 0x32, 0x7a, 0x0f, 0x4f, 0xad, 0x0a, 0xf7, 0x89,
	0xba, 0x8f, 0x12, 0x3d, 0xb1, 0x43, 0x6d, 0xae, 0x39, 0x8c, 0xd5, 0x8a, 0x97, 0x2c, 0xa7, 0x6b,
	0x1a, 0xe7, 0x51, 0x9a, 0xed, 0x1a, 0xdf, 0x50, 0x1c, 0xc1, 0x56, 0xcc, 0x64, 0x12, 0x49, 0x12,
	0x7f, 0x27, 0x37, 0xc6, 0x3d, 0xfa, 0xce, 0x77, 0xd6, 0x12, 0x31, 0x99, 0x5c, 0x55, 0x4d, 0x3c,
	0x8a, 0xd7, 0x49, 0x81, 0x5e, 0x41, 0x5f, 0xd1, 0x98, 0x49, 0xea, 0xf7, 0x26, 0x9d, 0xd0, 0x3b,
	0xd8, 0x6e, 0x66, 0xb0, 0x29, 0x9f, 0x6d, 0xe0, 0x1a, 0x80, 0x26, 0xe0, 0xc4, 0x79, 0x62, 0x9c,
	0xe4, 0x1d, 0x8c, 0xd6, 0xdc, 0x79, 0x72, 0xb6, 0x81, 0x75, 0x0b, 0xed, 0xc2, 0x20, 0xa1, 0x24,
	0xc9, 0x18, 0xa7, 0xc6, 0x4d, 0x3d, 0xdc, 0xe4, 0xef, 0xfa, 0xe0, 0x6a, 0xcd, 0x03, 0x0e, 0xfd,
	0x8a, 0x59, 0xdb, 0xb7, 0xbd, 0x75, 0xad, 0xa4, 0xd7, 0x5a, 0xaf, 0x81, 0xdc, 0x52, 0x55, 0x30,
	0xc1, 0x6b, 0x55, 0x0d, 0xe4, 0x63, 0x55, 0x42, 0xcf, 0x01, 0xa4, 0x12, 0x92, 0xaa, 0x92, 0xd1,
	0xea, 0x8d, 0x0d, 0x71, 0xab, 0x12, 0xec, 0x83, 0xb3, 0xc8, 0x13, 0x6d, 0x17, 0xfa, 0x93, 0xc6,
	0xd6, 0x2e, 0x3a, 0xd6, 0x35, 0xa2, 0x6e, 0x0a, 0xbf, 0x6b, 0xde, 0x8d, 0x89, 0x83, 0x2f, 0xe0,
	0x2d, 0xee, 0x2f, 0x50, 0xaf, 0x17, 0xb5, 0xdc, 0xe6, 0xd5, 0x35, 0xad, 0x86, 0x66, 0xd1, 0x2f,
	0xce, 0xbe, 0x79, 0x1d, 0x6b, 0x23, 0xda, 0x95, 0xab, 0x8d, 0x6c, 0xfa, 0xb5, 0x6f, 0xfe, 0x30,
	0x87, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xd1, 0x3e, 0xba, 0x8a, 0xd7, 0x04, 0x00, 0x00,
}
