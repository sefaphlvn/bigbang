package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCustomResourceList(c *gin.Context) {
	h.handleRequest(c, h.Custom.GetCustomResourceList)
}

func (h *Handler) GetCustomHttpFilterList(c *gin.Context) {
	h.handleRequest(c, h.Custom.GetCustomHttpFilterList)
}
