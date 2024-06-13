package handlers

import (
	"errors"
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

// This function handles a request in the Handler struct.
// It retrieves the necessary data from the context, including the groups and isAdmin parameters.
// It then sets the resourceDetails struct with the given parameters and decodes the resource.
// It then calls the dbFunc with the resource and resourceDetails, and stores the response in the response variable.
// Finally, it returns the response as a JSON object with the status OK.
func (h *Handler) handleRequest(c *gin.Context, dbFunc DBFunc) {
	userDetails, _ := getUserDetails(c)

	resourceDetails := models.ResourceDetails{
		CanonicalName: c.Param("canonical_name"),
		GType:         models.GTypes(c.Query("gtype")),
		Category:      c.Query("category"),
		Name:          c.Param("name"),
		Collection:    c.Query("collection"),
		SaveOrPublish: c.Query("save_or_publish"),
		User:          userDetails,
	}

	err := checkRole(c, userDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Check for a version parameter in the path and query parameters
	if version := c.Param("version"); version != "" {
		resourceDetails.Version = version
	} else if version := c.Query("version"); version != "" {
		resourceDetails.Version = version
	}

	// Check for a type parameter in the path and query parameters
	if ltype := c.Param("type"); ltype != "" {
		resourceDetails.Type = models.KnownTYPES(ltype)
	} else if ltype := c.Query("type"); ltype != "" {
		resourceDetails.Type = models.KnownTYPES(ltype)
	}

	// Decode the resource from the request
	resource, err := decodeResource(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Call the dbFunc with the resource and resourceDetails
	response, err := dbFunc(resource, resourceDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Return the response as a JSON object with the status OK
	c.JSON(http.StatusOK, response)
}

func getUserDetails(c *gin.Context) (models.UserDetails, error) {
	groups, _ := c.Get("groups")
	isAdmin, _ := c.Get("isAdmin")
	userRole, _ := c.Get("role")
	UserID, _ := c.Get("user_id")

	userGroup, ok := groups.([]string)
	if !ok {
		userGroup = []string{}
	}

	userIsAdmin, ok := isAdmin.(bool)
	if !ok {
		isAdmin = false
	}

	userRoleIs, ok := userRole.(string)
	if !ok {
		userRole = ""
	}

	userId, ok := UserID.(string)
	if !ok {
		userRole = ""
	}

	userDetails := models.UserDetails{
		Groups:  userGroup,
		Role:    userRoleIs,
		IsAdmin: userIsAdmin,
		UserID:  userId,
	}

	return userDetails, nil
}

func checkRole(c *gin.Context, userDetail models.UserDetails) (err error) {
	method := c.Request.Method
	switch userDetail.Role {
	case "admin":
		return nil
	case "editor":
		if method == "GET" || method == "POST" || method == "PUT" || method == "DELETE" {
			return nil
		}
		return errors.New("you are not authorized to perform this action")
	case "readonly":
		if method == "GET" {
			return nil
		}
		return errors.New("you are not authorized to perform this action")
	default:
		return errors.New("you are not authorized to perform this action")
	}
}
