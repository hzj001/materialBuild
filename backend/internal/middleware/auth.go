package middleware

import (
	"strings"

	"material-build/pkg/config"
	"material-build/pkg/jwtutil"
	"material-build/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtutil.Parse(tokenStr, cfg.JWT.Secret)
		if err != nil {
			response.Unauthorized(c, "登录已过期")
			c.Abort()
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ContextRole)
		roleStr, _ := role.(string)
		for _, r := range roles {
			if roleStr == r {
				c.Next()
				return
			}
		}
		response.Fail(c, 403, "无权限访问")
		c.Abort()
	}
}

func GetUserID(c *gin.Context) uint64 {
	id, _ := c.Get(ContextUserID)
	v, _ := id.(uint64)
	return v
}
