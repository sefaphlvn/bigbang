package models

type ArrayPath struct {
	ParentPath string
	IndexPath  string
}

type TypedConfigPath struct {
	ArrayPaths       []ArrayPath
	PathTemplate     string
	Kind             string
	IsPerTypedConfig bool
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

var RouteTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths:       []ArrayPath{},
		PathTemplate:     "typed_per_filter_config",
		Kind:             "route",
		IsPerTypedConfig: true,
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "virtual_hosts", IndexPath: "virtual_hosts.%d"},
		},
		PathTemplate:     "virtual_hosts.%d.typed_per_filter_config",
		Kind:             "route",
		IsPerTypedConfig: true,
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "virtual_hosts", IndexPath: "virtual_hosts.%d"},
			{ParentPath: "virtual_hosts.%d.routes", IndexPath: "routes.%d"},
		},
		PathTemplate:     "virtual_hosts.%d.routes.%d.typed_per_filter_config",
		Kind:             "route",
		IsPerTypedConfig: true,
	},
}

var VirtualHostTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths:       []ArrayPath{},
		PathTemplate:     "typed_per_filter_config",
		Kind:             "virtual_host",
		IsPerTypedConfig: true,
	},
}
