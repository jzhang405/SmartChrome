package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jzhang405/SmartChrome/backend/internal/handlers"
	"github.com/jzhang405/SmartChrome/backend/internal/middleware"
	"github.com/jzhang405/SmartChrome/backend/pkg/cache"
)

func main() {
	// Initialize configuration
	config := loadConfig()

	// Initialize Redis client
	redisClient, err := cache.NewRedisClient(config.RedisURL, "", 0)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(config.JWTSecret)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorMiddleware())

	// Initialize handlers
	h := handlers.NewHandlers(redisClient, jwtMiddleware)

	// API routes
	api := router.Group("/v1")
	{
		// Session endpoints
		api.POST("/sessions", h.CreateSession)
		api.GET("/sessions/:sessionId", jwtMiddleware.AuthMiddleware(), h.GetSession)
		api.DELETE("/sessions/:sessionId", jwtMiddleware.AuthMiddleware(), h.DeleteSession)

		// Conversation endpoints
		api.POST("/conversations", jwtMiddleware.AuthMiddleware(), h.CreateConversation)
		api.GET("/conversations/:conversationId", jwtMiddleware.AuthMiddleware(), h.GetConversation)
		api.GET("/conversations/:conversationId/messages", jwtMiddleware.AuthMiddleware(), h.GetConversationMessages)
		api.POST("/conversations/:conversationId/messages", jwtMiddleware.AuthMiddleware(), h.SendMessage)

		// Health check
		api.GET("/health", h.HealthCheck)
	}

	// WebSocket endpoint
	router.GET("/v1/stream", jwtMiddleware.AuthMiddleware(), h.StreamHandler)

	// Start server
	srv := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", config.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

type Config struct {
	Port      string
	RedisURL  string
	JWTSecret string
}

func loadConfig() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		RedisURL:  getEnv("REDIS_URL", "localhost:6379"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}