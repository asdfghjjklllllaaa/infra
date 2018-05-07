// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/tricium/api/v1/config.proto

/*
Package tricium is a generated protocol buffer package.

It is generated from these files:
	infra/tricium/api/v1/config.proto
	infra/tricium/api/v1/data.proto
	infra/tricium/api/v1/function.proto
	infra/tricium/api/v1/platform.proto
	infra/tricium/api/v1/tricium.proto

It has these top-level messages:
	ServiceConfig
	ProjectDetails
	ProjectConfig
	RepoDetails
	GitRepoDetails
	GerritDetails
	GerritProject
	GitRepo
	Acl
	Selection
	Config
	Data
	Function
	ConfigDef
	Impl
	Recipe
	Property
	Cmd
	CipdPackage
	Platform
	AnalyzeRequest
	GerritRevision
	GitCommit
	AnalyzeResponse
	ProgressRequest
	ProgressResponse
	FunctionProgress
	ProjectProgressRequest
	ProjectProgressResponse
	RunProgress
	ResultsRequest
	ResultsResponse
	FeedbackRequest
	FeedbackResponse
	ReportNotUsefulRequest
	ReportNotUsefulResponse
	GerritConsumerDetails
*/
package tricium

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

// Supported kinds of repositories.
// DEPRECATED, see https://crbug.com/824558
type RepoDetails_Kind int32

const (
	RepoDetails_GIT RepoDetails_Kind = 0
)

var RepoDetails_Kind_name = map[int32]string{
	0: "GIT",
}
var RepoDetails_Kind_value = map[string]int32{
	"GIT": 0,
}

func (x RepoDetails_Kind) String() string {
	return proto.EnumName(RepoDetails_Kind_name, int32(x))
}
func (RepoDetails_Kind) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

// Roles relevant to Tricium.
type Acl_Role int32

const (
	// Can read progress/results.
	Acl_READER Acl_Role = 0
	// Can request analysis.
	Acl_REQUESTER Acl_Role = 1
)

var Acl_Role_name = map[int32]string{
	0: "READER",
	1: "REQUESTER",
}
var Acl_Role_value = map[string]int32{
	"READER":    0,
	"REQUESTER": 1,
}

func (x Acl_Role) String() string {
	return proto.EnumName(Acl_Role_name, int32(x))
}
func (Acl_Role) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{8, 0} }

// Tricium service configuration.
//
// Listing supported platforms and analyzers shared between projects connected
// to Tricium.
type ServiceConfig struct {
	// Supported platforms.
	Platforms []*Platform_Details `protobuf:"bytes,1,rep,name=platforms" json:"platforms,omitempty"`
	// Supported data types.
	DataDetails []*Data_TypeDetails `protobuf:"bytes,2,rep,name=data_details,json=dataDetails" json:"data_details,omitempty"`
	// List of shared functions.
	Functions []*Function `protobuf:"bytes,3,rep,name=functions" json:"functions,omitempty"`
	// Details for connected projects.
	// DEPRECATED, see https://crbug.com/824558
	Projects []*ProjectDetails `protobuf:"bytes,4,rep,name=projects" json:"projects,omitempty"`
	// Base recipe command used for workers implemented as recipes.
	//
	// Specific recipe details for the worker will be added as flags at the
	// end of the argument list.
	RecipeCmd *Cmd `protobuf:"bytes,5,opt,name=recipe_cmd,json=recipeCmd" json:"recipe_cmd,omitempty"`
	// Base recipe packages used for workers implemented as recipes.
	//
	// These packages will be adjusted for the platform in question, by appending
	// platform name details to the end of the package name.
	RecipePackages []*CipdPackage `protobuf:"bytes,6,rep,name=recipe_packages,json=recipePackages" json:"recipe_packages,omitempty"`
	// Swarming server to use for this service instance.
	//
	// This should be a full URL with no trailing slash.
	SwarmingServer string `protobuf:"bytes,7,opt,name=swarming_server,json=swarmingServer" json:"swarming_server,omitempty"`
	// Isolate server to use for this service instance.
	//
	// This should be a full URL with no trailing slash.
	IsolateServer string `protobuf:"bytes,8,opt,name=isolate_server,json=isolateServer" json:"isolate_server,omitempty"`
}

