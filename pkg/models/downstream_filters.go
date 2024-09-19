package models

import "github.com/sefaphlvn/bigbang/pkg/filters"

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
	case VirtualHost:
		return filters.VirtualHostDownstreamFilters(name)
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
	case UpstreamTlsContext:
		return []filters.MongoFilters{
			filters.UpstreamTlsDownstreamFilters(name),
		}
	case HealthCheckEventFileSink:
		return filters.HCEFSDownstreamFilters(name)
	default:
		return nil
	}
}
