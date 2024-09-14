package models

/*
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
*/
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
			"rds.route_config_name": Route,
		}
	case Route:
		return map[string]GTypes{
			"virtual_hosts.#.routes.#.route.cluster":                           Cluster,
			"virtual_hosts.#.routes.#.route.weighted_clusters.clusters.#.name": Cluster,
			"virtual_hosts.#.request_mirror_policies.#.cluster":                Cluster,
			"request_mirror_policies.#.cluster":                                Cluster,
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
