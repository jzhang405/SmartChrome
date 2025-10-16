package models

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type WebpageContent struct {
	ID                  string                 `json:"id"`
	URL                 string                 `json:"url"`
	Title               string                 `json:"title"`
	ExtractedText       string                 `json:"extracted_text"`
	ExtractionTimestamp time.Time              `json:"extraction_timestamp"`
	ContentHash         string                 `json:"content_hash"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

func NewWebpageContent(url, title, extractedText string) *WebpageContent {
	contentHash := generateContentHash(extractedText)
	
	return &WebpageContent{
		ID:                  generateUUID(),
		URL:                 url,
		Title:               title,
		ExtractedText:       extractedText,
		ExtractionTimestamp: time.Now(),
		ContentHash:         contentHash,
		Metadata:            make(map[string]interface{}),
	}
}

func (w *WebpageContent) SetMetadata(key string, value interface{}) {
	if w.Metadata == nil {
		w.Metadata = make(map[string]interface{})
	}
	w.Metadata[key] = value
}

func (w *WebpageContent) GetMetadata(key string) (interface{}, bool) {
	if w.Metadata == nil {
		return nil, false
	}
	value, exists := w.Metadata[key]
	return value, exists
}

func (w *WebpageContent) IsContentMatch(content string) bool {
	return w.ContentHash == generateContentHash(content)
}

func generateContentHash(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}