package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI       string
	APP_PORT       string
	EMAIL_ID       string
	EMAIL_PASS     string
	JWT_SECRET_KEY string
	BASE_URL       string
	GIN_MODE       string
}

func Loadconfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		MongoURI:       os.Getenv("MONGO_URI"),
		APP_PORT:       os.Getenv("APP_PORT"),
		EMAIL_ID:       os.Getenv("EMAIL_ID"),
		EMAIL_PASS:     os.Getenv("EMAIL_PASS"),
		JWT_SECRET_KEY: os.Getenv("JWT_SECRET_KEY"),
		BASE_URL:       os.Getenv("BASE_URL"),
		GIN_MODE:       os.Getenv("GIN_MODE"),
	}
}
