package config

import (
	"fmt"
	"log"
	"os"

	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logs.Warn(".env file not found")
		return
	}

	logs.Info("load .env file to os env")
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

// JWTKey ..
func JWTKey() string {
	val, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("JWT_SECRET not provided")
	}
	return val
}

// BaseURL ..
func BaseURL() string {
	if val, ok := os.LookupEnv("BASE_URL"); ok {
		return val
	}

	return fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))
}

// PostgresDSN ..
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

// WorkerConcurrency ..
func WorkerConcurrency() uint {
	return 5
}

// RedisWorkerHost ..
func RedisWorkerHost() string {
	return os.Getenv("REDIS_WORKER_HOST")
}

// NewRedisPool ..
func NewRedisPool(host string) *redis.Pool {
	return &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(host)
		},
	}
}

// FileUploadPath ..
func FileUploadPath() string {
	val, ok := os.LookupEnv("FILE_UPLOAD_PATH")
	if ok {
		return val
	}

	return "file_upload_path"
}

func AutogradAuthToken() string {
	val, _ := os.LookupEnv("AUTOGRAD_AUTH_TOKEN")
	return val
}

func AutogradServerURL() string {
	val, _ := os.LookupEnv("AUTOGRAD_SERVER_URL")
	return val
}
