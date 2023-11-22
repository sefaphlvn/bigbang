package resources

import (
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/sefaphlvn/bigbang/grpc/models"
	"github.com/sefaphlvn/bigbang/pkg/db"
	"github.com/sirupsen/logrus"
)

type UsedResources struct {
	Name string
	Type string
}

type AllResources struct {
	UsedResource []UsedResources
	NodeID       string
	Version      string
	Listener     []*listener.Listener
	Cluster      []*cluster.Cluster
	Route        route.RouteConfiguration
	Endpoint     []*endpoint.Endpoint
	Secret       tls.Secret
	Extensions   []*core.TypedExtensionConfig
}

func NewResources() *AllResources {
	return &AllResources{}
}

func SetSnapshot(cur *models.DBResource, nodeID string, db *db.MongoDB, logger *logrus.Logger) (*AllResources, error) {
	resourceAll := NewResources()
	resourceAll.NodeID = nodeID
	resourceAll.DecodeListener(cur, db, logger)
	return resourceAll, nil
}
