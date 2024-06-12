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
)

func (kt KnownTYPES) String() string {
	return string(kt)
}

type DBResourceClass interface {
	GetGeneral() General
	SetGeneral(*General)
	GetResource() interface{}
	SetResource(interface{})
	GetVersion() interface{}
	GetConfigDiscovery() []*ConfigDiscovery
	GetTypedConfig() []*TypedConfig
	SetTypedConfig([]*TypedConfig)
	SetVersion(interface{})
}

type ResourceDetails struct {
	Collection    string
	Type          KnownTYPES
	GType         GTypes
	CanonicalName string
	Name          string
	Category      string
	Version       string
	User          UserDetails
	SaveOrPublish string
}

type UserDetails struct {
	Groups  []string
	Role    string
	IsAdmin bool
	UserID  string
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
	CanonicalName   string                 `json:"canonical_name" bson:"canonical_name"`
	Category        string                 `json:"category" bson:"category"`
	Service         GeneralService         `json:"service,omitempty" bson:"service,omitempty"`
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

type GeneralService struct {
	Name    string `json:"name" bson:"name"`
	Enabled bool   `json:"enabled" bson:"enabled"`
}

type Extensions struct {
	GType         GTypes `json:"gtype" bson:"gtype"`
	Name          string `json:"name" bson:"name"`
	Priority      int    `json:"priority" bson:"priority"`
	Category      string `json:"category" bson:"category"`
	CanonicalName string `json:"canonical_name" bson:"canonical_name"`
}

type ConfigDiscovery struct {
	ParentName   string       `json:"parent_name,omitempty" bson:"parent_name,omitempty"`
	Extensions   []Extensions `json:"extensions,omitempty" bson:"extensions,omitempty"`
	MainResource string       `json:"main_resource,omitempty" bson:"main_resource,omitempty"`
}

type TypedConfig struct {
	Name          string `json:"name" bson:"name"`
	CanonicalName string `json:"canonical_name" bson:"canonical_name"`
	Gtype         string `json:"gtype" bson:"gtype"`
	Type          string `json:"type" bson:"type"`
	Category      string `json:"category" bson:"category"`
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
