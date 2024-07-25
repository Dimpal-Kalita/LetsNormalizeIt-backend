package middleware

import (
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/config"
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		// validate token
		valid, err := utils.NewJWTServices(config.Loadconfig().JWT_SECRET_KEY).ValidateToken(token)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if !valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
