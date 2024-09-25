package models

import (
	bootstrap "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	al_file "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/file/v3"
	al_fluentd "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/fluentd/v3"
	al_stream "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/stream/v3"
	basic_auth "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/basic_auth/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tcp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	hcefs "github.com/envoyproxy/go-control-plane/envoy/extensions/health_check/event_sinks/file/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/sefaphlvn/bigbang/pkg/filters"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type GTypeMapping struct {
	Collection            string
	URL                   string
	PrettyName            string
	Message               proto.Message
	DownstreamFiltersFunc func(string) []filters.MongoFilters
	TypedConfigPaths      []TypedConfigPath
	UpstreamPaths         map[string]GTypes
}

var gTypeMappings = map[GTypes]GTypeMapping{
	BootStrap: {
		PrettyName:            "Bootstrap",
		Collection:            "bootstrap",
		URL:                   "/resource/bootstrap/",
		Message:               &bootstrap.Bootstrap{},
		DownstreamFiltersFunc: nil,
		TypedConfigPaths:      BootstrapTypedConfigPaths,
		UpstreamPaths:         nil,
	},
	HTTPConnectionManager: {
		PrettyName:            "Http Connection Manager",
		Collection:            "extensions",
		URL:                   "/filters/network/hcm/",
		Message:               &hcm.HttpConnectionManager{},
		DownstreamFiltersFunc: filters.HcmDownstreamFilters,
		TypedConfigPaths:      GeneralAccessLogTypedConfigPaths,
		UpstreamPaths:         HTTPConnectionManagerUpstreams,
	},
	Router: {
		PrettyName:            "Router",
		Collection:            "extensions",
		URL:                   "/filters/http/http_router/",
		Message:               &router.Router{},
		DownstreamFiltersFunc: filters.RouterDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	Cluster: {
		PrettyName:            "Cluster",
		Collection:            "clusters",
		URL:                   "/resource/cluster/",
		Message:               &cluster.Cluster{},
		DownstreamFiltersFunc: filters.ClusterDownstreamFilters,
		TypedConfigPaths:      ClusterTypedConfigPaths,
		UpstreamPaths:         ClusterUpstreams,
	},
	Listener: {
		PrettyName:            "Listener",
		Collection:            "listeners",
		URL:                   "/resource/listener/",
		Message:               &listener.Listener{},
		DownstreamFiltersFunc: nil,
		TypedConfigPaths:      ListenerTypedConfigPaths,
		UpstreamPaths:         nil,
	},
	Endpoint: {
		PrettyName:            "Endpoint",
		Collection:            "endpoints",
		URL:                   "/resource/endpoint",
		Message:               &endpoint.ClusterLoadAssignment{},
		DownstreamFiltersFunc: filters.EdsDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	Route: {
		PrettyName:            "Route",
		Collection:            "routes",
		URL:                   "/resource/route",
		Message:               &route.RouteConfiguration{},
		DownstreamFiltersFunc: filters.RouteDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         RouteUpstreams,
	},
	VirtualHost: {
		PrettyName:            "Virtual Host",
		Collection:            "virtual_host",
		URL:                   "/resource/virtual_host",
		Message:               &route.VirtualHost{},
		DownstreamFiltersFunc: filters.VirtualHostDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         VirtualHostUpstreams,
	},
	TcpProxy: {
		PrettyName:            "Tcp Proxy",
		Collection:            "extensions",
		URL:                   "/filters/network/tcp_proxy/",
		Message:               &tcp.TcpProxy{},
		DownstreamFiltersFunc: filters.TcpProxyDownstreamFilters,
		TypedConfigPaths:      GeneralAccessLogTypedConfigPaths,
		UpstreamPaths:         TcpProxyUpstreams,
	},
	FluentdAccessLog: {
		PrettyName:            "Access Log(Fluentd)",
		Collection:            "others",
		URL:                   "/others/access_log/",
		Message:               &al_fluentd.FluentdAccessLogConfig{},
		DownstreamFiltersFunc: filters.ALSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         FluentdAccessLogUpstreams,
	},
	FileAccessLog: {
		PrettyName:            "Access Log(File)",
		Collection:            "others",
		URL:                   "/others/access_log/",
		Message:               &al_file.FileAccessLog{},
		DownstreamFiltersFunc: filters.ALSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	StdoutAccessLog: {
		PrettyName:            "Access Log(StdOut)",
		Collection:            "others",
		URL:                   "/others/access_log/",
		Message:               &al_stream.StdoutAccessLog{},
		DownstreamFiltersFunc: filters.ALSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	StdErrAccessLog: {
		PrettyName:            "Access Log(StdErr)",
		Collection:            "others",
		URL:                   "/others/access_log/",
		Message:               &al_stream.StderrAccessLog{},
		DownstreamFiltersFunc: filters.ALSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	DownstreamTlsContext: {
		PrettyName:            "Downstream TLS",
		Collection:            "secrets",
		URL:                   "/resource/secret/",
		Message:               &tls.DownstreamTlsContext{},
		DownstreamFiltersFunc: filters.DownstreamTlsDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         DownstreamTlsContextUpstreams,
	},
	UpstreamTlsContext: {
		PrettyName:            "Upstream TLS",
		Collection:            "secrets",
		URL:                   "/resource/secret/",
		Message:               &tls.UpstreamTlsContext{},
		DownstreamFiltersFunc: filters.UpstreamTlsDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         UpstreamTlsContextUpstreams,
	},
	TlsCertificate: {
		PrettyName:            "TLS Certificate",
		Collection:            "secrets",
		URL:                   "/resource/secret/",
		Message:               &tls.TlsCertificate{},
		DownstreamFiltersFunc: filters.TlsCertificateDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	CertificateValidationContext: {
		PrettyName:            "Certificate Validation",
		Collection:            "secrets",
		URL:                   "/resource/secret/",
		Message:               &tls.CertificateValidationContext{},
		DownstreamFiltersFunc: filters.ContextValidateDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	HealthCheckEventFileSink: {
		PrettyName:            "Health Check Event File Sink",
		Collection:            "others",
		URL:                   "/others/hcefs/",
		Message:               &hcefs.HealthCheckEventFileSink{},
		DownstreamFiltersFunc: filters.HCEFSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	BasicAuth: {
		PrettyName:            "Basic Auth",
		Collection:            "extensions",
		URL:                   "/filters/http/basic_auth/",
		Message:               &basic_auth.BasicAuth{},
		DownstreamFiltersFunc: filters.BasicAuthDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
}

func (gt GTypes) String() string {
	return string(gt)
}

func (gt GTypes) CollectionString() string {
	if mapping, exists := gTypeMappings[gt]; exists {
		return mapping.Collection
	}
	return "unknown"
}

func (gt GTypes) URL() string {
	if mapping, exists := gTypeMappings[gt]; exists {
		return mapping.URL
	}
	return "unknown"
}

func (gt GTypes) PrettyName() string {
	if mapping, exists := gTypeMappings[gt]; exists {
		return mapping.PrettyName
	}
	return "unknown"
}

func (gt GTypes) ProtoMessage() proto.Message {
	if mapping, exists := gTypeMappings[gt]; exists {
		return mapping.Message
	}
	return &anypb.Any{}
}

func (gt GTypes) DownstreamFilters(name string) []filters.MongoFilters {
	if mapping, exists := gTypeMappings[gt]; exists && mapping.DownstreamFiltersFunc != nil {
		return mapping.DownstreamFiltersFunc(name)
	}
	return nil
}

func (gt GTypes) TypedConfigPaths() []TypedConfigPath {
	if mapping, exists := gTypeMappings[gt]; exists && mapping.TypedConfigPaths != nil {
		return mapping.TypedConfigPaths
	}
	return nil
}

func (gt GTypes) UpstreamPaths() map[string]GTypes {
	if mapping, exists := gTypeMappings[gt]; exists && mapping.UpstreamPaths != nil {
		return mapping.UpstreamPaths
	}
	return nil
}
