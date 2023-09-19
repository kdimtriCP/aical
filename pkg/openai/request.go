package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatCompletionRequest struct {
	Model            string                      `json:"model"`
	Messages         []ChatCompletionMessage     `json:"messages"`
	Functions        []FunctionDescription       `json:"functions,omitempty"`
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

func (r *ChatCompletionRequest) httpRequest(token string) (*http.Request, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// AddFunctionCall adds a function call to the request
func (r *ChatCompletionRequest) AddFunctionCall(name string, arguments string, result string) {
	r.Messages = append(r.Messages, ChatCompletionMessage{
		Role: "assistant",
		FunctionCall: &ChatCompletionFunctionCall{
			Name:      name,
			Arguments: arguments,
		}})
	r.Messages = append(r.Messages, ChatCompletionMessage{
		Role:    "function",
		Content: fmt.Sprintf("result: %s", result),
		Name:    name,
	})
}
