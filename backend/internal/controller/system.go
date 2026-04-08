package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
)

type SystemHandler struct {
	systemService *service.SystemService
}

func NewSystemHandler(svc *service.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: svc,
	}
}

func (h *SystemHandler) GetConfig(c *gin.Context) {
	data := h.systemService.GetGlobalConfig()
	response.Success(c, gin.H{
		"code": 200,
		"data": data,
	})
}
