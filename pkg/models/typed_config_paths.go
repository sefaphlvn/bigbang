package models

type ArrayPath struct {
	ParentPath string
	IndexPath  string
}

type TypedConfigPath struct {
	ArrayPaths   []ArrayPath
	PathTemplate string
	Kind         string
}

var BootstrapTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "admin.access_log", IndexPath: "admin.access_log.%d"},
		},
		PathTemplate: "admin.access_log.%d.typed_config",
		Kind:         "access_log",
	},
}

var ListenerTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "filter_chains", IndexPath: "filter_chains.%d"},
		},
		PathTemplate: "filter_chains.%d.transport_socket.typed_config",
		Kind:         "downstream_tls",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "access_log", IndexPath: "access_log.%d"},
		},
		PathTemplate: "access_log.%d.typed_config",
		Kind:         "access_log",
	},
}

var GeneralAccessLogTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "access_log", IndexPath: "access_log.%d"},
		},
		PathTemplate: "access_log.%d.typed_config",
		Kind:         "access_log",
	},
}

var ClusterTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths:   []ArrayPath{},
		PathTemplate: "transport_socket.typed_config",
		Kind:         "upstream_tls",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "health_checks", IndexPath: "health_checks.%d"},
			{ParentPath: "health_checks.%d.event_logger", IndexPath: "event_logger.%d"},
		},
		PathTemplate: "health_checks.%d.event_logger.%d.typed_config",
		Kind:         "hcefs",
	},
}
