package handlers

import (
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	DB        *mongo.Database
	JWTSecret string
}

type User struct {
	ID                      string `json:"id" bson:"_id,omitempty"`
	Email                   string `json:"email" bson:"email"`
	Password                string `json:"password" bson:"password"`
	Verified                bool   `json:"verified" bson:"verified"`
	VerificationToken       string `json:"verification_token" bson:"verification_token,omitempty"`
	ResetPasswordToken      string `json:"reset_password_token" bson:"reset_password_token,omitempty"`
	RestPasswordTokenExpiry int64  `json:"reset_password_token_expiry" bson:"reset_password_token_expiry,omitempty"`
}

func (s *UserHandler) GetEmailFromToken(c *gin.Context) {
	token := c.Param("token")
	email, err := utils.NewJWTServices(s.JWTSecret).GetEmailFromToken(token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, email)
}
