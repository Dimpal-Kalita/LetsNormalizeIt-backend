package handlers

import (
	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *auth.AuthHandler
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := h.AuthService.Register(request.Email, request.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	err := h.AuthService.VerifyEmail(token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error})
	}
	c.JSON(200, gin.H{"message": "Email Verified Successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := h.AuthService.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Logged in Successfully", "token": token})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	err := h.AuthService.RequestResetPassword(request.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	c.JSON(200, gin.H{"message": "Password Reset Link Sent Successfully"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var request struct {
		NewPassword string `json:"password" binding:"required"`
	}
	token := c.Param("token")
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := h.AuthService.ResetPassword(token, request.NewPassword)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Password Reset Successfully"})
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	valid, err := h.AuthService.ValidateToken(token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(200, gin.H{"message": "Token is valid"})
}
