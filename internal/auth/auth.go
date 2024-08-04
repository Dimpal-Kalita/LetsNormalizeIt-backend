package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Dimpal-Kalita/LetsNormalizeIt-backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB        *mongo.Database
	JWTSecret string
}

type User struct {
	ID                       string `json:"id" bson:"_id,omitempty"`
	Email                    string `json:"email" bson:"email"`
	Password                 string `json:"password" bson:"password"`
	Verified                 bool   `json:"verified" bson:"verified"`
	VerificationToken        string `json:"verification_token" bson:"verification_token,omitempty"`
	ResetPasswordToken       string `json:"reset_password_token" bson:"reset_password_token,omitempty"`
	ResetPasswordTokenExpiry int64  `json:"reset_password_token_expiry" bson:"reset_password_token_expiry,omitempty"`
}

func (s *AuthHandler) Register(Email, password string) (string, error) {
	// check if the user already exists
	var existingUser User
	err := s.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": Email}).Decode(&existingUser)

	if err == nil {
		if existingUser.Verified {
			return "", fmt.Errorf("User already exists")
		} else {
			err = utils.SendVerificationEmail(Email, existingUser.VerificationToken)
			return existingUser.VerificationToken, err
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// Generate a verification token
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		return "", err
	}
	verificationToken := hex.EncodeToString(token)

	user := User{
		Email:             Email,
		Password:          string(hashedPassword),
		Verified:          false,
		VerificationToken: verificationToken,
	}
	_, err = s.DB.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return "", err
	}
	err = utils.SendVerificationEmail(Email, verificationToken)
	return verificationToken, err
}

func (s *AuthHandler) VerifyEmail(token string) error {
	var user User
	err := s.DB.Collection("users").FindOne(context.TODO(), bson.M{"verification_token": token}).Decode(&user)
	if err != nil {
		return fmt.Errorf("Invalid Token")
	}
	_, err = s.DB.Collection("users").UpdateOne(context.TODO(), bson.M{"email": user.Email}, bson.M{"$set": bson.M{"verified": true, "verification_token": ""}})
	return err
}

func (s *AuthHandler) Login(Email, password string) (string, error) {
	var user User
	err := s.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": Email}).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("Invalid Email")
	}
	if !user.Verified {
		return "", fmt.Errorf("Email not verified")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("Invalid Password")
	}
	token, err := utils.NewJWTServices(s.JWTSecret).GenerateToken(user.Email)
	return token, err
}

func (s *AuthHandler) RequestResetPassword(email string) error {
	var user User
	err := s.DB.Collection("users").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return fmt.Errorf("Invalid Email")
	}
	token := make([]byte, 32)
	_, err = rand.Read(token)
	if err != nil {
		return err
	}
	resetPasswordToken := hex.EncodeToString(token)
	expires := time.Now().Add(time.Hour)
	_, err = s.DB.Collection("users").UpdateOne(context.TODO(), bson.M{"email": email}, bson.M{"$set": bson.M{"reset_password_token": resetPasswordToken, "reset_password_token_expiry": expires}})
	if err != nil {
		return err
	}
	err = utils.SendResetPasswordEmail(email, resetPasswordToken)
	return err
}

func (s *AuthHandler) ResetPassword(token, password string) error {
	var user User
	err := s.DB.Collection("users").FindOne(context.TODO(), bson.M{"reset_password_token": token}).Decode(&user)
	if err != nil {
		return fmt.Errorf("Invalid Token")
	}
	if user.ResetPasswordTokenExpiry < time.Now().Unix() {
		return fmt.Errorf("Token expired")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.DB.Collection("users").UpdateOne(context.TODO(), bson.M{"reset_password_token": token}, bson.M{"$set": bson.M{"password": string(hashedPassword), "reset_password_token": "", "reset_password_token_expiry": 0}})
	return err
}

func (s *AuthHandler) ValidateToken(token string) (bool, error) {
	return utils.NewJWTServices(s.JWTSecret).ValidateToken(token)
}
