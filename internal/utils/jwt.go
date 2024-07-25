package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTServices struct {
	SecretKey string
}

func NewJWTServices(secretKey string) *JWTServices {
	return &JWTServices{SecretKey: secretKey}
}

func (j *JWTServices) GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTServices) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

func (j *JWTServices) GetEmailFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}
	return claims["email"].(string), nil
}
