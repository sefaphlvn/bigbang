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
	SetVersion(interface{})
}

type ResourceDetails struct {
	Collection string
	Type       string
	SubType    string
	Name       string
	Version    string
	User       UserDetails
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
	Name      string                 `json:"name" bson:"name"`
	Version   string                 `json:"version" bson:"version"`
	Type      string                 `json:"type" bson:"type"`
	SubType   string                 `json:"subtype" bson:"subtype"`
	Extra     map[string]interface{} `json:"extra" bson:"extra"`
	Groups    []string               `json:"groups" bson:"groups"`
	CreatedAt primitive.DateTime     `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt primitive.DateTime     `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
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

func (d *DBResource) SetVersion(res interface{}) {
	d.Resource.Version = res.(string)
}

func (d *DBResource) SetResource(res interface{}) {
	d.Resource.Resource = res
}

func (d *DBResource) SetGeneral(g *General) {
	d.General = *g
}
