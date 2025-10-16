package models

import (
	"time"
)

type Conversation struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id"`
	URL           string    `json:"url"`
	Title         string    `json:"title"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	MessageCount  int       `json:"message_count"`
	IsActive      bool      `json:"is_active"`
}

func NewConversation(sessionID, url, title string) *Conversation {
	now := time.Now()
	return &Conversation{
		ID:           generateUUID(),
		SessionID:    sessionID,
		URL:          url,
		Title:        title,
		CreatedAt:    now,
		UpdatedAt:    now,
		MessageCount: 0,
		IsActive:     true,
	}
}

func (c *Conversation) AddMessage() {
	c.MessageCount++
	c.UpdatedAt = time.Now()
}

func (c *Conversation) SetActive(active bool) {
	c.IsActive = active
	c.UpdatedAt = time.Now()
}

func (c *Conversation) UpdateActivity() {
	c.UpdatedAt = time.Now()
}