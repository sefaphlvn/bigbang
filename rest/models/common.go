package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBResourceClass interface {
	GetGeneral() General
	SetGeneral(*General)
	GetResource() interface{}
	SetResource(interface{})
	GetVersion() interface{}
	GetAdditionalResources() []AdditionalResource
	SetVersion(interface{})
}

type ResourceDetails struct {
	Collection    string
	Type          string
	CanonicalName string
	Name          string
	Category      string
	Version       string
	User          UserDetails
}

type UserDetails struct {
	Groups  []string
	IsAdmin bool
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
	Name                string                 `json:"name" bson:"name"`
	Version             string                 `json:"version" bson:"version"`
	Type                string                 `json:"type" bson:"type"`
	GType               string                 `json:"gtype" bson:"gtype"`
	CanonicalName       string                 `json:"canonical_name" bson:"canonical_name"`
	Category            string                 `json:"category" bson:"category"`
	Extra               map[string]interface{} `json:"extra,omitempty" bson:"extra,omitempty"`
	Groups              []string               `json:"groups" bson:"groups"`
	AdditionalResources []AdditionalResource   `json:"additional_resources,omitempty" bson:"additional_resources,omitempty"`
	CreatedAt           primitive.DateTime     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt           primitive.DateTime     `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Extra struct {
	Agent   string `json:"agent,omitempty" bson:"agent,omitempty"`
	Team    string `json:"team,omitempty" bson:"team,omitempty"`
	Service string `json:"service,omitempty" bson:"service,omitempty"`
}

type Extensions struct {
	GType         string `json:"gtype" bson:"gtype"`
	Name          string `json:"name" bson:"name"`
	Priority      int    `json:"priority" bson:"priority"`
	Category      string `json:"category" bson:"category"`
	CanonicalName string `json:"canonical_name" bson:"canonical_name"`
}

type AdditionalResource struct {
	ParentName   string       `json:"parent_name,omitempty" bson:"parent_name,omitempty"`
	Extensions   []Extensions `json:"extensions,omitempty" bson:"extensions,omitempty"`
	MainResource string       `json:"main_resource,omitempty" bson:"main_resource,omitempty"`
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

func (d *DBResource) GetAdditionalResources() []AdditionalResource {
	return d.General.AdditionalResources
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
