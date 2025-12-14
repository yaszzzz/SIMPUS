package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type ServerConfig struct {
	Host string
	Port string
}

type AppConfig struct {
	Name string
	Env  string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// .env file is optional in production
	}

	expiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "simpus"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "simpus-secret-key"),
			Expiry: expiry,
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		App: AppConfig{
			Name: getEnv("APP_NAME", "SIMPUS"),
			Env:  getEnv("APP_ENV", "development"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
