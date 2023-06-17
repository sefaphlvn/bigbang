package resources

import (
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/sefaphlvn/bigbang/grpcServer/db"
	"github.com/sefaphlvn/bigbang/grpcServer/models"
)

type ResourcesUsed struct {
	Name string
	Type string
}

type AllResources struct {
	UsedResource []ResourcesUsed
	NodeID       string
	Version      string
	Listener     []*listener.Listener
	Cluster      []*cluster.Cluster
	Route        route.RouteConfiguration
	Endpoint     []*endpoint.Endpoint
	Secret       tls.Secret
}

func NewResources() *AllResources {
	return &AllResources{}
}

func SetSnapshot(cur *models.Resource, nodeID string, db *db.MongoDB) (*AllResources, error) {
	resourceAll := NewResources()
	resourceAll.NodeID = nodeID
	resourceAll.DecodeListener(cur, db)
	return resourceAll, nil
}
