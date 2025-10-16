package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jzhang405/SmartChrome/backend/internal/models"
)

type SessionCache struct {
	client *RedisClient
}

func NewSessionCache(client *RedisClient) *SessionCache {
	return &SessionCache{client: client}
}

func (s *SessionCache) StoreSession(ctx context.Context, session *models.UserSession) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := fmt.Sprintf("session:%s", session.ID)
	return s.client.Set(ctx, key, sessionJSON, 24*time.Hour)
}

func (s *SessionCache) GetSession(ctx context.Context, sessionID string) (*models.UserSession, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	sessionJSON, err := s.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session models.UserSession
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func (s *SessionCache) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.client.Delete(ctx, key)
}

func (s *SessionCache) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.client.Exists(ctx, key)
}

func (s *SessionCache) UpdateSessionActivity(ctx context.Context, sessionID string) error {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.UpdateActivity()
	return s.StoreSession(ctx, session)
}

func (s *SessionCache) StoreConversation(ctx context.Context, conversation *models.Conversation) error {
	conversationJSON, err := json.Marshal(conversation)
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	key := fmt.Sprintf("conversation:%s", conversation.ID)
	return s.client.Set(ctx, key, conversationJSON, 30*24*time.Hour) // 30 days
}

func (s *SessionCache) GetConversation(ctx context.Context, conversationID string) (*models.Conversation, error) {
	key := fmt.Sprintf("conversation:%s", conversationID)
	conversationJSON, err := s.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	var conversation models.Conversation
	if err := json.Unmarshal([]byte(conversationJSON), &conversation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversation: %w", err)
	}

	return &conversation, nil
}

func (s *SessionCache) StoreMessage(ctx context.Context, message *models.Message) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	key := fmt.Sprintf("message:%s", message.ID)
	return s.client.Set(ctx, key, messageJSON, 30*24*time.Hour) // 30 days
}

func (s *SessionCache) GetConversationMessages(ctx context.Context, conversationID string, limit, offset int) ([]*models.Message, error) {
	// This is a simplified implementation - in production, you might want to use
	// Redis sorted sets for better pagination support
	pattern := fmt.Sprintf("message:*")
	
	// For now, return empty list - this would need to be implemented based on your specific Redis setup
	// and how you want to store message lists
	return []*models.Message{}, nil
}