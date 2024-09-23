package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KnownTYPES string

const (
	EDS        KnownTYPES = "endpoints"
	CDS        KnownTYPES = "clusters"
	LDS        KnownTYPES = "listeners"
	ROUTE      KnownTYPES = "routes"
	EXTENSIONS KnownTYPES = "extensions"
	ACCESSLOG  KnownTYPES = "access_log"
)

func (kt KnownTYPES) String() string {
	return string(kt)
}

type DBResourceClass interface {
	GetGeneral() General
	GetGtype() GTypes
	SetGeneral(*General)
	GetResource() interface{}
	SetResource(interface{})
	GetVersion() interface{}
	GetConfigDiscovery() []*ConfigDiscovery
	GetTypedConfig() []*TypedConfig
	SetTypedConfig([]*TypedConfig)
	SetVersion(interface{})
	SetPermissions(*Permissions)
}

type RequestDetails struct {
	Collection    string
	Type          KnownTYPES
	GType         GTypes
	CanonicalName string
	Name          string
	Category      string
	Version       string
	User          UserDetails
	SaveOrPublish string
	Project       string
}

type UserDetails struct {
	Groups    []string
	Projects  []string
	BaseGroup string
	Role      Role
	IsOwner   bool
	UserID    string
	UserName  string
}

type Service struct {
	Name     string    `json:"name" bson:"name"`
	Machines []Machine `json:"Machines" bson:"Machines"`
}

type Machine struct {
	Name              string `json:"name" bson:"name"`
	Ifname            string `json:"ifname" bson:"ifname"`
	DownstreamAddress string `json:"downstream_address" bson:"downstream_address"`
}

type General struct {
	Name            string                 `json:"name" bson:"name"`
	Version         string                 `json:"version" bson:"version"`
	Type            KnownTYPES             `json:"type" bson:"type"`
	GType           GTypes                 `json:"gtype" bson:"gtype"`
	Project         string                 `json:"project" bson:"project"`
	Collection      string                 `json:"collection" bson:"collection"`
	CanonicalName   string                 `json:"canonical_name" bson:"canonical_name"`
	Category        string                 `json:"category" bson:"category"`
	Managed         bool                   `json:"managed,omitempty" bson:"managed,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Permissions     Permissions            `json:"permissions" bson:"permissions"`
	ConfigDiscovery []*ConfigDiscovery     `json:"config_discovery,omitempty" bson:"config_discovery,omitempty"`
	TypedConfig     []*TypedConfig         `json:"typed_config,omitempty" bson:"typed_config,omitempty"`
	CreatedAt       primitive.DateTime     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       primitive.DateTime     `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Permissions struct {
	Users  []string `json:"users" bson:"users"`
	Groups []string `json:"groups" bson:"groups"`
}

type ConfigDiscovery struct {
	ParentName    string `json:"parent_name,omitempty" bson:"parent_name,omitempty"`
	GType         GTypes `json:"gtype" bson:"gtype"`
	Name          string `json:"name" bson:"name"`
	Priority      int    `json:"priority" bson:"priority"`
	Category      string `json:"category" bson:"category"`
	CanonicalName string `json:"canonical_name" bson:"canonical_name"`
}

type TypedConfig struct {
	Name          string `json:"name" bson:"name"`
	CanonicalName string `json:"canonical_name" bson:"canonical_name"`
	Gtype         GTypes `json:"gtype" bson:"gtype"`
	Type          string `json:"type" bson:"type"`
	Category      string `json:"category" bson:"category"`
	Collection    string `json:"collection" bson:"collection"`
}

type DBResource struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	General  General            `json:"general" bson:"general"`
	Resource Resource           `json:"resource" bson:"resource"`
}

type Resource struct {
	Version  string      `json:"version" bson:"version"`
	Resource interface{} `json:"resource" bson:"resource"`
}

func (d *DBResource) GetGeneral() General {
	return d.General
}

func (d *DBResource) GetGtype() GTypes {
	return d.General.GType
}

func (d *DBResource) GetResource() interface{} {
	return d.Resource.Resource
}

func (d *DBResource) GetVersion() interface{} {
	return d.Resource.Version
}

func (d *DBResource) GetConfigDiscovery() []*ConfigDiscovery {
	return d.General.ConfigDiscovery
}

func (d *DBResource) GetTypedConfig() []*TypedConfig {
	return d.General.TypedConfig
}

func (d *DBResource) SetTypedConfig(res []*TypedConfig) {
	d.General.TypedConfig = res
}

func (d *DBResource) SetVersion(res interface{}) {
	d.Resource.Version = res.(string)
}

func (d *DBResource) SetResource(res interface{}) {
	d.Resource.Resource = res
}

func (d *DBResource) SetGeneral(g *General) {
	d.General = *g
}

func (d *DBResource) SetPermissions(p *Permissions) {
	d.General.Permissions = *p
}
