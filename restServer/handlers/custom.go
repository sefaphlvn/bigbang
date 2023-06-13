package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetFilterChainFilters(c *gin.Context) {
	h.handleResource(c, h.Custom.GetFilterChainFilters)
}
