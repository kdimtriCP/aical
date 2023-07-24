package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/conf"
)

// OpenAI .
type OpenAI struct {
}

// NewOpenAIService .
func NewOpenAIService(c *conf.OpenAI, logger log.Logger) (*OpenAI, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the openai resources")
	}
	return &OpenAI{}, cleanup, nil
}
