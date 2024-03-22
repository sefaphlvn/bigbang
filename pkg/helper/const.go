package helper

var Collections = []string{
	"clusters",
	"listeners",
	"endpoints",
	"routes",
	"lb_endpoints",
	"extensions",
	"vhds",
}

var AllowedEndpoints = []string{
	"/logout",
	"/auth/signup",
	"/auth/login",
	"/refresh",
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
	"/api/v3/extensions/:type",
	"/api/v3/extensions/:type/:canonical_name",
	"/api/v3/extensions/:type/:canonical_name/:name",
	"/api/v3/custom/resource_list",
}
