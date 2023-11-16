package middleware

import (
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/rest/api/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("user_id", claims.User_id)
		c.Set("groups", claims.Groups)
		c.Set("isAdmin", helper.Contains(claims.Groups, "admin"))
		c.Next()
	}
}

func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("refresh-token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, "Refresh token required")
			c.Abort()
			return
		}

		claims, err := auth.ValidateRefreshToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("user_id", claims.User_id)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("groups", claims.Groups)

		c.Next()
	}
}

// PathCheck Path Allow
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
