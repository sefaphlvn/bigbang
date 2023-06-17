package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/restServer/auth"
	"github.com/sefaphlvn/bigbang/restServer/crud/custom"
	"github.com/sefaphlvn/bigbang/restServer/crud/extension"
	"github.com/sefaphlvn/bigbang/restServer/crud/xds"
	"github.com/sefaphlvn/bigbang/restServer/models"
)

type DBFunc func(resource models.DBResourceClass, resourceType models.ResourceDetails) (interface{}, error)

type Handler struct {
	XDS       *xds.DBHandler
	Extension *extension.DBHandler
	Custom    *custom.DBHandler
	Auth      *auth.DBHandler
}

func NewHandler(XDS *xds.DBHandler, extension *extension.DBHandler, custom *custom.DBHandler, auth *auth.DBHandler) *Handler {
	return &Handler{
		XDS:       XDS,
		Extension: extension,
		Custom:    custom,
		Auth:      auth,
	}
}

func decodeResource(c *gin.Context) (models.DBResourceClass, error) {
	var body models.DBResource
	if c.Request.Method != "GET" && c.Request.Method != "DELETE" {
		err := c.BindJSON(&body)
		if err != nil {
			return nil, err
		}
	}
	return &body, nil
}

func (h *Handler) handleResource(c *gin.Context, dbFunc DBFunc) {
	resourceDetails := models.ResourceDetails{
		Type:    c.Param("type"),
		SubType: c.Param("subtype"),
		Name:    c.Param("name"),
		Version: c.Param("version"),
	}
	resource, err := decodeResource(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	response, err := dbFunc(resource, resourceDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}
