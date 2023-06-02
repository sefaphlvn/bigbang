package helper

var Collections = []string{
	"clusters",
	"listeners",
	"endpoints",
	"routes",
	"lb_endpoints",
	"extensions",
}

var AllowedEndpoints = []string{
	"/api/v3/listeners",
	"/api/v3/listeners/:name",
	"/api/v3/routes",
	"/api/v3/routes/:name",
	"/api/v3/clusters",
	"/api/v3/clusters/:name",
	"/api/v3/hcm",
	"/api/v3/hcm/:name",
	"/api/v3/endpoints",
	"/api/v3/endpoints/:name",
	"/api/v3/lb_endpoints",
	"/api/v3/lb_endpoints/:name",
	"/api/v3/extensions/:type",
	"/api/v3/extensions/:type/:subtype",
	"/api/v3/extensions/:type/:subtype/:name",
}
