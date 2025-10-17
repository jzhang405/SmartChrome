package models

import (
	"time"
)

type MessageType string

const (
	UserQuestion MessageType = "user_question"
	LLMReply     MessageType = "llm_response"
)

type Message struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	Type           MessageType            `json:"type"`
	Content        string                 `json:"content"`
	Timestamp      time.Time              `json:"timestamp"`
	SequenceNumber int                    `json:"sequence_number"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

func NewMessage(conversationID string, msgType MessageType, content string, sequenceNumber int) *Message {
	return &Message{
		ID:             generateUUID(),
		ConversationID: conversationID,
		Type:           msgType,
		Content:        content,
		Timestamp:      time.Now(),
		SequenceNumber: sequenceNumber,
		Metadata:       make(map[string]interface{}),
	}
}

func (m *Message) SetMetadata(key string, value interface{}) {
	if m.Metadata == nil {
		m.Metadata = make(map[string]interface{})
	}
	m.Metadata[key] = value
}

func (m *Message) GetMetadata(key string) (interface{}, bool) {
	if m.Metadata == nil {
		return nil, false
	}
	value, exists := m.Metadata[key]
	return value, exists
}