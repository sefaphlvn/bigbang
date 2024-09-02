package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/rest/api/middleware"
	"github.com/sefaphlvn/bigbang/rest/handlers"
	"github.com/sirupsen/logrus"
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

	apiAuth := e.Group("/auth")
	apiSettings := e.Group("/api/v3/setting")
	apiCustom := e.Group("/api/v3/custom")
	apiExtension := e.Group("/api/v3/eo")
	apiResource := e.Group("/api/v3/xds")
	apiDependency := e.Group("/api/v3/dependency")

	apiSettings.Use(middleware.Authentication())
	apiCustom.Use(middleware.Authentication())
	apiExtension.Use(middleware.Authentication())
	apiResource.Use(middleware.Authentication())
	//apiDependency.Use(middleware.Authentication())

	initAuthRoutes(apiAuth, h)
	initSettingRoutes(apiSettings, h)
	initCustomRoutes(apiCustom, h)
	initExtensionRoutes(apiExtension, h)
	initResourceRoutes(apiResource, h)
	initDependencyRoutes(apiDependency, h)

	return e
}

func initAuthRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"POST", "/login", h.Auth.Login()},
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
		{"GET", "/:collection/others/:type", h.GetExtensions},
		{"POST", "/:collection/others/:type", h.SetExtension},
		{"GET", "/:collection/others/:type/:name", h.GetOtherExtension},
		{"PUT", "/:collection/others/:type/:name", h.UpdateOtherExtensions},
		{"DELETE", "/:collection/others/:type/:canonical_name/:name", h.GetExtension},

		{"GET", "/:collection/:type", h.GetExtensions},
		{"GET", "/:collection/:type/:canonical_name", h.ListExtensions},
		{"POST", "/:collection/:type/:canonical_name", h.SetExtension},
		{"GET", "/:collection/:type/:canonical_name/:name", h.GetExtension},
		{"PUT", "/:collection/:type/:canonical_name/:name", h.UpdateExtension},
		{"DELETE", "/:collection/:type/:canonical_name/:name", h.GetExtension},
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
}) {
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
