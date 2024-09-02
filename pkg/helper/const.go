package helper

import "os"

var Collections = []string{
	"clusters",
	"listeners",
	"endpoints",
	"routes",
	"lb_endpoints",
	"extensions",
	"vhds",
	"others",
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
	"/api/v3/xds/vhds",
	"/api/v3/xds/vhds/:name",
	"/api/v3/xds/clusters",
	"/api/v3/xds/clusters/:name",
	"/api/v3/xds/hcm",
	"/api/v3/xds/hcm/:name",
	"/api/v3/xds/endpoints",
	"/api/v3/xds/endpoints/:name",
	"/api/v3/eo/:collection/:type",
	"/api/v3/eo/:collection/:type/:canonical_name",
	"/api/v3/eo/:collection/:type/:canonical_name/:name",
	"/api/v3/eo/:collection/others/:type",
	"/api/v3/eo/:collection/others/:type/:name",
	"/api/v3/custom/resource_list",
	"/api/v3/dependency/:name",
}

var SECRET_KEY string = os.Getenv("secret")
