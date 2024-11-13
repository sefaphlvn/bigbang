package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCustomResourceList(c *gin.Context) {
	h.handleRequest(c, h.Custom.GetCustomResourceList)
}

func (h *Handler) GetCustomHTTPFilterList(c *gin.Context) {
	h.handleRequest(c, h.Custom.GetCustomHTTPFilterList)
}
