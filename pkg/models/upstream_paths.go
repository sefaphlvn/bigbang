package models

func (gt GTypes) GetUpstreamPaths() map[string]GTypes {
	switch gt {
	case Cluster:
		return map[string]GTypes{
			"eds_cluster_config.service_name": Endpoint,
		}
	case TcpProxy:
		return map[string]GTypes{
			"cluster":                           Cluster,
			"weighted_clusters.clusters.#.name": Cluster,
		}
	case HTTPConnectionManager:
		return map[string]GTypes{
			"rds.route_config_name":                                                         Route,
			"route_config.virtual_hosts.#.routes.#.route.cluster":                           Cluster,
			"route_config.virtual_hosts.#.routes.#.route.weighted_clusters.clusters.#.name": Cluster,
			"route_config.virtual_hosts.#.request_mirror_policies.#.cluster":                Cluster,
			"route_config.request_mirror_policies.#.cluster":                                Cluster,
		}
	case Route:
		return map[string]GTypes{
			"virtual_hosts.#.routes.#.route.cluster":                           Cluster,
			"virtual_hosts.#.routes.#.route.weighted_clusters.clusters.#.name": Cluster,
			"virtual_hosts.#.request_mirror_policies.#.cluster":                Cluster,
			"request_mirror_policies.#.cluster":                                Cluster,
		}
	case VirtualHost:
		return map[string]GTypes{
			"routes.#.route.cluster":                           Cluster,
			"routes.#.route.weighted_clusters.clusters.#.name": Cluster,
			"request_mirror_policies.#.cluster":                Cluster,
		}
	case FluentdAccessLog:
		return map[string]GTypes{
			"cluster": Cluster,
		}
	case DownstreamTlsContext:
		return map[string]GTypes{
			"common_tls_context.tls_certificate_sds_secret_configs.#.name": TlsCertificate,
			"common_tls_context.validation_context_sds_secret_config.name": CertificateValidationContext,
		}
	case UpstreamTlsContext:
		return map[string]GTypes{
			"common_tls_context.tls_certificate_sds_secret_configs.#.name": TlsCertificate,
			"common_tls_context.validation_context_sds_secret_config.name": CertificateValidationContext,
		}
	default:
		return nil
	}
}
