package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/internal/service"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
)

type SystemHandler struct {
	sysService service.SystemService
}

func (h *SystemHandler) GetConfig(c *gin.Context) {
	data := h.sysService.GetGlobalConfig()
	response.Success(c, gin.H{
		"code": 200,
		"data": data,
	})
}
