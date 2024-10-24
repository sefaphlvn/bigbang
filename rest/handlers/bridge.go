package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSnapshotResources(c *gin.Context) {
	h.handleRequest(c, h.Bridge.GetSnapshotResources)
}

func (h *Handler) GetSnapshotKeys(c *gin.Context) {
	h.handleRequest(c, h.Bridge.GetSnapshotKeys)
}

func (h *Handler) GetErrors(c *gin.Context) {
	h.handleRequest(c, h.Bridge.GetErrors)
}
