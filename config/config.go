package config

import (
	"fmt"
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

// BaseURL ..
func BaseURL() string {
	if val, ok := os.LookupEnv("BASE_URL"); ok {
		return val
	}

	return "localhost:" + os.Getenv("PORT")
}

// PostgresDSN :nodoc:
func PostgresDSN() string {
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode)
}

// WorkerNamespace ..
func WorkerNamespace() string {
	return "autograd_worker"
}

// WorkerConcurrency :nodoc:
func WorkerConcurrency() uint {
	return 5
}

// RedisWorkerHost :nodoc:
func RedisWorkerHost() string {
	return os.Getenv("REDIS_WORKER_HOST")
}
