package infra

import (
	"fmt"
	"os"
)

type Config struct {
	// Server
	ServerPort   string `envconfig:"SERVER_PORT" default:"8080"`
	Environment  string `envconfig:"ENVIRONMENT" default:"development"`

	// Database
	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     string `envconfig:"DB_PORT" default:"5432"`
	DBUser     string `envconfig:"DB_USER" default:"postgres"`
	DBPassword string `envconfig:"DB_PASSWORD" default:"postgres"`
	DBName     string `envconfig:"DB_NAME" default:"voicechat"`

	// Redis
	RedisAddr string `envconfig:"REDIS_ADDR" default:"localhost:6379"`

	// JWT
	JWTSecret string `envconfig:"JWT_SECRET" default:"dev-secret-key-change-in-prod"`
}

func LoadConfig() *Config {
	// Load .env file if it exists
	_ = os.Getenv("ENV_FILE") // env vars can override .env

	cfg := &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		DBName:       getEnv("DB_NAME", "voicechat"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:    getEnv("JWT_SECRET", "dev-secret-key-change-in-prod"),
	}

	return cfg
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
