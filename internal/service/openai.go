package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/sashabaranov/go-openai"
)

// OpenAI .
type OpenAI struct {
	client *openai.Client
	model  string
}

// NewOpenAIService .
func NewOpenAIService(c *conf.OpenAI, logger log.Logger) (*OpenAI, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the openai resources")
	}
	return &OpenAI{
		client: openai.NewClient(c.Api.Key),
		model:  openai.GPT40613,
	}, cleanup, nil
}

// GenerateNextWeekEvents generates next week events based on the input events
func (s *OpenAI) GenerateNextWeekEvents(ctx context.Context, events []*biz.Event) ([]*biz.Event, error) {
	return nil, nil
}
