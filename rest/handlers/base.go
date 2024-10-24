package handlers

import (
	"errors"
	"net/http"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/api/auth"
	"github.com/sefaphlvn/bigbang/rest/bridge"
	"github.com/sefaphlvn/bigbang/rest/crud/custom"
	"github.com/sefaphlvn/bigbang/rest/crud/extension"
	"github.com/sefaphlvn/bigbang/rest/crud/xds"
	"github.com/sefaphlvn/bigbang/rest/dependency"

	"github.com/gin-gonic/gin"
)

type DBFunc func(resource models.DBResourceClass, requestDetails models.RequestDetails) (interface{}, error)
type DepFunc func(requestDetails models.RequestDetails) (*dependency.DependencyGraph, error)

type Handler struct {
	XDS        *xds.AppHandler
	Extension  *extension.AppHandler
	Custom     *custom.AppHandler
	Auth       *auth.AppHandler
	dependency *dependency.AppHandler
	Bridge     *bridge.AppHandler
}

func NewHandler(XDS *xds.AppHandler, extension *extension.AppHandler, custom *custom.AppHandler, auth *auth.AppHandler, dependency *dependency.AppHandler, stats *bridge.AppHandler) *Handler {
	return &Handler{
		XDS:        XDS,
		Extension:  extension,
		Custom:     custom,
		Auth:       auth,
		dependency: dependency,
		Bridge:     stats,
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

// This function handles a request in the Handler struct.
// It retrieves the necessary data from the context, including the groups and isOwner parameters.
// It then sets the requestDetails struct with the given parameters and decodes the resource.
// It then calls the dbFunc with the resource and requestDetails, and stores the response in the response variable.
// Finally, it returns the response as a JSON object with the status OK.
func (h *Handler) handleRequest(c *gin.Context, dbFunc DBFunc) {
	userDetails, _ := GetUserDetails(c)

	requestDetails := models.RequestDetails{
		CanonicalName: getParamOrQuery(c, "canonical_name"),
		GType:         models.GTypes(c.Query("gtype")),
		Category:      c.Query("category"),
		Name:          c.Param("name"),
		Collection:    getParamOrQuery(c, "collection"),
		SaveOrPublish: c.Query("save_or_publish"),
		User:          userDetails,
		Project:       c.Query("project"),
		Metadata:      extractMetadata(c),
		Version:       getOptionalParam(c, "version"),
		Type:          models.KnownTYPES(getOptionalParam(c, "type")),
	}

	if err := checkRole(c, userDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	resource, err := decodeResource(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, err := dbFunc(resource, requestDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": response})
		return
	}

	c.JSON(http.StatusOK, response)
}

func GetUserDetails(c *gin.Context) (models.UserDetails, error) {
	groups, _ := c.Get("groups")
	isOwner, _ := c.Get("isOwner")
	userRole, _ := c.Get("role")
	UserID, _ := c.Get("user_id")
	projects, _ := c.Get("projects")
	userName, _ := c.Get("user_name")
	BaseGroup, _ := c.Get("base_group")

	userGroup, ok := groups.([]string)
	if !ok {
		userGroup = []string{}
	}

	userProjects, ok := projects.([]string)
	if !ok {
		userProjects = []string{}
	}

	userIsOwner, ok := isOwner.(bool)
	if !ok {
		userIsOwner = false
	}

	userRolePtr, ok := userRole.(*models.Role)
	var userRoleIs models.Role
	if ok && userRolePtr != nil {
		userRoleIs = *userRolePtr
	} else {
		userRoleIs = models.RoleViewer
	}

	userId, ok := UserID.(string)
	if !ok {
		userId = ""
	}

	user, ok := userName.(string)
	if !ok {
		user = ""
	}

	userBaseGroup, ok := BaseGroup.(string)
	if !ok {
		userBaseGroup = ""
	}

	userDetails := models.UserDetails{
		Groups:    userGroup,
		Role:      userRoleIs,
		IsOwner:   userIsOwner,
		UserID:    userId,
		Projects:  userProjects,
		UserName:  user,
		BaseGroup: userBaseGroup,
	}

	return userDetails, nil
}

func checkRole(c *gin.Context, userDetail models.UserDetails) (err error) {
	method := c.Request.Method
	switch userDetail.Role {
	case models.RoleAdmin, models.RoleOwner:
		return nil
	case models.RoleEditor:
		if method == "GET" || method == "POST" || method == "PUT" || method == "DELETE" {
			return nil
		}
		return errors.New("you are not authorized to perform this action")
	case models.RoleViewer:
		if method == "GET" {
			return nil
		}
		return errors.New("you are not authorized to perform this action")
	default:
		return errors.New("you are not authorized to perform this action")
	}
}

func (h *Handler) handleDepRequest(c *gin.Context, depFunc DepFunc) {
	userDetails, _ := GetUserDetails(c)
	requestDetails := models.RequestDetails{
		GType:      models.GTypes(c.Query("gtype")),
		Name:       c.Param("name"),
		Collection: c.Query("collection"),
		Project:    c.Query("project"),
		User:       userDetails,
	}

	err := checkRole(c, userDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, err := depFunc(requestDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func extractMetadata(c *gin.Context) map[string]string {
	metadata := make(map[string]string)

	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 && len(key) >= 9 && key[:9] == "metadata_" {
			metadata[key[9:]] = values[0]
		}
	}

	return metadata
}

func getParamOrQuery(c *gin.Context, key string) string {
	if value := c.Param(key); value != "" {
		return value
	}
	return c.Query(key)
}

func getOptionalParam(c *gin.Context, key string) string {
	if value := c.Param(key); value != "" {
		return value
	}
	return c.Query(key)
}
