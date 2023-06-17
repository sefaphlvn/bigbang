package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sefaphlvn/bigbang/restServer/auth"
	"github.com/sefaphlvn/bigbang/restServer/helper"
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

		fmt.Println(claims.User_id)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("user_id", claims.User_id)
		c.Set("groups", claims.Groups)
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

// CORS Allow
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token, refresh-token")
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
