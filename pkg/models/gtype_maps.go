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
	brotli_compressor "github.com/envoyproxy/go-control-plane/envoy/extensions/compression/brotli/compressor/v3"
	gzip_compressor "github.com/envoyproxy/go-control-plane/envoy/extensions/compression/gzip/compressor/v3"
	zstd_compressor "github.com/envoyproxy/go-control-plane/envoy/extensions/compression/zstd/compressor/v3"
	bandwidth_limit "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/bandwidth_limit/v3"
	basic_auth "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/basic_auth/v3"
	compressor "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/compressor/v3"
	cors "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tcp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	hcefs "github.com/envoyproxy/go-control-plane/envoy/extensions/health_check/event_sinks/file/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	http_protocol_options "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"
	"github.com/sefaphlvn/bigbang/pkg/models/downstream_filters"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type GTypeMapping struct {
	Collection            string
	URL                   string
	PrettyName            string
	Message               proto.Message
	DownstreamFiltersFunc func(string) []downstream_filters.MongoFilters
	TypedConfigPaths      []TypedConfigPath
	UpstreamPaths         map[string]GTypes
}

var URLs = map[string]string{
	"bootstrap":             "/resource/bootstrap/",
	"clusters":              "/resource/cluster/",
	"endpoints":             "/resource/endpoint/",
	"listeners":             "/resource/listener/",
	"routes":                "/resource/route/",
	"virtual_host":          "/resource/virtual_host/",
	"tcp_proxy":             "/filters/network/tcp_proxy/",
	"hcm":                   "/filters/network/hcm/",
	"secrets":               "/resource/secret/",
	"access_log":            "/others/access_log/",
	"http_router":           "/filters/http/http_router/",
	"hcefs":                 "/others/hcefs/",
	"basic_auth":            "/filters/http/basic_auth/",
	"cors":                  "/filters/http/cors/",
	"bandwidth_limit":       "/filters/http/bandwidth_limit/",
	"compressor":            "/filters/http/compressor/",
	"compressor_library":    "/others/compressor_library/",
	"http_protocol_options": "/others/http_protocol_options/",
}

