package llm

import (
	"context"
)

// LLMProvider defines the interface for an LLM provider
type LLMProvider interface {
	GetModel() string
	GetProvider() string
	Generate(ctx context.Context, prompt string, options ...GenerateOption) (<-chan StreamResponse, error)
	GenerateStream(ctx context.Context, prompt string, options ...GenerateOption) (<-chan StreamResponse, error)
	Validate() error
}

// StreamResponse represents a single response from the LLM stream
type StreamResponse struct {
	Content     string
	Done        bool
	Error       error
	Usage       Usage
	FinishReason string
}

// Usage represents token usage
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// GenerateOption represents an option for generating text
type GenerateOption struct {
	MaxTokens   *int
	Temperature *float64
	TopP        *float64
	Stop        []string
}

// Option functions
func WithMaxTokens(maxTokens int) GenerateOption {
	return GenerateOption{MaxTokens: &maxTokens}
}

func WithTemperature(temperature float64) GenerateOption {
	return GenerateOption{Temperature: &temperature}
}

func WithTopP(topP float64) GenerateOption {
	return GenerateOption{TopP: &topP}
}

func WithStop(stop []string) GenerateOption {
	return GenerateOption{Stop: stop}
}

// LLMClient manages multiple LLM providers
type LLMClient struct {
	providers map[string]LLMProvider
	defaultProvider string
}

func NewLLMClient() *LLMClient {
	return &LLMClient{
		providers: make(map[string]LLMProvider),
	}
}

func (c *LLMClient) RegisterProvider(name string, provider LLMProvider) {
	c.providers[name] = provider
	if provider != nil {
		// Set the first provider as the default if no default is set
		if c.defaultProvider == "" {
			c.defaultProvider = name
		}
	}
}

func (c *LLMClient) SetDefaultProvider(name string) {
	if _, exists := c.providers[name]; exists {
		c.defaultProvider = name
	}
}

func (c *LLMClient) GetProvider(name string) (LLMProvider, bool) {
	provider, exists := c.providers[name]
	return provider, exists
}

func (c *LLMClient) GetDefaultProvider() (LLMProvider, bool) {
	return c.GetProvider(c.defaultProvider)
}

func (c *LLMClient) Generate(ctx context.Context, providerName, prompt string, options ...GenerateOption) (<-chan StreamResponse, error) {
	provider, exists := c.GetProvider(providerName)
	if !exists {
		// Try to use the default provider
		provider, exists = c.GetDefaultProvider()
		if !exists {
			return nil, NewProviderNotFoundError(providerName)
		}
	}

	return provider.GenerateStream(ctx, prompt, options...)
}

// ProviderNotFoundError indicates that the requested provider was not found
type ProviderNotFoundError struct {
	ProviderName string
}

func (e *ProviderNotFoundError) Error() string {
	return "LLM provider not found: " + e.ProviderName
}

func NewProviderNotFoundError(providerName string) error {
	return &ProviderNotFoundError{ProviderName: providerName}
}