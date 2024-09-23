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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type GTypeMapping struct {
	Collection string
	URL        string
	PrettyName string
	Message    proto.Message
}

var gTypeMappings = map[GTypes]GTypeMapping{
	BootStrap:                    {PrettyName: "Bootstrap", Collection: "bootstrap", URL: "/resource/bootstrap/", Message: &bootstrap.Bootstrap{}},
	HTTPConnectionManager:        {PrettyName: "Http Connection Manager", Collection: "extensions", URL: "/filters/network/hcm/", Message: &hcm.HttpConnectionManager{}},
	Router:                       {PrettyName: "Router", Collection: "extensions", URL: "/filters/http/http_router/", Message: &router.Router{}},
	Cluster:                      {PrettyName: "Cluster", Collection: "clusters", URL: "/resource/cluster/", Message: &cluster.Cluster{}},
	Listener:                     {PrettyName: "Listener", Collection: "listeners", URL: "/resource/listener/", Message: &listener.Listener{}},
	Endpoint:                     {PrettyName: "Endpoint", Collection: "endpoints", URL: "/resource/endpoint", Message: &endpoint.ClusterLoadAssignment{}},
	Route:                        {PrettyName: "Route", Collection: "routes", URL: "/resource/route", Message: &route.RouteConfiguration{}},
	VirtualHost:                  {PrettyName: "Virtual Host", Collection: "virtual_host", URL: "/resource/virtual_host", Message: &route.VirtualHost{}},
	TcpProxy:                     {PrettyName: "Tcp Proxy", Collection: "extensions", URL: "/filters/network/tcp_proxy/", Message: &tcp.TcpProxy{}},
	FluentdAccessLog:             {PrettyName: "Access Log(Fluentd)", Collection: "others", URL: "/others/access_log/", Message: &al_fluentd.FluentdAccessLogConfig{}},
	FileAccessLog:                {PrettyName: "Access Log(File)", Collection: "others", URL: "/others/access_log/", Message: &al_file.FileAccessLog{}},
	StdoutAccessLog:              {PrettyName: "Access Log(StdOut)", Collection: "others", URL: "/others/access_log/", Message: &al_stream.StdoutAccessLog{}},
	StdErrAccessLog:              {PrettyName: "Access Log(StdErr)", Collection: "others", URL: "/others/access_log/", Message: &al_stream.StderrAccessLog{}},
	DownstreamTlsContext:         {PrettyName: "Downstream TLS", Collection: "secrets", URL: "/resource/secret/", Message: &tls.DownstreamTlsContext{}},
	UpstreamTlsContext:           {PrettyName: "Upstream TLS", Collection: "secrets", URL: "/resource/secret/", Message: &tls.UpstreamTlsContext{}},
	TlsCertificate:               {PrettyName: "TLS Certificate", Collection: "secrets", URL: "/resource/secret/", Message: &tls.TlsCertificate{}},
	CertificateValidationContext: {PrettyName: "Certificate Validation", Collection: "secrets", URL: "/resource/secret/", Message: &tls.CertificateValidationContext{}},
	HealthCheckEventFileSink:     {PrettyName: "Health Check Event File Sink", Collection: "others", URL: "/others/hcefs/", Message: &hcefs.HealthCheckEventFileSink{}},
	BasicAuth:                    {PrettyName: "Basic Auth", Collection: "extensions", URL: "/filters/http/basic_auth/", Message: &basic_auth.BasicAuth{}},
}

func (gt GTypes) String() string {
	return string(gt)
}

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

func (gt GTypes) ProtoMessage() proto.Message {
	if str, exists := gTypeMappings[gt]; exists {
		return str.Message
	}
	return &anypb.Any{}
}
