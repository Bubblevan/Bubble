package middleware

import (
	"backend/models"
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			utils.Fail(c, http.StatusUnauthorized, "用户未登录")
			c.Abort()
			return
		}

		// 解析和验证JWT
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "无效的token")
			c.Abort()
			return
		}

		// 在上下文中存储用户信息
		var user models.User
		if err := utils.DB.First(&user, claims.UserID).Error; err != nil {
			utils.Fail(c, http.StatusUnauthorized, "用户不存在")
			c.Abort()
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
