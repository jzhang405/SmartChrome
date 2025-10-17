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
	"github.com/jzhang405/SmartChrome/backend/config"
	"github.com/jzhang405/SmartChrome/backend/internal/handlers"
	"github.com/jzhang405/SmartChrome/backend/internal/middleware"
	"github.com/jzhang405/SmartChrome/backend/pkg/cache"
	"github.com/jzhang405/SmartChrome/backend/pkg/llm"
)

func main() {
	// Initialize configuration
	config := config.Load()

	// Initialize Redis client
	redisClient, err := cache.NewRedisClient(config.Redis.URL, config.Redis.Password, config.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize LLM client with multiple providers
	llmClient := llm.NewLLMClient()
	
	// Register configured LLM providers
	for _, llmConfig := range config.LLMs {
		var provider llm.LLMProvider
		var err error

		switch llmConfig.Provider {
		case "openai":
			provider, err = llm.NewOpenAIProvider(
				llmConfig.APIKey,
				llmConfig.BaseURL,
				llmConfig.Model,
			)
		case "deepseek":
			provider, err = llm.NewDeepSeekProvider(
				llmConfig.APIKey,
				llmConfig.BaseURL,
				llmConfig.Model,
			)
		case "douban":
			provider, err = llm.NewDoubanProvider(
				llmConfig.APIKey,
				llmConfig.BaseURL,
				llmConfig.Model,
			)
		default:
			log.Printf("Unsupported LLM provider: %s", llmConfig.Provider)
			continue
		}

		if err != nil {
			log.Printf("Failed to initialize %s provider: %v", llmConfig.Provider, err)
			continue
		}

		llmClient.RegisterProvider(llmConfig.Provider, provider)
		
		// Set as default if this is the default provider
		if llmConfig.IsDefault {
			llmClient.SetDefaultProvider(llmConfig.Provider)
		}
	}

	// Initialize JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(config.Auth.JWTSecret)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorMiddleware())

	// Initialize handlers with LLM client
	h := handlers.NewHandlers(redisClient, jwtMiddleware, llmClient)

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
		Addr:    ":" + config.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", config.Server.Port)

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