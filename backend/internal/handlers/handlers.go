package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jzhang405/SmartChrome/backend/internal/middleware"
	"github.com/jzhang405/SmartChrome/backend/internal/models"
	"github.com/jzhang405/SmartChrome/backend/internal/websocket"
	"github.com/jzhang405/SmartChrome/backend/pkg/cache"
	"github.com/jzhang405/SmartChrome/backend/pkg/llm"
)

type Handlers struct {
	sessionCache     *cache.SessionCache
	jwtMiddleware    *middleware.JWTMiddleware
	streamManager    *websocket.StreamManager
	llmClient        *llm.LLMClient
}

func NewHandlers(redisClient *cache.RedisClient, jwtMiddleware *middleware.JWTMiddleware, llmClient *llm.LLMClient) *Handlers {
	sessionCache := cache.NewSessionCache(redisClient)
	streamManager := websocket.NewStreamManager()

	return &Handlers{
		sessionCache:  sessionCache,
		jwtMiddleware: jwtMiddleware,
		streamManager: streamManager,
		llmClient:     llmClient,
	}
}

func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": "2025-10-16T00:00:00Z",
		"version":   "1.0.0",
	})
}

func (h *Handlers) CreateSession(c *gin.Context) {
	// Create a new session
	session := models.NewUserSession("") // Empty user ID for now
	
	// Store session in cache
	ctx := context.Background()
	if err := h.sessionCache.StoreSession(ctx, session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Generate JWT token
	token, err := h.jwtMiddleware.GenerateToken("", session.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session": session,
		"token":   token,
	})
}

func (h *Handlers) GetSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	
	ctx := context.Background()
	session, err := h.sessionCache.GetSession(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *Handlers) DeleteSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	
	ctx := context.Background()
	if err := h.sessionCache.DeleteSession(ctx, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handlers) CreateConversation(c *gin.Context) {
	var req struct {
		URL           string                `json:"url" binding:"required"`
		Title         string                `json:"title" binding:"required"`
		WebpageContent *models.WebpageContent `json:"webpage_content,omitempty"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get session ID from context
	sessionID, exists := c.Get("session_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session ID not found"})
		return
	}

	// Create conversation
	conversation := models.NewConversation(sessionID.(string), req.URL, req.Title)
	
	// Store conversation in cache
	ctx := context.Background()
	if err := h.sessionCache.StoreConversation(ctx, conversation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	c.JSON(http.StatusCreated, conversation)
}

func (h *Handlers) GetConversation(c *gin.Context) {
	conversationID := c.Param("conversationId")
	
	ctx := context.Background()
	conversation, err := h.sessionCache.GetConversation(ctx, conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}

	c.JSON(http.StatusOK, conversation)
}

func (h *Handlers) GetConversationMessages(c *gin.Context) {
	conversationID := c.Param("conversationId")
	// limit := c.DefaultQuery("limit", "50")
	// offset := c.DefaultQuery("offset", "0")
	
	// Parse limit and offset
	// In a real implementation, you would use these values for pagination
	
	ctx := context.Background()
	messages, err := h.sessionCache.GetConversationMessages(ctx, conversationID, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    len(messages),
	})
}

func (h *Handlers) SendMessage(c *gin.Context) {
	conversationID := c.Param("conversationId")
	
	var req struct {
		Content string `json:"content" binding:"required"`
		Type    string `json:"type" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create message
	message := models.NewMessage(conversationID, models.MessageType(req.Type), req.Content, 0) // Sequence number would be determined in real implementation
	
	// Store message in cache
	ctx := context.Background()
	if err := h.sessionCache.StoreMessage(ctx, message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store message"})
		return
	}

	// If this is a user question, generate LLM response
	if req.Type == string(models.UserQuestion) {
		// Get conversation to get context
		// conversation, err := h.sessionCache.GetConversation(ctx, conversationID)
		// if err == nil {
			// In a real implementation, you would:
			// 1. Get previous messages for context
			// 2. Get webpage content if available
			// 3. Generate prompt with context
			// 4. Call LLM to generate response
			// 5. Stream response back to client
			
			// For now, we'll just log that we would generate a response
			log.Printf("Would generate LLM response for conversation %s with content: %s", conversationID, req.Content)
		// }
	}

	c.JSON(http.StatusCreated, message)
}

func (h *Handlers) StreamHandler(c *gin.Context) {
	h.streamManager.HandleWebSocket(c)
}

// GenerateLLMResponse generates a response using the configured LLM provider
func (h *Handlers) GenerateLLMResponse(ctx context.Context, providerName, prompt string) (<-chan llm.StreamResponse, error) {
	return h.llmClient.Generate(ctx, providerName, prompt)
}