var gTypeMappings = map[GTypes]GTypeMapping{
	BootStrap: {
		PrettyName:            "Bootstrap",
		Collection:            "bootstrap",
		URL:                   URLs["bootstrap"],
		Message:               &bootstrap.Bootstrap{},
		DownstreamFiltersFunc: nil,
		TypedConfigPaths:      BootstrapTypedConfigPaths,
		UpstreamPaths:         nil,
	},
	HTTPConnectionManager: {
		PrettyName:            "Http Connection Manager",
		Collection:            "extensions",
		URL:                   URLs["hcm"],
		Message:               &hcm.HttpConnectionManager{},
		DownstreamFiltersFunc: downstream_filters.ConfigDiscoveryListenerDownstreamFilters,
		TypedConfigPaths:      HttpConnectionManagerTypedConfigPaths,
		UpstreamPaths:         HTTPConnectionManagerUpstreams,
	},
	Router: {
		PrettyName:            "Router",
		Collection:            "extensions",
		URL:                   URLs["http_routerhcm"],
		Message:               &router.Router{},
		DownstreamFiltersFunc: downstream_filters.ConfigDiscoveryHttpFilterDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	Cluster: {
		PrettyName:            "Cluster",
		Collection:            "clusters",
		URL:                   URLs["clusters"],
		Message:               &cluster.Cluster{},
		DownstreamFiltersFunc: downstream_filters.ClusterDownstreamFilters,
		TypedConfigPaths:      ClusterTypedConfigPaths,
		UpstreamPaths:         ClusterUpstreams,
	},
	Listener: {
		PrettyName:            "Listener",
		Collection:            "listeners",
		URL:                   URLs["listeners"],
		Message:               &listener.Listener{},
		DownstreamFiltersFunc: nil,
		TypedConfigPaths:      ListenerTypedConfigPaths,
		UpstreamPaths:         nil,
	},
	Endpoint: {
		PrettyName:            "Endpoint",
		Collection:            "endpoints",
		URL:                   URLs["endpoints"],
		Message:               &endpoint.ClusterLoadAssignment{},
		DownstreamFiltersFunc: downstream_filters.EdsDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	Route: {
		PrettyName:            "Route",
		Collection:            "routes",
		URL:                   URLs["routes"],
		Message:               &route.RouteConfiguration{},
		DownstreamFiltersFunc: downstream_filters.RouteDownstreamFilters,
		TypedConfigPaths:      RouteTypedConfigPaths,
		UpstreamPaths:         RouteUpstreams,
	},
	VirtualHost: {
		PrettyName:            "Virtual Host",
		Collection:            "virtual_host",
		URL:                   URLs["virtual_host"],
		Message:               &route.VirtualHost{},
		DownstreamFiltersFunc: downstream_filters.VirtualHostDownstreamFilters,
		TypedConfigPaths:      VirtualHostTypedConfigPaths,
		UpstreamPaths:         VirtualHostUpstreams,
	},
	TcpProxy: {
		PrettyName:            "Tcp Proxy",
		Collection:            "extensions",
		URL:                   URLs["tcp_proxy"],
		Message:               &tcp.TcpProxy{},
		DownstreamFiltersFunc: downstream_filters.ConfigDiscoveryListenerDownstreamFilters,
		TypedConfigPaths:      GeneralAccessLogTypedConfigPaths,
		UpstreamPaths:         TcpProxyUpstreams,
	},
	FluentdAccessLog: {
		PrettyName:            "Access Log(Fluentd)",
		Collection:            "others",
		URL:                   URLs["access_log"],
		Message:               &al_fluentd.FluentdAccessLogConfig{},
		DownstreamFiltersFunc: downstream_filters.ALSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         FluentdAccessLogUpstreams,
	},
	FileAccessLog: {
		PrettyName:            "Access Log(File)",
		Collection:            "others",
		URL:                   URLs["access_log"],
		Message:               &al_file.FileAccessLog{},
		DownstreamFiltersFunc: downstream_filters.ALSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	StdoutAccessLog: {
		PrettyName:            "Access Log(StdOut)",
		Collection:            "others",
		URL:                   URLs["access_log"],
		Message:               &al_stream.StdoutAccessLog{},
		DownstreamFiltersFunc: downstream_filters.ALSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	StdErrAccessLog: {
		PrettyName:            "Access Log(StdErr)",
		Collection:            "others",
		URL:                   URLs["access_log"],
		Message:               &al_stream.StderrAccessLog{},
		DownstreamFiltersFunc: downstream_filters.ALSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	DownstreamTlsContext: {
		PrettyName:            "Downstream TLS",
		Collection:            "secrets",
		URL:                   URLs["secrets"],
		Message:               &tls.DownstreamTlsContext{},
		DownstreamFiltersFunc: downstream_filters.DownstreamTlsDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         DownstreamTlsContextUpstreams,
	},
	UpstreamTlsContext: {
		PrettyName:            "Upstream TLS",
		Collection:            "secrets",
		URL:                   URLs["secrets"],
		Message:               &tls.UpstreamTlsContext{},
		DownstreamFiltersFunc: downstream_filters.UpstreamTlsDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         UpstreamTlsContextUpstreams,
	},
	TlsCertificate: {
		PrettyName:            "TLS Certificate",
		Collection:            "secrets",
		URL:                   URLs["secrets"],
		Message:               &tls.TlsCertificate{},
		DownstreamFiltersFunc: downstream_filters.TlsCertificateDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	CertificateValidationContext: {
		PrettyName:            "Certificate Validation",
		Collection:            "secrets",
		URL:                   URLs["secrets"],
		Message:               &tls.CertificateValidationContext{},
		DownstreamFiltersFunc: downstream_filters.ContextValidateDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	HealthCheckEventFileSink: {
		PrettyName:            "Health Check Event File Sink",
		Collection:            "others",
		URL:                   URLs["hcefs"],
		Message:               &hcefs.HealthCheckEventFileSink{},
		DownstreamFiltersFunc: downstream_filters.HCEFSDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	BasicAuth: {
		PrettyName:            "Basic Auth",
		Collection:            "extensions",
		URL:                   URLs["basic_auth"],
		Message:               &basic_auth.BasicAuth{},
		DownstreamFiltersFunc: downstream_filters.ConfigDiscoveryHttpFilterDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	BasicAuthPerRoute: {
		PrettyName:            "Basic Auth Per Route",
		Collection:            "extensions",
		URL:                   URLs["basic_auth"],
		Message:               &basic_auth.BasicAuthPerRoute{},
		DownstreamFiltersFunc: downstream_filters.TypedHttpFilterDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	Cors: {
		PrettyName:            "Cors",
		Collection:            "extensions",
		URL:                   URLs["cors"],
		Message:               &cors.Cors{},
		DownstreamFiltersFunc: downstream_filters.ConfigDiscoveryHttpFilterDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	CorsPolicy: {
		PrettyName:            "Cors Policy",
		Collection:            "extensions",
		URL:                   URLs["cors"],
		Message:               &cors.CorsPolicy{},
		DownstreamFiltersFunc: downstream_filters.TypedHttpFilterDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	BandwidthLimit: {
		PrettyName:            "Bandwidth Limit",
		Collection:            "extensions",
		URL:                   URLs["bandwidth_limit"],
		Message:               &bandwidth_limit.BandwidthLimit{},
		DownstreamFiltersFunc: downstream_filters.DiscoverAndTypedHttpFilterDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	Compressor: {
		PrettyName:            "Compressor",
		Collection:            "extensions",
		URL:                   URLs["compressor"],
		Message:               &compressor.Compressor{},
		DownstreamFiltersFunc: downstream_filters.ConfigDiscoveryHttpFilterDownstreamFilters,
		TypedConfigPaths:      CompressorTypedConfigPaths,
		UpstreamPaths:         nil,
	},
	CompressorPerRoute: {
		PrettyName:            "Compressor Per Route",
		Collection:            "extensions",
		URL:                   URLs["compressor"],
		Message:               &compressor.CompressorPerRoute{},
		DownstreamFiltersFunc: downstream_filters.TypedHttpFilterDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	GzipCompressor: {
		PrettyName:            "Gzip Compressor",
		Collection:            "others",
		URL:                   URLs["compressor_library"],
		Message:               &gzip_compressor.Gzip{},
		DownstreamFiltersFunc: downstream_filters.CompressorLibraryDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	BrotliCompressor: {
		PrettyName:            "Brotli Compressor",
		Collection:            "others",
		URL:                   URLs["compressor_library"],
		Message:               &brotli_compressor.Brotli{},
		DownstreamFiltersFunc: downstream_filters.CompressorLibraryDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	ZstdCompressor: {
		PrettyName:            "Zstd Compressor",
		Collection:            "others",
		URL:                   URLs["compressor_library"],
		Message:               &zstd_compressor.Zstd{},
		DownstreamFiltersFunc: downstream_filters.CompressorLibraryDownstreamFilters,
		TypedConfigPaths:      nil,
		UpstreamPaths:         nil,
	},
	HttpProtocolOptions: {
		PrettyName:            "Http Protocol Options",
		Collection:            "others",
		URL:                   URLs["http_protocol_options"],
		Message:               &http_protocol_options.HttpProtocolOptions{},
		DownstreamFiltersFunc: downstream_filters.TypedHttpProtocolDownstreamFilters,
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

func (gt GTypes) DownstreamFilters(name string) []downstream_filters.MongoFilters {
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

func (gt GTypes) Validate() map[string]GTypes {
	if mapping, exists := gTypeMappings[gt]; exists && mapping.UpstreamPaths != nil {
		return mapping.UpstreamPaths
	}
	return nil
}
