package models

import (
	"github.com/sefaphlvn/bigbang/pkg/filters"
)

type GTypes string

const (
	APITypePrefix                GTypes = "type.googleapis.com/"
	BootStrap                    GTypes = "envoy.config.bootstrap.v3.Bootstrap"
	HTTPConnectionManager        GTypes = "envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
	Router                       GTypes = "envoy.extensions.filters.http.router.v3.Router"
	Cluster                      GTypes = "envoy.config.cluster.v3.Cluster"
	Listener                     GTypes = "envoy.config.listener.v3.Listener"
	Endpoint                     GTypes = "envoy.config.endpoint.v3.ClusterLoadAssignment"
	Route                        GTypes = "envoy.config.route.v3.RouteConfiguration"
	TcpProxy                     GTypes = "envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy"
	FluentdAccessLog             GTypes = "envoy.extensions.access_loggers.fluentd.v3.FluentdAccessLogConfig"
	FileAccessLog                GTypes = "envoy.extensions.access_loggers.file.v3.FileAccessLog"
	StdoutAccessLog              GTypes = "envoy.extensions.access_loggers.stream.v3.StdoutAccessLog"
	StdErrAccessLog              GTypes = "envoy.extensions.access_loggers.stream.v3.StderrAccessLog"
	DownstreamTlsContext         GTypes = "envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext"
	UpstreamTlsContext           GTypes = "envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext"
	TlsCertificate               GTypes = "envoy.extensions.transport_sockets.tls.v3.TlsCertificate"
	CertificateValidationContext GTypes = "envoy.extensions.transport_sockets.tls.v3.CertificateValidationContext"
)

func (gt GTypes) CollectionString() string {
	if str, exists := gTypeMappings[gt]; exists {
		return str.Collection
	}
	return "unknown"
}

func (gt GTypes) URL() string {
	if str, exists := gTypeMappings[gt]; exists {
		return str.URL
	}
	return "unknown"
}

func (gt GTypes) PrettyName() string {
	if str, exists := gTypeMappings[gt]; exists {
		return str.PrettyName
	}
	return "unknown"
}

func (gt GTypes) String() string {
	return string(gt)
}

type GTypeMapping struct {
	Collection string
	URL        string
	PrettyName string
}

var gTypeMappings = map[GTypes]GTypeMapping{
	BootStrap:                    {PrettyName: "Bootstrap", Collection: "bootstrap", URL: "/resource/bootstrap/"},
	HTTPConnectionManager:        {PrettyName: "Http Connection Manager", Collection: "extensions", URL: "/filters/network/hcm/"},
	Router:                       {PrettyName: "Router", Collection: "extensions", URL: "/filters/http/http_router/"},
	Cluster:                      {PrettyName: "Cluster", Collection: "clusters", URL: "/resource/cluster/"},
	Listener:                     {PrettyName: "Listener", Collection: "listeners", URL: "/resource/listener/"},
	Endpoint:                     {PrettyName: "Endpoint", Collection: "endpoints", URL: "/resource/endpoint"},
	Route:                        {PrettyName: "Route", Collection: "routes", URL: "/resource/route"},
	TcpProxy:                     {PrettyName: "Tcp Proxy", Collection: "extensions", URL: "/filters/network/tcp_proxy/"},
	FluentdAccessLog:             {PrettyName: "Access Log(Fluentd)", Collection: "others", URL: "/others/access_log/"},
	FileAccessLog:                {PrettyName: "Access Log(File)", Collection: "others", URL: "/others/access_log/"},
	StdoutAccessLog:              {PrettyName: "Access Log(StdOut)", Collection: "others", URL: "/others/access_log/"},
	StdErrAccessLog:              {PrettyName: "Access Log(StdErr)", Collection: "others", URL: "/others/access_log/"},
	DownstreamTlsContext:         {PrettyName: "Downstream TLS", Collection: "secrets", URL: "/resource/secret/"},
	UpstreamTlsContext:           {PrettyName: "Upstream TLS", Collection: "secrets", URL: "/resource/secret/"},
	TlsCertificate:               {PrettyName: "TLS Certificate", Collection: "secrets", URL: "/resource/secret/"},
	CertificateValidationContext: {PrettyName: "Certificate Validation", Collection: "secrets", URL: "/resource/secret/"},
}

func (gt GTypes) GetUpstreamPaths() map[string]GTypes {
	switch gt {
	case Cluster:
		return map[string]GTypes{
			"resource.resource.eds_cluster_config.service_name": Endpoint,
		}
	case TcpProxy:
		return map[string]GTypes{
			"resource.resource.cluster":                           Cluster,
			"resource.resource.weighted_clusters.clusters.#.name": Cluster,
		}
	case HTTPConnectionManager:
		return map[string]GTypes{
			"resource.resource.rds.route_config_name": Route,
		}
	case Route:
		return map[string]GTypes{
			"resource.resource.virtual_hosts.#.routes.#.route.cluster":                           Cluster,
			"resource.resource.virtual_hosts.#.routes.#.route.weighted_clusters.clusters.#.name": Cluster,
			"resource.resource.virtual_hosts.#.request_mirror_policies.#.cluster":                Cluster,
			"resource.resource.request_mirror_policies.#.cluster":                                Cluster,
		}
	case FluentdAccessLog:
		return map[string]GTypes{
			"resource.resource.cluster": Cluster,
		}
	case DownstreamTlsContext:
		return map[string]GTypes{
			"resource.resource.common_tls_context.tls_certificate_sds_secret_configs.#.name": TlsCertificate,
			"resource.resource.common_tls_context.validation_context_sds_secret_config.name": CertificateValidationContext,
		}
	case UpstreamTlsContext:
		return map[string]GTypes{
			"resource.resource.common_tls_context.tls_certificate_sds_secret_configs.#.name": TlsCertificate,
			"resource.resource.common_tls_context.validation_context_sds_secret_config.name": CertificateValidationContext,
		}
	default:
		return nil
	}
}

func (gt GTypes) GetDownstreamFilters(name string) []filters.MongoFilters {
	switch gt {
	case TcpProxy:
		return []filters.MongoFilters{
			filters.TcpProxyDownstreamFilters(name),
		}
	case Route:
		return []filters.MongoFilters{
			filters.RouteDownstreamFilters(name),
		}
	case HTTPConnectionManager:
		return []filters.MongoFilters{
			filters.HcmDownstreamFilters(name),
		}
	case Cluster:
		return filters.ClusterDownstreamFilters(name)
	case Endpoint:
		return []filters.MongoFilters{
			filters.EdsDownstreamFilters(name),
		}
	case FileAccessLog, FluentdAccessLog, StdErrAccessLog, StdoutAccessLog:
		return filters.ALSDownstreamFilters(name)
	case Router:
		return []filters.MongoFilters{
			filters.RouterDownstreamFilters(name),
		}
	case DownstreamTlsContext:
		return []filters.MongoFilters{
			filters.DownstreamTlsDownstreamFilters(name),
		}
	case CertificateValidationContext:
		return []filters.MongoFilters{
			filters.ContextValidateDownstreamFilters(name),
		}
	case TlsCertificate:
		return []filters.MongoFilters{
			filters.TlsCertificateDownstreamFilters(name),
		}
	default:
		return nil
	}
}