func (m *ServiceConfig) Reset()                    { *m = ServiceConfig{} }
func (m *ServiceConfig) String() string            { return proto.CompactTextString(m) }
func (*ServiceConfig) ProtoMessage()               {}
func (*ServiceConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ServiceConfig) GetPlatforms() []*Platform_Details {
	if m != nil {
		return m.Platforms
	}
	return nil
}

func (m *ServiceConfig) GetDataDetails() []*Data_TypeDetails {
	if m != nil {
		return m.DataDetails
	}
	return nil
}

func (m *ServiceConfig) GetFunctions() []*Function {
	if m != nil {
		return m.Functions
	}
	return nil
}

func (m *ServiceConfig) GetProjects() []*ProjectDetails {
	if m != nil {
		return m.Projects
	}
	return nil
}

func (m *ServiceConfig) GetRecipeCmd() *Cmd {
	if m != nil {
		return m.RecipeCmd
	}
	return nil
}

func (m *ServiceConfig) GetRecipePackages() []*CipdPackage {
	if m != nil {
		return m.RecipePackages
	}
	return nil
}

func (m *ServiceConfig) GetSwarmingServer() string {
	if m != nil {
		return m.SwarmingServer
	}
	return ""
}

func (m *ServiceConfig) GetIsolateServer() string {
	if m != nil {
		return m.IsolateServer
	}
	return ""
}

// DEPRECATED, see https://crbug.com/824558
type ProjectDetails struct {
	// Project name used to map these project details to the config for a project.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// General service account for this project.
	// Used for any service interaction, with the exception of swarming.
	ServiceAccount string `protobuf:"bytes,2,opt,name=service_account,json=serviceAccount" json:"service_account,omitempty"`
	// Project-specific swarming service account.
	SwarmingServiceAccount string `protobuf:"bytes,3,opt,name=swarming_service_account,json=swarmingServiceAccount" json:"swarming_service_account,omitempty"`
	// Details of the repository connected to the project. This should be the
	// repository hosting the files that should be analyzed for this project.
	RepoDetails *RepoDetails `protobuf:"bytes,4,opt,name=repo_details,json=repoDetails" json:"repo_details,omitempty"`
	// Gerrit details for a project.
	//
	// This field should only be included if there is a Gerrit host for a
	// project and that host should be polled for changes and used for
	// reporting of analyzer progress and results.
	GerritDetails *GerritDetails `protobuf:"bytes,5,opt,name=gerrit_details,json=gerritDetails" json:"gerrit_details,omitempty"`
}

func (m *ProjectDetails) Reset()                    { *m = ProjectDetails{} }
func (m *ProjectDetails) String() string            { return proto.CompactTextString(m) }
func (*ProjectDetails) ProtoMessage()               {}
func (*ProjectDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ProjectDetails) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProjectDetails) GetServiceAccount() string {
	if m != nil {
		return m.ServiceAccount
	}
	return ""
}

func (m *ProjectDetails) GetSwarmingServiceAccount() string {
	if m != nil {
		return m.SwarmingServiceAccount
	}
	return ""
}

func (m *ProjectDetails) GetRepoDetails() *RepoDetails {
	if m != nil {
		return m.RepoDetails
	}
	return nil
}

func (m *ProjectDetails) GetGerritDetails() *GerritDetails {
	if m != nil {
		return m.GerritDetails
	}
	return nil
}

// Tricium project configuration.
//
// Specifies details needed to connect a project to Tricium.
// Adds project-specific functions and selects shared function
// implementations.
type ProjectConfig struct {
	// Project name,
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Access control rules for the project.
	Acls []*Acl `protobuf:"bytes,2,rep,name=acls" json:"acls,omitempty"`
	// Project-specific function details.
	//
	// This includes project-specific analyzer implementations and full
	// project-specific analyzer specifications.
	Functions []*Function `protobuf:"bytes,3,rep,name=functions" json:"functions,omitempty"`
	// Selection of function implementations to run for this project.
	Selections []*Selection `protobuf:"bytes,4,rep,name=selections" json:"selections,omitempty"`
	// Repositories, including Git and Gerrit details.
	Repos []*RepoDetails `protobuf:"bytes,5,rep,name=repos" json:"repos,omitempty"`
	// General service account for this project.
	// Used for any service interaction, with the exception of swarming.
	ServiceAccount string `protobuf:"bytes,6,opt,name=service_account,json=serviceAccount" json:"service_account,omitempty"`
	// Project-specific swarming service account.
	SwarmingServiceAccount string `protobuf:"bytes,7,opt,name=swarming_service_account,json=swarmingServiceAccount" json:"swarming_service_account,omitempty"`
}

