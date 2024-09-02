package models

import (
	accessLog "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/stream/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"google.golang.org/protobuf/proto"
)

type TypedConfigPath struct {
	JsonPath     string
	PathTemplate string
	Kind         string
}

var BootstrapTypedConfigPaths = []TypedConfigPath{
	{JsonPath: "admin.access_log", PathTemplate: "admin.access_log.%d.typed_config", Kind: "access_log"},
}

var ListenerTypedConfigPaths = []TypedConfigPath{
	{JsonPath: "filter_chains", PathTemplate: "filter_chains.%d.transport_socket.typed_config", Kind: "downstream_tls"},
	{JsonPath: "access_log", PathTemplate: "access_log.%d.typed_config", Kind: "access_log"},
}

var GeneralAccessLogTypedConfigPaths = []TypedConfigPath{
	{JsonPath: "access_log", PathTemplate: "access_log.%d.typed_config", Kind: "access_log"},
}

// TypedConfigMap is a map of string to proto.Message type
var TypedConfigMap = map[string]proto.Message{
	"envoy.extensions.access_loggers.stream.v3.StdoutAccessLog":      &accessLog.StdoutAccessLog{},
	"envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext": &tls.DownstreamTlsContext{},
}
