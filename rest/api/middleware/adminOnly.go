package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sefaphlvn/bigbang/rest/handlers"
)

func InitSettingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userDetails, _ := handlers.GetUserDetails(c)
		if !userDetails.IsOwner && userDetails.Role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}
