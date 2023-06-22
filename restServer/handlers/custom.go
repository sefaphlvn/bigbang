package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetFilterChainFilters(c *gin.Context) {
	fmt.Println(c.Get("username"))
	h.handleRequest(c, h.Custom.GetFilterChainFilters)
}
