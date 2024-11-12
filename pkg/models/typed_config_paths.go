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

var (
	access_log              = "access_log.%d"
	access_log_typed_config = "access_log.%d.typed_config"
	routes                  = "routes.%d"
)

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
			{ParentPath: "access_log", IndexPath: access_log},
		},
		PathTemplate: access_log_typed_config,
		Kind:         "access_log",
	},
}

var GeneralAccessLogTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "access_log", IndexPath: access_log},
		},
		PathTemplate: access_log_typed_config,
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
	{
		ArrayPaths:       []ArrayPath{},
		PathTemplate:     "typed_extension_protocol_options",
		Kind:             "http_protocol",
		IsPerTypedConfig: true,
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
			{ParentPath: "virtual_hosts.%d.routes", IndexPath: routes},
		},
		PathTemplate:     "virtual_hosts.%d.routes.%d.typed_per_filter_config",
		Kind:             "route",
		IsPerTypedConfig: true,
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "virtual_hosts", IndexPath: "virtual_hosts.%d"},
			{ParentPath: "virtual_hosts.%d.routes", IndexPath: routes},
		},
		PathTemplate: "virtual_hosts.%d.routes.%d.match.path_match_policy.typed_config",
		Kind:         "route",
	},
}

var VirtualHostTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths:       []ArrayPath{},
		PathTemplate:     "typed_per_filter_config",
		Kind:             "virtual_host",
		IsPerTypedConfig: true,
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "routes", IndexPath: routes},
		},
		PathTemplate:     "routes.%d.typed_per_filter_config",
		Kind:             "virtual_host",
		IsPerTypedConfig: true,
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "routes", IndexPath: "routes.%d"},
		},
		PathTemplate: "routes.%d.match.path_match_policy.typed_config",
		Kind:         "virtual_host",
	},
}

var HttpConnectionManagerTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "access_log", IndexPath: access_log},
		},
		PathTemplate: access_log_typed_config,
		Kind:         "access_log",
	},
	{
		ArrayPaths:       []ArrayPath{},
		PathTemplate:     "route_config.typed_per_filter_config",
		Kind:             "hcm",
		IsPerTypedConfig: true,
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "route_config.virtual_hosts", IndexPath: "route_config.virtual_hosts.%d"},
		},
		PathTemplate:     "route_config.virtual_hosts.%d.typed_per_filter_config",
		Kind:             "hcm",
		IsPerTypedConfig: true,
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "route_config.virtual_hosts", IndexPath: "route_config.virtual_hosts.%d"},
			{ParentPath: "route_config.virtual_hosts.%d.routes", IndexPath: routes},
		},
		PathTemplate:     "route_config.virtual_hosts.%d.routes.%d.typed_per_filter_config",
		Kind:             "hcm",
		IsPerTypedConfig: true,
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "route_config.virtual_hosts", IndexPath: "route_config.virtual_hosts.%d"},
			{ParentPath: "route_config.virtual_hosts.%d.routes", IndexPath: routes},
		},
		PathTemplate: "route_config.virtual_hosts.%d.routes.%d.match.path_match_policy.typed_config",
		Kind:         "hcm",
	},
}

var CompressorTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths:   []ArrayPath{},
		PathTemplate: "compressor_library.typed_config",
		Kind:         "compressor_library",
	},
}

var RBACTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.not_rule.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.and_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.and_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.or_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.or_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.not_rule.and_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.not_rule.and_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.not_rule.or_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.not_rule.or_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.and_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.and_rules.rules.%d.and_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.and_rules.rules.%d.and_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.or_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.or_rules.rules.%d.or_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.or_rules.rules.%d.or_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.not_rule.and_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.not_rule.and_rules.rules.%d.or_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.not_rule.and_rules.rules.%d.or_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.and_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.and_rules.rules.%d.not_rule.and_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.and_rules.rules.%d.not_rule.and_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.or_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.or_rules.rules.%d.not_rule.uri_template.typed_config", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.or_rules.rules.%d.not_rule.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rules.policies.*.permissions", IndexPath: "rules.policies.*.permissions.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.not_rule.or_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rules.policies.*.permissions.%d.not_rule.or_rules.rules.%d.not_rule.uri_template.typed_config", IndexPath: "rules.%d"},
		},
		PathTemplate: "rules.policies.*.permissions.%d.not_rule.or_rules.rules.%d.not_rule.uri_template.typed_config",
		Kind:         "uri_template",
	},
}

var RBACPerRouteTypedConfigPaths = []TypedConfigPath{
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.not_rule.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.and_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.and_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.or_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.or_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.not_rule.and_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.not_rule.and_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.not_rule.or_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.not_rule.or_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.and_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.and_rules.rules.%d.and_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.and_rules.rules.%d.and_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.or_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.or_rules.rules.%d.or_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.or_rules.rules.%d.or_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.not_rule.and_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.not_rule.and_rules.rules.%d.or_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.not_rule.and_rules.rules.%d.or_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.and_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.and_rules.rules.%d.not_rule.and_rules.rules", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.and_rules.rules.%d.not_rule.and_rules.rules.%d.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.or_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.or_rules.rules.%d.not_rule.uri_template.typed_config", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.or_rules.rules.%d.not_rule.uri_template.typed_config",
		Kind:         "uri_template",
	},
	{
		ArrayPaths: []ArrayPath{
			{ParentPath: "rbac.rules.policies.*.permissions", IndexPath: "rbac.rules.policies.*.permissions.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.not_rule.or_rules.rules", IndexPath: "rules.%d"},
			{ParentPath: "rbac.rules.policies.*.permissions.%d.not_rule.or_rules.rules.%d.not_rule.uri_template.typed_config", IndexPath: "rules.%d"},
		},
		PathTemplate: "rbac.rules.policies.*.permissions.%d.not_rule.or_rules.rules.%d.not_rule.uri_template.typed_config",
		Kind:         "uri_template",
	},
}
