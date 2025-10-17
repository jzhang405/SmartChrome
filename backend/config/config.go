package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	LLMs     []LLMConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	URL                string
	MaxConnections     int
	MaxIdleConnections int
}

type AuthConfig struct {
	JWTSecret     string
	JWTExpiration int
}

type LLMConfig struct {
	Provider    string
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature float64
	IsDefault   bool
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

func Load() *Config {
	// Load default OpenAI config for backward compatibility
	openaiConfig := LLMConfig{
		Provider:    "openai",
		APIKey:      getEnv("OPENAI_API_KEY", ""),
		BaseURL:     getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		Model:       getEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
		MaxTokens:   getEnvAsInt("OPENAI_MAX_TOKENS", 1000),
		Temperature: getEnvAsFloat("OPENAI_TEMPERATURE", 0.7),
		IsDefault:   true,
	}

	// Load additional LLM providers from environment variables
	// This is a simplified approach - in production, you might want to load from a config file
	llmConfigs := []LLMConfig{openaiConfig}

	// Example for DeepSeek
	if deepseekAPIKey := getEnv("DEEPSEEK_API_KEY", ""); deepseekAPIKey != "" {
		deepseekConfig := LLMConfig{
			Provider:    "deepseek",
			APIKey:      deepseekAPIKey,
			BaseURL:     getEnv("DEEPSEEK_BASE_URL", "https://api.deepseek.com/v1"),
			Model:       getEnv("DEEPSEEK_MODEL", "deepseek-chat"),
			MaxTokens:   getEnvAsInt("DEEPSEEK_MAX_TOKENS", 1000),
			Temperature: getEnvAsFloat("DEEPSEEK_TEMPERATURE", 0.7),
			IsDefault:   false,
		}
		llmConfigs = append(llmConfigs, deepseekConfig)
	}

	// Example for Douban
	if doubanAPIKey := getEnv("DOUBAN_API_KEY", ""); doubanAPIKey != "" {
		doubanConfig := LLMConfig{
			Provider:    "douban",
			APIKey:      doubanAPIKey,
			BaseURL:     getEnv("DOUBAN_BASE_URL", "https://api.douban.com/v1"),
			Model:       getEnv("DOUBAN_MODEL", "douban-chat"),
			MaxTokens:   getEnvAsInt("DOUBAN_MAX_TOKENS", 1000),
			Temperature: getEnvAsFloat("DOUBAN_TEMPERATURE", 0.7),
			IsDefault:   false,
		}
		llmConfigs = append(llmConfigs, doubanConfig)
	}

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
		LLMs: llmConfigs,
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