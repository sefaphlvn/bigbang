package models

type ArrayPath struct {
	ParentPath string // Dizinin bulunduğu üst seviye path
	IndexPath  string // Dizinin kendisi, %d yer tutucusu ile gösterilir
}

type TypedConfigPath struct {
	ArrayPaths   []ArrayPath // Dizileri tanımlayan tüm seviyeler
	PathTemplate string      // Nihai path için template
	Kind         string      // Typed config türü
}

var ConfigGetters = map[GTypes][]TypedConfigPath{
	Listener:              ListenerTypedConfigPaths,
	HTTPConnectionManager: GeneralAccessLogTypedConfigPaths,
	TcpProxy:              GeneralAccessLogTypedConfigPaths,
	BootStrap:             BootstrapTypedConfigPaths,
	Cluster:               ClusterTypedConfigPaths,
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
		ArrayPaths:   []ArrayPath{}, // Array olmadığı için boş bırakıldı
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