func (m *ProjectConfig) Reset()                    { *m = ProjectConfig{} }
func (m *ProjectConfig) String() string            { return proto.CompactTextString(m) }
func (*ProjectConfig) ProtoMessage()               {}
func (*ProjectConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ProjectConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProjectConfig) GetAcls() []*Acl {
	if m != nil {
		return m.Acls
	}
	return nil
}

func (m *ProjectConfig) GetFunctions() []*Function {
	if m != nil {
		return m.Functions
	}
	return nil
}

func (m *ProjectConfig) GetSelections() []*Selection {
	if m != nil {
		return m.Selections
	}
	return nil
}

func (m *ProjectConfig) GetRepos() []*RepoDetails {
	if m != nil {
		return m.Repos
	}
	return nil
}

func (m *ProjectConfig) GetServiceAccount() string {
	if m != nil {
		return m.ServiceAccount
	}
	return ""
}

func (m *ProjectConfig) GetSwarmingServiceAccount() string {
	if m != nil {
		return m.SwarmingServiceAccount
	}
	return ""
}

// Repository details for a project.
// DEPRECATED, see https://crbug.com/824558
type RepoDetails struct {
	// DEPRECATED, see https://crbug.com/824558
	Kind RepoDetails_Kind `protobuf:"varint,1,opt,name=kind,enum=tricium.RepoDetails_Kind" json:"kind,omitempty"`
	// If repository kind is GIT then provide Git details.
	// DEPRECATED, see https://crbug.com/824558
	GitDetails *GitRepoDetails `protobuf:"bytes,2,opt,name=git_details,json=gitDetails" json:"git_details,omitempty"`
	// If there is an associated Gerrit project, details can be added here.
	// By adding GerritDetails, Gerrit polling is enabled.
	// DEPRECATED, see https://crbug.com/824558
	GerritDetails *GerritDetails `protobuf:"bytes,3,opt,name=gerrit_details,json=gerritDetails" json:"gerrit_details,omitempty"`
	// Could be renamed to kind when the above kind is removed.
	//
	// Types that are valid to be assigned to Source:
	//	*RepoDetails_GerritProject
	//	*RepoDetails_GitRepo
	Source isRepoDetails_Source `protobuf_oneof:"source"`
}

func (m *RepoDetails) Reset()                    { *m = RepoDetails{} }
func (m *RepoDetails) String() string            { return proto.CompactTextString(m) }
func (*RepoDetails) ProtoMessage()               {}
func (*RepoDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type isRepoDetails_Source interface {
	isRepoDetails_Source()
}

type RepoDetails_GerritProject struct {
	GerritProject *GerritProject `protobuf:"bytes,4,opt,name=gerrit_project,json=gerritProject,oneof"`
}
type RepoDetails_GitRepo struct {
	GitRepo *GitRepo `protobuf:"bytes,5,opt,name=git_repo,json=gitRepo,oneof"`
}

func (*RepoDetails_GerritProject) isRepoDetails_Source() {}
func (*RepoDetails_GitRepo) isRepoDetails_Source()       {}

func (m *RepoDetails) GetSource() isRepoDetails_Source {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *RepoDetails) GetKind() RepoDetails_Kind {
	if m != nil {
		return m.Kind
	}
	return RepoDetails_GIT
}

func (m *RepoDetails) GetGitDetails() *GitRepoDetails {
	if m != nil {
		return m.GitDetails
	}
	return nil
}

func (m *RepoDetails) GetGerritDetails() *GerritDetails {
	if m != nil {
		return m.GerritDetails
	}
	return nil
}

func (m *RepoDetails) GetGerritProject() *GerritProject {
	if x, ok := m.GetSource().(*RepoDetails_GerritProject); ok {
		return x.GerritProject
	}
	return nil
}

func (m *RepoDetails) GetGitRepo() *GitRepo {
	if x, ok := m.GetSource().(*RepoDetails_GitRepo); ok {
		return x.GitRepo
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*RepoDetails) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _RepoDetails_OneofMarshaler, _RepoDetails_OneofUnmarshaler, _RepoDetails_OneofSizer, []interface{}{
		(*RepoDetails_GerritProject)(nil),
		(*RepoDetails_GitRepo)(nil),
	}
}

func _RepoDetails_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*RepoDetails)
	// source
	switch x := m.Source.(type) {
	case *RepoDetails_GerritProject:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.GerritProject); err != nil {
			return err
		}
	case *RepoDetails_GitRepo:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.GitRepo); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("RepoDetails.Source has unexpected type %T", x)
	}
	return nil
}

