package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) SetExtension(c *gin.Context) {
	h.handleRequest(c, h.Extension.SetExtension)
}

func (h *Handler) GetExtension(c *gin.Context) {
	h.handleRequest(c, h.Extension.GetExtension)
}

func (h *Handler) GetExtensions(c *gin.Context) {
	h.handleRequest(c, h.Extension.GetExtensions)
}

func (h *Handler) ListExtensions(c *gin.Context) {
	h.handleRequest(c, h.Extension.ListExtensions)
}

func (h *Handler) UpdateExtension(c *gin.Context) {
	h.handleRequest(c, h.Extension.UpdateExtensions)
}
