package resources

import (
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/sefaphlvn/bigbang/grpcServer/models"
)

type AllResources struct {
	Listener []*listener.Listener
	Cluster  []*cluster.Cluster
	Route    route.Route
	Endpoint []*endpoint.Endpoint
	Secret   tls.Secret
}

func NewResources() *AllResources {
	return &AllResources{}
}

func SetSnapshot(cur *models.Resource) (*AllResources, error) {
	resourceAll := NewResources()
	resourceAll.DecodeListener(cur)
	return resourceAll, nil
}
