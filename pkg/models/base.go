package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KnownTYPES string

const (
	EDS        KnownTYPES = "endpoint"
	CDS        KnownTYPES = "cluster"
	LDS        KnownTYPES = "listener"
	ROUTE      KnownTYPES = "route"
	EXTENSIONS KnownTYPES = "extensions"
	FILTERS    KnownTYPES = "filters"
	ACCESSLOG  KnownTYPES = "access_log"
)

func (kt KnownTYPES) String() string {
	return string(kt)
}

type ScenarioBody map[string]interface{}

type DBResourceClass interface {
	GetGeneral() General
	GetGtype() GTypes
	SetGeneral(general *General)
	GetResource() interface{}
	SetResource(resource interface{})
	GetVersion() interface{}
	GetConfigDiscovery() []*ConfigDiscovery
	GetTypedConfig() []*TypedConfig
	SetTypedConfig(typedConfig []*TypedConfig)
	SetVersion(versionRaw interface{})
	SetPermissions(permissions *Permissions)

	SetBootstrapClusters(clusters []interface{})
	SetBootstrapAccessLoggers(accessLoggers []interface{})
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
	Metadata      map[string]string
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
	Metadata        map[string]interface{} `json:"metadata" bson:"metadata"`
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
	Disabled      bool   `json:"disabled" bson:"disabled"`
	Priority      int    `json:"priority" bson:"priority"`
	ParentName    string `json:"parent_name" bson:"parent_name"`
}

type TC struct {
	Name        string                 `json:"name" bson:"name"`
	TypedConfig map[string]interface{} `json:"typed_config" bson:"typed_config"`
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

func (d *DBResource) SetTypedConfig(typedConfig []*TypedConfig) {
	d.General.TypedConfig = typedConfig
}

func (d *DBResource) SetVersion(versionRaw interface{}) {
	version, ok := versionRaw.(string)
	if !ok {
		d.Resource.Version = "0"
		return
	}
	d.Resource.Version = version
}

func (d *DBResource) SetResource(resource interface{}) {
	d.Resource.Resource = resource
}

func (d *DBResource) SetBootstrapClusters(clusters []interface{}) {
	resourceMap, ok := d.Resource.Resource.(primitive.M)
	if !ok {
		fmt.Errorf("failed to parse Resource.Resource as map[string]interface{}, got type: %T", d.Resource.Resource)
	}

	//fmt.Printf("Type: %T\nValue: %+v\n", resourceMap, resourceMap)
	staticResources, ok := resourceMap["static_resources"].(primitive.M)
	if !ok || staticResources == nil {
		staticResources = make(primitive.M)
	}

	staticResources["clusters"] = clusters
	resourceMap["static_resources"] = staticResources
	d.Resource.Resource = resourceMap
}

func (d *DBResource) SetBootstrapAccessLoggers(accessLoggers []interface{}) {
	resourceMap, ok := d.Resource.Resource.(primitive.M)
	if !ok {
		fmt.Errorf("failed to parse Resource.Resource as map[string]interface{}, got type: %T", d.Resource.Resource)
	}

	admin, ok := resourceMap["admin"].(primitive.M)
	if !ok || admin == nil {
		admin = make(primitive.M)
	}

	admin["access_log"] = accessLoggers
	resourceMap["admin"] = admin
	d.Resource.Resource = resourceMap
}

func (d *DBResource) SetGeneral(general *General) {
	d.General = *general
}

func (d *DBResource) SetPermissions(permissions *Permissions) {
	d.General.Permissions = *permissions
}
