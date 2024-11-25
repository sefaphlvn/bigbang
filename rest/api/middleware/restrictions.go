package middleware

var Collections = []string{
	"clusters",
	"listeners",
	"endpoints",
	"routes",
	"lb_endpoints",
	"extensions",
	"filters",
	"vhds",
}

var AllowedEndpoints = []string{
	"/logout",
	"/refresh",
	"/auth/login",
	"/api/v3/setting/user_list",
	"/api/v3/setting/user/:user_id",
	"/api/v3/setting/group_list",
	"/api/v3/setting/group/:group_id",
	"/api/v3/setting/project_list",
	"/api/v3/setting/project/:project_id",
	"/api/v3/setting/permissions/:kind/:type/:id",
	"/api/v3/xds/secrets",
	"/api/v3/xds/secrets/:name",
	"/api/v3/xds/bootstrap",
	"/api/v3/xds/bootstrap/:name",
	"/api/v3/xds/listeners",
	"/api/v3/xds/listeners/:name",
	"/api/v3/xds/routes",
	"/api/v3/xds/routes/:name",
	"/api/v3/xds/tls",
	"/api/v3/xds/tls/:name",
	"/api/v3/xds/virtual_hosts",
	"/api/v3/xds/virtual_hosts/:name",
	"/api/v3/xds/clusters",
	"/api/v3/xds/clusters/:name",
	"/api/v3/xds/hcm",
	"/api/v3/xds/hcm/:name",
	"/api/v3/xds/endpoints",
	"/api/v3/xds/endpoints/:name",
	"/api/v3/eo/:collection/:type",
	"/api/v3/eo/:collection/:type/:canonical_name",
	"/api/v3/eo/:collection/:type/:canonical_name/:name",
	"/api/v3/eo/:collection/extensions/:type",
	"/api/v3/eo/:collection/extensions/:type/:name",
	"/api/v3/custom/resource_list",
	"/api/v3/custom/http_filter_list",
	"/api/v3/dependency/:name",
	"/api/v3/bridge/stats/:name",
	"/api/v3/bridge/poke/:name",
	"/api/v3/bridge/errors",
}
