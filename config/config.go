package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn(".env file not found")
		return
	}

	logrus.Info("load .env file to os env")
}

// Port ..
func Port() string {
	return os.Getenv("PORT")
}

// Env ..
func Env() string {
	if val, ok := os.LookupEnv("ENV"); ok {
		return val
	}

	return "development"
}

// JWTSecret ..
func JWTSecret() string {
	val, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		logrus.Fatal("JWT_SECRET not provided")
	}
	return val
}
