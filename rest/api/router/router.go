package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/sefaphlvn/bigbang/rest/api/middleware"
	"github.com/sefaphlvn/bigbang/rest/handlers"
)

func InitRouter(h *handlers.Handler, logger *logrus.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()

	e.HandleMethodNotAllowed = true
	e.ForwardedByClientIP = true

	e.Use(middleware.CORS())
	e.Use(middleware.PathCheck())
	e.Use(middleware.GinLog(logger), gin.Recovery())

	e.POST("/logout", middleware.Authentication(), h.Auth.Logout())
	e.POST("/refresh", middleware.Refresh(), h.Auth.Refresh())

	api := e.Group("/api")
	v3 := api.Group("/v3")
	v3.Use(middleware.Authentication())

	apiAuth := e.Group("/auth")
	apiSettings := v3.Group("/setting")
	apiCustom := v3.Group("/custom")
	apiExtension := v3.Group("/eo")
	apiResource := v3.Group("/xds")
	apiDependency := v3.Group("/dependency")
	apiScenario := v3.Group("/scenario")
	apiBridge := v3.Group("/bridge")

	initAuthRoutes(apiAuth, h)
	initSettingRoutes(apiSettings, h)
	initCustomRoutes(apiCustom, h)
	initExtensionRoutes(apiExtension, h)
	initScenarioRoutes(apiScenario, h)
	initResourceRoutes(apiResource, h)
	initDependencyRoutes(apiDependency, h)
	initBridgeRoutes(apiBridge, h)

	logRoutes(e)
	return e
}

func initAuthRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	enableDemo := h.Auth.Context.Config.BigbangEnableDemo == "true"
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"POST", "/login", h.Auth.Login()},
	}

	if enableDemo {
		routes = append(routes, struct {
			method  string
			path    string
			handler gin.HandlerFunc
		}{
			"POST", "/demo/:email", h.Auth.DemoAccount,
		})
	}

	initRoutes(rg, routes)
}

func initBridgeRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/stats/snapshot-keys", h.GetSnapshotKeys},
		{"GET", "/stats/:name", h.GetSnapshotResources},
		{"POST", "/poke/:name", h.GetSnapshotResources},
		{"GET", "/snapshot_details", h.GetSnapshotDetails},
	}

	initRoutes(rg, routes)
}

func initSettingRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	rg.Use(middleware.InitSettingMiddleware())

	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/user_list", h.Auth.ListUsers},
		{"GET", "/user/:user_id", h.Auth.GetUser},
		{"PUT", "/user/:user_id", h.Auth.SetUpdateUser},

		{"GET", "/group_list", h.Auth.ListGroups},
		{"GET", "/group/:group_id", h.Auth.GetGroup},
		{"PUT", "/group/:group_id", h.Auth.SetUpdateGroup},

		{"GET", "/project_list", h.Auth.ListProjects},
		{"GET", "/project/:project_id", h.Auth.GetProject},
		{"PUT", "/project/:project_id", h.Auth.SetUpdateProject},

		{"GET", "/permissions/:kind/:type/:id", h.Auth.GetPermissions},
	}

	initRoutes(rg, routes)
}

func initCustomRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/resource_list", h.GetCustomResourceList},
		{"GET", "/http_filter_list", h.GetCustomHTTPFilterList},

		{"GET", "/count/all", h.GetResourceCounts},
		{"GET", "/count/filters", h.GetFilterCounts},
	}

	initRoutes(rg, routes)
}

func initScenarioRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/scenario_list", h.GetScenarios},
		{"GET", "/scenario", h.GetScenario},
		{"POST", "/scenario", h.SetScenario},
	}

	initRoutes(rg, routes)
}

func initDependencyRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/:name", h.GetResourceDependencies},
	}

	initRoutes(rg, routes)
}

func initExtensionRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/:collection/extensions/:type", h.GetExtensions},
		{"POST", "/:collection/extensions/:type", h.SetExtension},
		{"GET", "/:collection/extensions/:type/:name", h.GetOtherExtension},
		{"PUT", "/:collection/extensions/:type/:name", h.UpdateOtherExtensions},
		{"DELETE", "/:collection/extensions/:type/:name", h.DelExtension},

		{"GET", "/:collection/:type/:canonical_name", h.ListExtensions},
		{"POST", "/:collection/:type/:canonical_name", h.SetExtension},
		{"GET", "/:collection/:type/:canonical_name/:name", h.GetExtension},
		{"PUT", "/:collection/:type/:canonical_name/:name", h.UpdateExtension},
		{"DELETE", "/:collection/:type/:canonical_name/:name", h.DelExtension},
	}

	initRoutes(rg, routes)
}

func initResourceRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/:collection", h.ListResource},
		{"POST", "/:collection", h.SetResource},
		{"GET", "/:collection/:name", h.GetResource},
		{"PUT", "/:collection/:name", h.UpdateResource},
		{"DELETE", "/:collection/:name", h.DelResource},
	}

	initRoutes(rg, routes)
}

func initRoutes(rg *gin.RouterGroup, routes []struct {
	method  string
	path    string
	handler gin.HandlerFunc
},
) {
	for _, route := range routes {
		switch route.method {
		case "GET":
			rg.GET(route.path, route.handler)
		case "POST":
			rg.POST(route.path, route.handler)
		case "PUT":
			rg.PUT(route.path, route.handler)
		case "DELETE":
			rg.DELETE(route.path, route.handler)
		}
	}
}

func logRoutes(r *gin.Engine) {
	log.Println("Registered Routes:")
	for _, route := range r.Routes() {
		log.Printf("Method: %s, Path: %s, Handler: %s\n", route.Method, route.Path, route.Handler)
	}
}
