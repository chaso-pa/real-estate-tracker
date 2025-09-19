package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// OpenAIService handles interactions with OpenAI API
type OpenAIService struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents the request payload for chat completion
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	// For JSON mode
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

// ResponseFormat specifies the format of the response
type ResponseFormat struct {
	Type       string      `json:"type"`                  // "json_object" for JSON mode, "json_schema" for structured outputs
	JSONSchema *JSONSchema `json:"json_schema,omitempty"` // For structured outputs
}

// JSONSchema defines the structure for structured outputs
type JSONSchema struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Schema      map[string]interface{} `json:"schema"`
	Strict      bool                   `json:"strict,omitempty"`
}

// ChatCompletionResponse represents the response from OpenAI
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// NewOpenAIService creates a new OpenAI service instance
func NewOpenAIService() *OpenAIService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is required")
	}

	return &OpenAIService{
		APIKey:  apiKey,
		BaseURL: "https://api.openai.com/v1",
		Client:  &http.Client{},
	}
}

// ChatCompletion performs a regular chat completion request
func (s *OpenAIService) ChatCompletion(messages []Message, model string, options ...func(*ChatCompletionRequest)) (string, error) {
	if model == "" {
		model = "gpt-5-mini"
	}

	req := &ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	// Apply optional parameters
	for _, option := range options {
		option(req)
	}

	return s.sendChatRequest(req)
}

// ChatCompletionWithJSON performs a chat completion request with structured JSON output
func (s *OpenAIService) ChatCompletionWithJSON(messages []Message, model string, options ...func(*ChatCompletionRequest)) (string, error) {
	if model == "" {
		model = "gpt-5-mini" // JSON mode requires newer models
	}

	req := &ChatCompletionRequest{
		Model:    model,
		Messages: messages,
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
		},
	}

	// Apply optional parameters
	for _, option := range options {
		option(req)
	}

	return s.sendChatRequest(req)
}

// ChatCompletionWithStructuredOutput performs a chat completion request with Structured Outputs
func (s *OpenAIService) ChatCompletionWithStructuredOutput(messages []Message, schema *JSONSchema, model string, options ...func(*ChatCompletionRequest)) (string, error) {
	if model == "" {
		model = "gpt-5-mini"
	}

	req := &ChatCompletionRequest{
		Model:    model,
		Messages: messages,
		ResponseFormat: &ResponseFormat{
			Type:       "json_schema",
			JSONSchema: schema,
		},
	}

	// Apply optional parameters
	for _, option := range options {
		option(req)
	}

	return s.sendChatRequest(req)
}

// sendChatRequest sends the actual HTTP request to OpenAI
func (s *OpenAIService) sendChatRequest(req *ChatCompletionRequest) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", s.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.APIKey)

	resp, err := s.Client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// Option functions for customizing requests

// WithTemperature sets the temperature for the request
func WithTemperature(temp float64) func(*ChatCompletionRequest) {
	return func(req *ChatCompletionRequest) {
		req.Temperature = temp
	}
}

// WithMaxTokens sets the max tokens for the request
func WithMaxTokens(tokens int) func(*ChatCompletionRequest) {
	return func(req *ChatCompletionRequest) {
		req.MaxTokens = tokens
	}
}

// WithTopP sets the top_p for the request
func WithTopP(topP float64) func(*ChatCompletionRequest) {
	return func(req *ChatCompletionRequest) {
		req.TopP = topP
	}
}

// NewCustomSchema creates a custom JSON schema
func NewCustomSchema(name, description string, schema map[string]interface{}) *JSONSchema {
	return &JSONSchema{
		Name:        name,
		Description: description,
		Strict:      true,
		Schema:      schema,
	}
}
