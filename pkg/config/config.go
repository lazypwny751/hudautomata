package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Host string
	Port string

	// Database
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPath     string

	// JWT
	JWTSecret     string
	JWTExpiration string

	// CORS
	CORSOrigins string

	// Environment
	Environment string
}

var AppConfig *Config

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Host:          getEnv("HOST", "0.0.0.0"),
		Port:          getEnv("PORT", "8080"),
		DBDriver:      getEnv("DB_DRIVER", "sqlite"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "huduser"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "hudautomata"),
		DBPath:        getEnv("DB_PATH", "./hudautomata.db"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
		CORSOrigins:   getEnv("CORS_ORIGINS", "*"),
		Environment:   getEnv("GIN_MODE", "debug"),
	}

	return AppConfig
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