func _RepoDetails_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*RepoDetails)
	switch tag {
	case 4: // source.gerrit_project
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(GerritProject)
		err := b.DecodeMessage(msg)
		m.Source = &RepoDetails_GerritProject{msg}
		return true, err
	case 5: // source.git_repo
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(GitRepo)
		err := b.DecodeMessage(msg)
		m.Source = &RepoDetails_GitRepo{msg}
		return true, err
	default:
		return false, nil
	}
}

func _RepoDetails_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*RepoDetails)
	// source
	switch x := m.Source.(type) {
	case *RepoDetails_GerritProject:
		s := proto.Size(x.GerritProject)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *RepoDetails_GitRepo:
		s := proto.Size(x.GitRepo)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Git repository details.
// DEPRECATED, see https://crbug.com/824558
type GitRepoDetails struct {
	// URL to repository.
	Repository string `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	// Default ref to use to get files to analyze.
	Ref string `protobuf:"bytes,2,opt,name=ref" json:"ref,omitempty"`
}

func (m *GitRepoDetails) Reset()                    { *m = GitRepoDetails{} }
func (m *GitRepoDetails) String() string            { return proto.CompactTextString(m) }
func (*GitRepoDetails) ProtoMessage()               {}
func (*GitRepoDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GitRepoDetails) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func (m *GitRepoDetails) GetRef() string {
	if m != nil {
		return m.Ref
	}
	return ""
}

// Gerrit details for a project.
// DEPRECATED, see https://crbug.com/824558
type GerritDetails struct {
	// The Gerrit host to connect to.
	//
	// Value must not include protocol.
	Host string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	// Gerrit project name.
	Project string `protobuf:"bytes,2,opt,name=project" json:"project,omitempty"`
	// Disable reporting to Gerrit.
	//
	// Whether to send progress and results to Gerrit.
	DisableReporting bool `protobuf:"varint,3,opt,name=disable_reporting,json=disableReporting" json:"disable_reporting,omitempty"`
	// Whitelisted groups.
	//
	// The owner of a change will be checked for membership of a whitelisted
	// group. An asterisk entry means all groups are whitelisted.
	//
	// Group names must be known to the Chrome infra auth service,
	// https://chrome-infra-auth.appspot.com. Contact a Chromium trooper
	// if you need to add or modify a group: g.co/bugatrooper.
	//
	// Note that no presence of this field effectively blacklists all groups.
	WhitelistedGroup []string `protobuf:"bytes,4,rep,name=whitelisted_group,json=whitelistedGroup" json:"whitelisted_group,omitempty"`
}

func (m *GerritDetails) Reset()                    { *m = GerritDetails{} }
func (m *GerritDetails) String() string            { return proto.CompactTextString(m) }
func (*GerritDetails) ProtoMessage()               {}
func (*GerritDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GerritDetails) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *GerritDetails) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *GerritDetails) GetDisableReporting() bool {
	if m != nil {
		return m.DisableReporting
	}
	return false
}

func (m *GerritDetails) GetWhitelistedGroup() []string {
	if m != nil {
		return m.WhitelistedGroup
	}
	return nil
}

// Specifies a Gerrit project and its corresponding git repo.
type GerritProject struct {
	// The Gerrit host to connect to.
	//
	// Value must not include the schema part; it will be assumed to be "https".
	Host string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	// Gerrit project name.
	Project string `protobuf:"bytes,2,opt,name=project" json:"project,omitempty"`
	// Full URL for the corresponding git repo.
	GitUrl string `protobuf:"bytes,3,opt,name=git_url,json=gitUrl" json:"git_url,omitempty"`
}

func (m *GerritProject) Reset()                    { *m = GerritProject{} }
func (m *GerritProject) String() string            { return proto.CompactTextString(m) }
func (*GerritProject) ProtoMessage()               {}
func (*GerritProject) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *GerritProject) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *GerritProject) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *GerritProject) GetGitUrl() string {
	if m != nil {
		return m.GitUrl
	}
	return ""
}

type GitRepo struct {
	// Full repository url, including schema.
	Url string `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
}

func (m *GitRepo) Reset()                    { *m = GitRepo{} }
func (m *GitRepo) String() string            { return proto.CompactTextString(m) }
func (*GitRepo) ProtoMessage()               {}
func (*GitRepo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *GitRepo) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

// Access control rules.
type Acl struct {
	// Role of a group or identity.
	Role Acl_Role `protobuf:"varint,1,opt,name=role,enum=tricium.Acl_Role" json:"role,omitempty"`
	// Name of group, as defined in the auth service. Specify either group or
	// identity, not both.
	Group string `protobuf:"bytes,2,opt,name=group" json:"group,omitempty"`
	// Identity, as defined by the auth service. Can be either an email address
	// or an identity string, for instance, "anonymous:anonymous" for anonymous
	// users. Specify either group or identity, not both.
	Identity string `protobuf:"bytes,3,opt,name=identity" json:"identity,omitempty"`
}

func (m *Acl) Reset()                    { *m = Acl{} }
func (m *Acl) String() string            { return proto.CompactTextString(m) }
func (*Acl) ProtoMessage()               {}
func (*Acl) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Acl) GetRole() Acl_Role {
	if m != nil {
		return m.Role
	}
	return Acl_READER
}

