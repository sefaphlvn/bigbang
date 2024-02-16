package models

type GTypes string

const (
	APITypePrefix         GTypes = "type.googleapis.com/"
	HTTPConnectionManager GTypes = "envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
	Router                GTypes = "envoy.extensions.filters.http.router.v3.Router"
	Cluster               GTypes = "envoy.config.cluster.v3.Cluster"
	Listener              GTypes = "envoy.config.listener.v3.Listener"
	Endpoint              GTypes = "envoy.config.endpoint.v3.ClusterLoadAssignment"
	Route                 GTypes = "envoy.config.route.v3.RouteConfiguration"
	TcpProxy              GTypes = "envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy"
)

func (kt GTypes) String() string {
	return string(kt)
}
