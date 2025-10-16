package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	LLM      LLMConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	URL             string
	MaxConnections  int
	MaxIdleConnections int
}

type AuthConfig struct {
	JWTSecret     string
	JWTExpiration int
}

type LLMConfig struct {
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "localhost"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
		},
		Database: DatabaseConfig{
			URL:                getEnv("DATABASE_URL", ""),
			MaxConnections:     getEnvAsInt("DB_MAX_CONNECTIONS", 25),
			MaxIdleConnections: getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 24),
		},
		LLM: LLMConfig{
			APIKey:      getEnv("OPENAI_API_KEY", ""),
			Model:       getEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
			MaxTokens:   getEnvAsInt("OPENAI_MAX_TOKENS", 1000),
			Temperature: getEnvAsFloat("OPENAI_TEMPERATURE", 0.7),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}