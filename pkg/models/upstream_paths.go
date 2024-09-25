package models

var ClusterUpstreams = map[string]GTypes{
	"eds_cluster_config.service_name": Endpoint,
}

var TcpProxyUpstreams = map[string]GTypes{
	"cluster":                           Cluster,
	"weighted_clusters.clusters.#.name": Cluster,
}

var HTTPConnectionManagerUpstreams = map[string]GTypes{
	"rds.route_config_name":                                                         Route,
	"route_config.virtual_hosts.#.routes.#.route.cluster":                           Cluster,
	"route_config.virtual_hosts.#.routes.#.route.weighted_clusters.clusters.#.name": Cluster,
	"route_config.virtual_hosts.#.request_mirror_policies.#.cluster":                Cluster,
	"route_config.request_mirror_policies.#.cluster":                                Cluster,
}

var RouteUpstreams = map[string]GTypes{
	"virtual_hosts.#.routes.#.route.cluster":                           Cluster,
	"virtual_hosts.#.routes.#.route.weighted_clusters.clusters.#.name": Cluster,
	"virtual_hosts.#.request_mirror_policies.#.cluster":                Cluster,
	"request_mirror_policies.#.cluster":                                Cluster,
}

var VirtualHostUpstreams = map[string]GTypes{
	"routes.#.route.cluster":                           Cluster,
	"routes.#.route.weighted_clusters.clusters.#.name": Cluster,
	"request_mirror_policies.#.cluster":                Cluster,
}

var FluentdAccessLogUpstreams = map[string]GTypes{
	"cluster": Cluster,
}

var DownstreamTlsContextUpstreams = map[string]GTypes{
	"common_tls_context.tls_certificate_sds_secret_configs.#.name": TlsCertificate,
	"common_tls_context.validation_context_sds_secret_config.name": CertificateValidationContext,
}

var UpstreamTlsContextUpstreams = map[string]GTypes{
	"common_tls_context.tls_certificate_sds_secret_configs.#.name": TlsCertificate,
	"common_tls_context.validation_context_sds_secret_config.name": CertificateValidationContext,
}
