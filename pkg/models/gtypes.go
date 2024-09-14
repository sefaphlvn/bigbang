package models

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
	HealthCheckEventFileSink     GTypes = "envoy.extensions.health_check.event_sinks.file.v3.HealthCheckEventFileSink"
)
