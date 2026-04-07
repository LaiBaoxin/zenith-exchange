package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: data,
		Msg:  "success",
	})
}

// Error 失败返回
func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}
