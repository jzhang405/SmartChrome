package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// DoubanProvider implements the LLMProvider interface for Douban
type DoubanProvider struct {
	client  *openai.Client
	model   string
	provider string
	baseURL string
	apiKey  string
}

// NewDoubanProvider creates a new Douban provider
func NewDoubanProvider(apiKey, baseURL, model string) (*DoubanProvider, error) {
	if apiKey == "" {
		return nil, errors.New("API key is required")
	}

	if model == "" {
		model = "douban-chat"
	}

	// Douban is compatible with OpenAI API
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	} else {
		config.BaseURL = "https://api.douban.com/v1"
	}

	client := openai.NewClientWithConfig(config)

	return &DoubanProvider{
		client:   client,
		model:    model,
		provider: "douban",
		baseURL:  config.BaseURL,
		apiKey:   apiKey,
	}, nil
}

func (p *DoubanProvider) GetModel() string {
	return p.model
}

func (p *DoubanProvider) GetProvider() string {
	return p.provider
}

func (p *DoubanProvider) Validate() error {
	if p.apiKey == "" {
		return errors.New("API key is required")
	}
	
	if p.client == nil {
		return errors.New("Douban client is not initialized")
	}
	
	return nil
}

func (p *DoubanProvider) Generate(ctx context.Context, prompt string, options ...GenerateOption) (<-chan StreamResponse, error) {
	return p.GenerateStream(ctx, prompt, options...)
}

func (p *DoubanProvider) GenerateStream(ctx context.Context, prompt string, options ...GenerateOption) (<-chan StreamResponse, error) {
	// Apply options
	req := openai.ChatCompletionRequest{
		Model: p.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}

	// Apply options
	for _, option := range options {
		if option.MaxTokens != nil {
			req.MaxTokens = *option.MaxTokens
		}
		if option.Temperature != nil {
			req.Temperature = float32(*option.Temperature)
		}
		if option.TopP != nil {
			req.TopP = float32(*option.TopP)
		}
		if option.Stop != nil {
			req.Stop = option.Stop
		}
	}

	stream, err := p.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion stream: %w", err)
	}

	responseChan := make(chan StreamResponse)

	go func() {
		defer close(responseChan)
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, context.Canceled) {
				responseChan <- StreamResponse{
					Error: err,
				}
				return
			}
			
			if errors.Is(err, context.DeadlineExceeded) {
				responseChan <- StreamResponse{
					Error: err,
				}
				return
			}

			if err != nil {
				// Check if it's EOF (stream finished)
				if strings.Contains(err.Error(), "EOF") {
					responseChan <- StreamResponse{
						Done: true,
					}
					return
				}
				
				responseChan <- StreamResponse{
					Error: err,
				}
				return
			}

			if len(response.Choices) > 0 {
				content := response.Choices[0].Delta.Content
				finishReason := ""
				
				if response.Choices[0].FinishReason != "" {
					finishReason = string(response.Choices[0].FinishReason)
				}

				responseChan <- StreamResponse{
					Content:      content,
					Done:         finishReason != "",
					FinishReason: finishReason,
					// Note: Stream responses don't have usage data
					Usage: Usage{
						PromptTokens:     0,
						CompletionTokens: 0,
						TotalTokens:      0,
					},
				}
			}
		}
	}()

	return responseChan, nil
}