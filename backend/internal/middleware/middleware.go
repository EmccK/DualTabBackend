package middleware

import (
	"net/http"
	"strings"

	"dualtab-backend/pkg/jwt"
	"dualtab-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminAuth 管理员认证中间件
func AdminAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未授权，请先登录")
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Token 格式错误")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1], jwtSecret)
		if err != nil {
			response.Unauthorized(c, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 允许浏览器扩展和本地开发环境
		if isAllowedOrigin(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, secret")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isAllowedOrigin 检查 Origin 是否允许访问
func isAllowedOrigin(origin string) bool {
	// 允许浏览器扩展
	if isExtensionOrigin(origin) {
		return true
	}

	// 允许本地开发环境
	if strings.HasPrefix(origin, "http://localhost") ||
	   strings.HasPrefix(origin, "http://127.0.0.1") {
		return true
	}

	// 允许管理后台域名（如果有）
	allowedDomains := []string{
		"https://dualtab-admin.emcck.com",
		"https://dualtab.emcck.com",
		// 可以在这里添加更多允许的域名
	}

	for _, domain := range allowedDomains {
		if origin == domain {
			return true
		}
	}

	return false
}

// isExtensionOrigin 判断是否为浏览器扩展 Origin
func isExtensionOrigin(origin string) bool {
	return strings.HasPrefix(origin, "chrome-extension://") ||
		strings.HasPrefix(origin, "moz-extension://") ||
		strings.HasPrefix(origin, "extension://") ||
		strings.HasPrefix(origin, "safari-web-extension://")
}
