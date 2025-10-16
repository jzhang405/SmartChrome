package models

import (
	"time"
)

type LLMResponse struct {
	ID           string    `json:"id"`
	MessageID    string    `json:"message_id"`
	StreamID     string    `json:"stream_id"`
	Content      string    `json:"content"`
	IsComplete   bool      `json:"is_complete"`
	TokensUsed   int       `json:"tokens_used"`
	ModelUsed    string    `json:"model_used"`
	CreatedAt    time.Time `json:"created_at"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
}

func NewLLMResponse(messageID, streamID, modelUsed string) *LLMResponse {
	return &LLMResponse{
		ID:         generateUUID(),
		MessageID:  messageID,
		StreamID:   streamID,
		Content:    "",
		IsComplete: false,
		TokensUsed: 0,
		ModelUsed:  modelUsed,
		CreatedAt:  time.Now(),
	}
}

func (r *LLMResponse) AddContent(content string) {
	r.Content += content
}

func (r *LLMResponse) Complete() {
	r.IsComplete = true
	r.CompletedAt = time.Now()
}

func (r *LLMResponse) AddTokens(tokens int) {
	r.TokensUsed += tokens
}