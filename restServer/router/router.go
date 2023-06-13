package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/restServer/handlers"
	"github.com/sefaphlvn/bigbang/restServer/helper"
)

func InitRouter(h *handlers.Handler) *gin.Engine {
	r := gin.New()
	r.Use(CORS())
	r.Use(PathCheck())
	r.Use(gin.Logger())

	apiCustom := r.Group("/api/v3/custom")
	apiExtension := r.Group("/api/v3/extensions")
	apiResource := r.Group("/api/v3")

	initCustomRoutes(apiCustom, h)
	initExtensionRoutes(apiExtension, h)
	initResourceRoutes(apiResource, h)

	return r
}

func initCustomRoutes(rg *gin.RouterGroup, h *handlers.Handler) {
	routes := []struct {
		method  string
		path    string
		handler gin.HandlerFunc
	}{
		{"GET", "/filter_chain_filters", h.GetFilterChainFilters},
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
		{"PUT", "/:type/:subtype/:name", h.SetExtension},
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

// CORS Allow
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Path Allow
func PathCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		pathParts := strings.Split(path, "/")
		for _, allowedPath := range helper.AllowedEndpoints {
			allowedParts := strings.Split(allowedPath, "/")
			if len(pathParts) != len(allowedParts) {
				continue
			}

			matched := true
			for i := range pathParts {
				if allowedParts[i] != pathParts[i] && !strings.HasPrefix(allowedParts[i], ":") {
					matched = false
					break
				}
			}

			if matched {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid path"})
	}
}
