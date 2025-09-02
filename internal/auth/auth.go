package auth

import (
	"crypto/subtle"
	"log"
	"net/http"
	"strings"
	"wq_submitter/configs"

	"github.com/gin-gonic/gin"
)

var secretToken string

// init 函数在 main 函数之前自动执行。
// 我们在这里检查并加载必要的配置。
func init() {
	config := configs.GetGlobalConfig()
	// 从环境变量 "WQS_SECRET_TOKEN" 中读取令牌
	secretToken = config.CredentialConfig.Token
	if secretToken == "" {
		// 如果环境变量未设置，打印致命错误并终止程序。
		log.Fatal("FATAL: Environment variable WQS_SECRET_TOKEN is not set.")
	}
}

// APIKeyAuthMiddleware 验证 "Authorization: Bearer <token>" 请求头。
func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.String(http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.String(http.StatusUnauthorized, "Authorization header format must be 'Bearer {token}'")
			c.Abort()
			return
		}

		providedToken := parts[1]

		// 使用 subtle.ConstantTimeCompare 来安全地比较令牌，防止时序攻击。
		if subtle.ConstantTimeCompare([]byte(providedToken), []byte(secretToken)) != 1 {
			c.String(http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		c.Next()
	}
}