func (m *Acl) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *Acl) GetIdentity() string {
	if m != nil {
		return m.Identity
	}
	return ""
}

// Selection of function implementations to run for a project.
type Selection struct {
	// Name of function to run.
	Function string `protobuf:"bytes,1,opt,name=function" json:"function,omitempty"`
	// Name of platform to retrieve results from.
	Platform Platform_Name `protobuf:"varint,2,opt,name=platform,enum=tricium.Platform_Name" json:"platform,omitempty"`
	// Function configuration to use on this platform.
	Configs []*Config `protobuf:"bytes,3,rep,name=configs" json:"configs,omitempty"`
}

func (m *Selection) Reset()                    { *m = Selection{} }
func (m *Selection) String() string            { return proto.CompactTextString(m) }
func (*Selection) ProtoMessage()               {}
func (*Selection) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *Selection) GetFunction() string {
	if m != nil {
		return m.Function
	}
	return ""
}

func (m *Selection) GetPlatform() Platform_Name {
	if m != nil {
		return m.Platform
	}
	return Platform_ANY
}

func (m *Selection) GetConfigs() []*Config {
	if m != nil {
		return m.Configs
	}
	return nil
}

// Function configuration used when selecting a function implementation.
type Config struct {
	// Name of the configuration option.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Value of the configuration.
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Config) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Config) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*ServiceConfig)(nil), "tricium.ServiceConfig")
	proto.RegisterType((*ProjectDetails)(nil), "tricium.ProjectDetails")
	proto.RegisterType((*ProjectConfig)(nil), "tricium.ProjectConfig")
	proto.RegisterType((*RepoDetails)(nil), "tricium.RepoDetails")
	proto.RegisterType((*GitRepoDetails)(nil), "tricium.GitRepoDetails")
	proto.RegisterType((*GerritDetails)(nil), "tricium.GerritDetails")
	proto.RegisterType((*GerritProject)(nil), "tricium.GerritProject")
	proto.RegisterType((*GitRepo)(nil), "tricium.GitRepo")
	proto.RegisterType((*Acl)(nil), "tricium.Acl")
	proto.RegisterType((*Selection)(nil), "tricium.Selection")
	proto.RegisterType((*Config)(nil), "tricium.Config")
	proto.RegisterEnum("tricium.RepoDetails_Kind", RepoDetails_Kind_name, RepoDetails_Kind_value)
	proto.RegisterEnum("tricium.Acl_Role", Acl_Role_name, Acl_Role_value)
}

