package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
	"github.com/wwater/zenith-exchange/backend/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "请求未携带授权令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, "授权格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			log.Printf("[AUTH ERROR] Token 解析失败: %v", err)
			response.Error(c, http.StatusUnauthorized, "登录凭证已过期，请重新登录")
			c.Abort()
			return
		}

		c.Set("user_id", int64(claims.UserID))
		c.Set("user_address", claims.Address)

		c.Next()
	}
}
