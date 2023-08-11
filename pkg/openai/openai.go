package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	http.Client
	token string
	model string
}

type ChatCompletionFunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type ChatCompletionMessage struct {
	Role         string                      `json:"role"`
	Content      string                      `json:"content"`
	Name         string                      `json:"name,omitempty"`
	FunctionCall *ChatCompletionFunctionCall `json:"function_call,omitempty"`
}

type ChatCompletionMessageContext []ChatCompletionMessage

// NewChatCompletionMessageContext initializes the context with the given prompt
func NewChatCompletionMessageContext(role, prompt string) *ChatCompletionMessageContext {
	c := &ChatCompletionMessageContext{}
	return c.Add(role, prompt)
}

// Add adds a new message to the context
func (c *ChatCompletionMessageContext) Add(role, message string) *ChatCompletionMessageContext {
	*c = append(*c, ChatCompletionMessage{
		Role:    role,
		Content: message,
	})
	return c
}

// Messages returns the messages in the context
func (c *ChatCompletionMessageContext) Messages() []ChatCompletionMessage {
	return *c
}

type ChatCompletionErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int                   `json:"index"`
		Message      ChatCompletionMessage `json:"message"`
		FinishReason string                `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type ChatCompletionFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ChatCompletionRequest struct {
	Model            string                      `json:"model"`
	Messages         []ChatCompletionMessage     `json:"messages"`
	Functions        []ChatCompletionFunction    `json:"functions,omitempty"`
	FunctionCall     *ChatCompletionFunctionCall `json:"function_call,omitempty"`
	Temperature      float64                     `json:"temperature,omitempty"`
	TopP             float64                     `json:"top_p,omitempty"`
	N                int                         `json:"n,omitempty"`
	Stream           bool                        `json:"stream,omitempty"`
	Stop             []string                    `json:"stop,omitempty"`
	MaxTokens        int                         `json:"max_tokens,omitempty"`
	PresencePenalty  float64                     `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64                     `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64          `json:"logit_bias,omitempty"`
	User             string                      `json:"user,omitempty"`
}

// Body returns the request body bytes
func (r *ChatCompletionRequest) Body() ([]byte, error) {
	return json.Marshal(r)
}

func NewClient(apiToken string, model string) *Client {
	return &Client{
		Client: http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		},
		token: apiToken,
		model: model,
	}
}

func (c *Client) NewChatCompletionRequest(body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c *Client) GetCompletion(ctx context.Context, messages []ChatCompletionMessage, functions []ChatCompletionFunction) (*ChatCompletionResponse, error) {
	body, err := json.Marshal(ChatCompletionRequest{
		Model:     c.model,
		Messages:  messages,
		Functions: functions,
	})
	if err != nil {
		return nil, err
	}

	req, err := c.NewChatCompletionRequest(body)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var completionResponse ChatCompletionResponse
		if err := json.NewDecoder(resp.Body).Decode(&completionResponse); err != nil {
			return nil, err
		}
		return &completionResponse, nil
	}
	if resp.StatusCode == http.StatusBadRequest {
		var errorResponse ChatCompletionErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("bad request: %s", errorResponse.Error.Message)
	}
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// Model returns the model name
func (c *Client) Model() string {
	return c.model
}
