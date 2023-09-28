package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	http.Client
	token string
	model string
}

type chatCompletionFunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type ChatCompletionMessage struct {
	Role         string                      `json:"role"`
	Content      string                      `json:"content"`
	Name         string                      `json:"name,omitempty"`
	FunctionCall *chatCompletionFunctionCall `json:"function_call,omitempty"`
}

type chatCompletionErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
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

func (c *Client) DoRequest(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	request.Model = c.model
	req, err := request.httpRequest(c.token)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		var completionResponse ChatCompletionResponse
		if err := json.NewDecoder(resp.Body).Decode(&completionResponse); err != nil {
			return nil, err
		}
		return &completionResponse, nil
	}
	if resp.StatusCode == http.StatusBadRequest {
		var errorResponse chatCompletionErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("bad request: %s", errorResponse.Error.Message)
	}
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
