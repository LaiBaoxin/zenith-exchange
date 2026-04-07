package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wwater/zenith-exchange/backend/pkg/response"
)

// AuthMiddleware 身份验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "请求未携带授权令牌")
			c.Abort()
			return
		}

		// 检查 Bearer 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, "授权格式错误，需使用 Bearer Token")
			c.Abort()
			return
		}

		// 解析 Token
		tokenString := parts[1]
		userAddr, err := parseToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "无效的登录凭证或登录已过期")
			c.Abort()
			return
		}
		c.Set("user_address", userAddr)

		c.Next()
	}
}

// parseToken 这是一个占位函数，后续你应该在这里调用你的 JWT 验证方法
func parseToken(token string) (string, error) {
	// TODO: 使用 jwt-go 库验证 token 并返回其中的 address
	// 目前为了让你能跑通测试，我们先假设它返回成功
	// 在实际业务中，你会在这里 return "", errors.New("invalid token")
	return "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", nil
}
