package handlers

import (
	"net/http"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/api/auth"
	"github.com/sefaphlvn/bigbang/rest/crud/custom"
	"github.com/sefaphlvn/bigbang/rest/crud/extension"
	"github.com/sefaphlvn/bigbang/rest/crud/xds"

	"github.com/gin-gonic/gin"
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

func (h *Handler) handleRequest(c *gin.Context, dbFunc DBFunc) {
	groups, _ := c.Get("groups")
	isAdmin, _ := c.Get("isAdmin")
	userGroup, ok := groups.([]string)
	if !ok {
		userGroup = []string{}
	}

	userIsAdmin, ok := isAdmin.(bool)
	if !ok {
		userIsAdmin = false
	}

	resourceDetails := models.ResourceDetails{
		CanonicalName: c.Param("canonical_name"),
		Category:      c.Query("category"),
		Name:          c.Param("name"),
		Collection:    c.Query("collection"),
		User: models.UserDetails{
			Groups:  userGroup,
			IsAdmin: userIsAdmin,
		},
	}

	if version := c.Param("version"); version != "" {
		resourceDetails.Version = version
	} else if version := c.Query("version"); version != "" {
		resourceDetails.Version = version
	}

	if ltype := c.Param("type"); ltype != "" {
		resourceDetails.Type = models.KnownTYPES(ltype)
	} else if ltype := c.Query("type"); ltype != "" {
		resourceDetails.Type = models.KnownTYPES(ltype)
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
