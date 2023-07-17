package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/restServer/handlers"
	"github.com/sefaphlvn/bigbang/restServer/middleware"
)

func InitRouter(h *handlers.Handler) *gin.Engine {
	r := gin.New()
	r.Use(middleware.CORS())
	r.Use(middleware.PathCheck())
	r.Use(gin.Logger())

	r.POST("/logout", middleware.Authentication(), h.Auth.Logout())
	r.POST("/refresh", middleware.Refresh(), h.Auth.Refresh())

	apiCustom := r.Group("/api/v3/custom")
	apiExtension := r.Group("/api/v3/extensions")
	apiResource := r.Group("/api/v3")
	apiAuth := r.Group("/auth")

	apiCustom.Use(middleware.Authentication())
	apiExtension.Use(middleware.Authentication())
	apiResource.Use(middleware.Authentication())

	initAuthRoutes(apiAuth, h)
	initCustomRoutes(apiCustom, h)
	initExtensionRoutes(apiExtension, h)
	initResourceRoutes(apiResource, h)

	return r
}

func initAuthRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"POST", "/signup", h.Auth.SignUp()},
		{"POST", "/login", h.Auth.Login()},
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

func initExtensionRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/:type", h.GetExtensions},
		{"GET", "/:type/:subtype", h.ListExtensions},
		{"POST", "/:type/:subtype", h.SetExtension},
		{"GET", "/:type/:subtype/:name", h.GetExtension},
		{"PUT", "/:type/:subtype/:name", h.UpdateExtension},
		{"DELETE", "/:type/:subtype/:name", h.GetExtension},
	}

	initRoutes(rg, routes)
}

func initResourceRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/:type", h.ListResource},
		{"POST", "/:type", h.SetResource},
		{"GET", "/:type/:name", h.GetResource},
		{"PUT", "/:type/:name", h.UpdateResource},
		{"DELETE", "/:type/:name", h.DelResource},
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
