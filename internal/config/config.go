package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// Render provides DATABASE_URL
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return &Config{
			Port:      getEnv("PORT", "8080"),
			JWTSecret: getEnv("JWT_SECRET", "change-this-secret"),
			// For pgx, we can use DATABASE_URL directly
			DBHost:     dbURL,
			DBPort:     "",
			DBUser:     "",
			DBPassword: "",
			DBName:     "",
		}
	}

	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "ecommerce"),
		JWTSecret:  getEnv("JWT_SECRET", "change-this-secret-in-production"),
	}
}

func (c *Config) DBConnectionString() string {
	// If DATABASE_URL is set (Render), use it directly
	if c.DBHost != "" && c.DBPort == "" {
		return c.DBHost
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
