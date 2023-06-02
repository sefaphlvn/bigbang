package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/restApi/crud/extension"
	"github.com/sefaphlvn/bigbang/restApi/crud/xds"
	"github.com/sefaphlvn/bigbang/restApi/models"
)

type DBFunc func(resource models.DBResourceClass, resourceType models.ResourceDetails) (interface{}, error)

type Handler struct {
	XDS       *xds.DBHandler
	Extension *extension.DBHandler
}

func NewHandler(XDS *xds.DBHandler, extension *extension.DBHandler) *Handler {
	return &Handler{
		XDS:       XDS,
		Extension: extension,
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
