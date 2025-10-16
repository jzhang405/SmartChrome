package models

import (
	"time"
)

type UserSession struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	CreatedAt    time.Time              `json:"created_at"`
	LastActive   time.Time              `json:"last_active"`
	Preferences  map[string]interface{} `json:"preferences"`
	JWTToken     string                 `json:"jwt_token,omitempty"`
}

func NewUserSession(userID string) *UserSession {
	now := time.Now()
	return &UserSession{
		ID:          generateUUID(),
		UserID:      userID,
		CreatedAt:   now,
		LastActive:  now,
		Preferences: make(map[string]interface{}),
	}
}

func (s *UserSession) UpdateActivity() {
	s.LastActive = time.Now()
}

func (s *UserSession) SetPreference(key string, value interface{}) {
	if s.Preferences == nil {
		s.Preferences = make(map[string]interface{})
	}
	s.Preferences[key] = value
}

func (s *UserSession) GetPreference(key string) (interface{}, bool) {
	if s.Preferences == nil {
		return nil, false
	}
	value, exists := s.Preferences[key]
	return value, exists
}