func init() { proto.RegisterFile("infra/tricium/api/v1/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 875 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xdd, 0x6e, 0xdb, 0x36,
	0x14, 0x8e, 0x2c, 0xc7, 0x3f, 0xc7, 0xb1, 0xe3, 0x10, 0x41, 0xab, 0x65, 0xc0, 0x96, 0x6a, 0x28,
	0x96, 0xad, 0xa8, 0x8d, 0xb9, 0x17, 0xed, 0xc5, 0x8a, 0x21, 0x4d, 0xbc, 0x74, 0x18, 0x30, 0x74,
	0x4c, 0xba, 0x5b, 0x83, 0x95, 0x68, 0x95, 0xab, 0x2c, 0x0a, 0x14, 0x9d, 0x22, 0x97, 0xbb, 0xd9,
	0x2b, 0xec, 0x15, 0xf6, 0x0a, 0xbb, 0xdb, 0x9b, 0x6d, 0xe0, 0xaf, 0x94, 0xc0, 0x0d, 0xe0, 0x3b,
	0x9e, 0xf3, 0x7d, 0x1f, 0x79, 0xfe, 0x48, 0xc2, 0x23, 0x56, 0x2c, 0x05, 0x99, 0x4a, 0xc1, 0x12,
	0xb6, 0x5e, 0x4d, 0x49, 0xc9, 0xa6, 0xd7, 0xdf, 0x4d, 0x13, 0x5e, 0x2c, 0x59, 0x36, 0x29, 0x05,
	0x97, 0x1c, 0x75, 0x2d, 0x78, 0xf4, 0xe5, 0x46, 0x6e, 0x4a, 0x24, 0x31, 0xcc, 0xa3, 0xaf, 0x36,
	0x12, 0x96, 0xeb, 0x22, 0x91, 0x8c, 0x17, 0xf7, 0x92, 0xca, 0x9c, 0xc8, 0x25, 0x17, 0x2b, 0x43,
	0x8a, 0xff, 0x0e, 0x61, 0x78, 0x49, 0xc5, 0x35, 0x4b, 0xe8, 0x99, 0x8e, 0x05, 0x3d, 0x87, 0xbe,
	0xe3, 0x54, 0x51, 0x70, 0x1c, 0x9e, 0x0c, 0x66, 0x9f, 0x4d, 0xec, 0x26, 0x93, 0x37, 0x4e, 0x7d,
	0x4e, 0x25, 0x61, 0x79, 0x85, 0x6b, 0x2e, 0xfa, 0x1e, 0xf6, 0x54, 0x88, 0x8b, 0xd4, 0x40, 0x51,
	0xeb, 0x8e, 0xf6, 0x5c, 0xc5, 0x7f, 0x75, 0x53, 0x52, 0xa7, 0x1d, 0x28, 0xba, 0x35, 0xd0, 0x14,
	0xfa, 0x2e, 0xfe, 0x2a, 0x0a, 0xb5, 0xf4, 0xc0, 0x4b, 0x7f, 0xb4, 0x08, 0xae, 0x39, 0xe8, 0x19,
	0xf4, 0x4a, 0xc1, 0x7f, 0xa7, 0x89, 0xac, 0xa2, 0xb6, 0xe6, 0x3f, 0xac, 0xc3, 0x34, 0x80, 0x3b,
	0xc8, 0x13, 0xd1, 0x13, 0x00, 0x41, 0x13, 0x56, 0xd2, 0x45, 0xb2, 0x4a, 0xa3, 0xdd, 0xe3, 0xe0,
	0x64, 0x30, 0xdb, 0xf3, 0xb2, 0xb3, 0x55, 0x8a, 0xfb, 0x06, 0x3f, 0x5b, 0xa5, 0xe8, 0x25, 0xec,
	0x5b, 0x72, 0x49, 0x92, 0x0f, 0x24, 0xa3, 0x55, 0xd4, 0xd1, 0x07, 0x1d, 0xd6, 0x0a, 0x56, 0xa6,
	0x6f, 0x0c, 0x88, 0x47, 0x86, 0x6c, 0xcd, 0x0a, 0x7d, 0x0d, 0xfb, 0xd5, 0x47, 0x22, 0x56, 0xac,
	0xc8, 0x16, 0x15, 0x15, 0xd7, 0x54, 0x44, 0xdd, 0xe3, 0xe0, 0xa4, 0x8f, 0x47, 0xce, 0x7d, 0xa9,
	0xbd, 0xe8, 0x31, 0x8c, 0x58, 0xc5, 0x73, 0x22, 0xa9, 0xe3, 0xf5, 0x34, 0x6f, 0x68, 0xbd, 0x86,
	0x16, 0xff, 0x17, 0xc0, 0xe8, 0x76, 0x62, 0x08, 0x41, 0xbb, 0x20, 0x2b, 0x1a, 0x05, 0x9a, 0xaf,
	0xd7, 0xfa, 0x58, 0xd3, 0xd0, 0x05, 0x49, 0x12, 0xbe, 0x2e, 0x64, 0xd4, 0xb2, 0xc7, 0x1a, 0xf7,
	0xa9, 0xf1, 0xa2, 0x17, 0x10, 0xdd, 0x8a, 0xaf, 0xa9, 0x08, 0xb5, 0xe2, 0x41, 0x33, 0xd0, 0x86,
	0xf2, 0x39, 0xec, 0x09, 0x5a, 0x72, 0xdf, 0xe9, 0xb6, 0xae, 0x63, 0x5d, 0x15, 0x4c, 0x4b, 0xee,
	0x9b, 0x2c, 0x6a, 0x03, 0xbd, 0x84, 0x51, 0x46, 0x85, 0x60, 0xd2, 0x4b, 0x4d, 0x0b, 0x1e, 0x78,
	0xe9, 0x85, 0x86, 0x9d, 0x78, 0x98, 0x35, 0xcd, 0xf8, 0x9f, 0x16, 0x0c, 0x6d, 0x05, 0xec, 0xb0,
	0x6e, 0x2a, 0xc0, 0x31, 0xb4, 0x49, 0xe2, 0xe7, 0xaf, 0xee, 0xee, 0x69, 0x92, 0x63, 0x8d, 0x6c,
	0x3f, 0x6b, 0x33, 0x80, 0x8a, 0xe6, 0xd4, 0x2a, 0xcc, 0xb4, 0x21, 0xaf, 0xb8, 0x74, 0x10, 0x6e,
	0xb0, 0xd0, 0xb7, 0xb0, 0xab, 0x52, 0x57, 0x29, 0x86, 0x9f, 0xac, 0x8e, 0xa1, 0x6c, 0xea, 0x59,
	0x67, 0xeb, 0x9e, 0x75, 0xef, 0xeb, 0x59, 0xfc, 0x6f, 0x0b, 0x06, 0x8d, 0x93, 0xd1, 0x53, 0x68,
	0x7f, 0x60, 0x45, 0xaa, 0x2b, 0x37, 0x6a, 0xdc, 0xd2, 0x06, 0x67, 0xf2, 0x33, 0x2b, 0x52, 0xac,
	0x69, 0xe8, 0x05, 0x0c, 0xb2, 0x46, 0xdb, 0x5a, 0xba, 0x6d, 0xf5, 0x85, 0xbb, 0x60, 0xb2, 0x99,
	0x16, 0x64, 0xbe, 0x69, 0x1b, 0x7a, 0x1e, 0x6e, 0xd1, 0x73, 0xf4, 0x83, 0x97, 0xdb, 0x4b, 0x6c,
	0xa7, 0xed, 0xae, 0xdc, 0xce, 0xc5, 0xeb, 0x1d, 0xb7, 0x81, 0x75, 0xa0, 0xa7, 0xd0, 0x53, 0x91,
	0xab, 0x42, 0xdb, 0x69, 0x1b, 0xdf, 0x0d, 0xfb, 0xf5, 0x0e, 0xee, 0x66, 0x66, 0x19, 0xef, 0x43,
	0x5b, 0xa5, 0x8d, 0xba, 0x10, 0x5e, 0xfc, 0x74, 0x35, 0xde, 0x79, 0xd5, 0x83, 0x4e, 0xc5, 0xd7,
	0x22, 0xa1, 0xf1, 0x2b, 0x18, 0xdd, 0xce, 0x13, 0x7d, 0xa1, 0x9e, 0x93, 0x92, 0x57, 0x4c, 0x72,
	0x71, 0x63, 0x87, 0xb0, 0xe1, 0x41, 0x63, 0x08, 0x05, 0x5d, 0xda, 0xfb, 0xa7, 0x96, 0xf1, 0x5f,
	0x01, 0x0c, 0x6f, 0xe5, 0xab, 0x46, 0xf8, 0x3d, 0xaf, 0xa4, 0x1b, 0x61, 0xb5, 0x46, 0x11, 0x74,
	0x5d, 0xb6, 0x46, 0xeb, 0x4c, 0xf4, 0x04, 0x0e, 0x52, 0x56, 0x91, 0x77, 0x39, 0xd5, 0x19, 0x09,
	0xc9, 0x8a, 0x4c, 0x17, 0xb4, 0x87, 0xc7, 0x16, 0xc0, 0xce, 0xaf, 0xc8, 0x1f, 0xdf, 0x33, 0x49,
	0x73, 0x56, 0x49, 0x9a, 0x2e, 0x32, 0xc1, 0xd7, 0xa5, 0x9e, 0xde, 0x3e, 0x1e, 0x37, 0x80, 0x0b,
	0xe5, 0x8f, 0x7f, 0x73, 0x81, 0xb9, 0xc2, 0x6d, 0x17, 0xd8, 0x43, 0x50, 0x25, 0x5c, 0xac, 0x45,
	0x6e, 0x1f, 0x8f, 0x4e, 0xc6, 0xe4, 0x5b, 0x91, 0xc7, 0x9f, 0x43, 0xd7, 0x56, 0x4d, 0x95, 0xa3,
	0xc6, 0xd5, 0x32, 0xfe, 0x23, 0x80, 0xf0, 0x34, 0xc9, 0xd1, 0x63, 0x68, 0x0b, 0x9e, 0x53, 0x3b,
	0x8d, 0x07, 0xcd, 0x3b, 0x3b, 0xc1, 0x3c, 0xa7, 0x58, 0xc3, 0xe8, 0x10, 0x76, 0x4d, 0x12, 0xe6,
	0x70, 0x63, 0xa0, 0x23, 0xe8, 0xb1, 0x94, 0x16, 0x92, 0xc9, 0x1b, 0xbb, 0xb7, 0xb7, 0xe3, 0x47,
	0xd0, 0x56, 0x7a, 0x04, 0xd0, 0xc1, 0xf3, 0xd3, 0xf3, 0x39, 0x1e, 0xef, 0xa0, 0x21, 0xf4, 0xf1,
	0xfc, 0xd7, 0xb7, 0xf3, 0xcb, 0xab, 0x39, 0x1e, 0x07, 0xf1, 0x9f, 0x01, 0xf4, 0xfd, 0x15, 0x56,
	0x9b, 0xb9, 0x7b, 0x6f, 0x33, 0xf7, 0x36, 0x9a, 0x41, 0xcf, 0x7d, 0x77, 0x3a, 0x82, 0x51, 0x63,
	0x0a, 0xfd, 0xcf, 0xf8, 0x0b, 0x59, 0x51, 0xec, 0x79, 0xe8, 0x1b, 0xe8, 0x9a, 0x4f, 0xde, 0xbd,
	0x34, 0xfb, 0xf5, 0xe7, 0xa1, 0xfd, 0xd8, 0xe1, 0xf1, 0x0c, 0x3a, 0xf7, 0x3c, 0x6b, 0x87, 0xb0,
	0x7b, 0x4d, 0xf2, 0x35, 0x75, 0xb9, 0x6b, 0xe3, 0x5d, 0x47, 0x7f, 0xe3, 0xcf, 0xfe, 0x0f, 0x00,
	0x00, 0xff, 0xff, 0x8b, 0x91, 0x6e, 0x34, 0x5f, 0x08, 0x00, 0x00,
}
