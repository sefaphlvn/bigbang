package common

import (
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

type Resources struct {
	NodeID     string
	Version    string
	Listener   []types.Resource
	Cluster    []types.Resource
	Route      []types.Resource
	Endpoint   []types.Resource
	Secret     *tls.Secret
	Extensions []types.Resource
}

type AllResources interface {
	SetNodeID(nodeID string)
	GetNodeID() string

	SetVersion(version string)
	GetVersion() string

	SetListener(listener []types.Resource)
	GetListener() []*listener.Listener
	GetListenerT() []types.Resource

	SetCluster(cluster []types.Resource)
	GetCluster() []*cluster.Cluster
	GetClusterT() []types.Resource

	SetRoute(route *route.RouteConfiguration)
	GetRoute() *route.RouteConfiguration
	AppendRoute(route *route.RouteConfiguration)

	SetEndpoint(endpoint []types.Resource)
	GetEndpoint() []*endpoint.Endpoint
	GetEndpointT() []types.Resource

	SetSecret(secret *tls.Secret)
	GetSecret() *tls.Secret

	SetExtensions(extensions []types.Resource)
	GetExtensions() []*core.TypedExtensionConfig
	GetExtensionsT() []types.Resource
}

func (ar *Resources) SetNodeID(nodeID string) {
	ar.NodeID = nodeID
}

func (ar *Resources) GetNodeID() string {
	return ar.NodeID
}

func (ar *Resources) SetVersion(version string) {
	ar.Version = version
}

func (ar *Resources) GetVersion() string {
	return ar.Version
}

func (ar *Resources) SetListener(listener []types.Resource) {
	ar.Listener = listener
}

func (ar *Resources) GetListener() []*listener.Listener {
	listeners := make([]*listener.Listener, len(ar.Listener))
	for i, res := range ar.Listener {
		listeners[i] = res.(*listener.Listener)
	}
	return listeners
}

func (ar *Resources) GetListenerT() []types.Resource {
	return ar.Listener
}

func (ar *Resources) SetCluster(cluster []types.Resource) {
	ar.Cluster = cluster
}

func (ar *Resources) GetCluster() []*cluster.Cluster {
	clusters := make([]*cluster.Cluster, len(ar.Cluster))
	for i, res := range ar.Cluster {
		clusters[i] = res.(*cluster.Cluster)
	}
	return clusters
}

func (ar *Resources) GetClusterT() []types.Resource {
	return ar.Cluster
}

func (ar *Resources) SetRoute(route []types.Resource) {
	ar.Route = route
}

func (ar *Resources) AppendRoute(route *route.RouteConfiguration) {
	ar.Route = append(ar.Route, route)
}

func (ar *Resources) GetRoute() []*route.RouteConfiguration {
	routes := make([]*route.RouteConfiguration, len(ar.Route))
	for i, res := range ar.Route {
		routes[i] = res.(*route.RouteConfiguration)
	}
	return routes
}

func (ar *Resources) GetRouteT() []types.Resource {
	return ar.Route
}

func (ar *Resources) SetEndpoint(endpoint []types.Resource) {
	ar.Endpoint = endpoint
}

func (ar *Resources) GetEndpoint() []*endpoint.Endpoint {
	endpoints := make([]*endpoint.Endpoint, len(ar.Endpoint))
	for i, res := range ar.Endpoint {
		endpoints[i] = res.(*endpoint.Endpoint)
	}
	return endpoints
}

func (ar *Resources) GetEndpointT() []types.Resource {
	return ar.Endpoint
}

func (ar *Resources) SetSecret(secret *tls.Secret) {
	ar.Secret = secret
}

func (ar *Resources) GetSecret() *tls.Secret {
	return ar.Secret
}

func (ar *Resources) SetExtensions(extensions []types.Resource) {
	ar.Extensions = extensions
}

func (ar *Resources) GetExtensions() []*core.TypedExtensionConfig {
	extensions := make([]*core.TypedExtensionConfig, len(ar.Extensions))
	for i, res := range ar.Extensions {
		extensions[i] = res.(*core.TypedExtensionConfig)
	}
	return extensions
}

func (ar *Resources) GetExtensionsT() []types.Resource {
	return ar.Extensions
}
