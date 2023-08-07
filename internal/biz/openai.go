package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/pkg/openai"
)

// OpenAIService .
type OpenAIUseCase struct {
	log    *log.Helper
	client *openai.Client
}

// NewOpenAIUseCase .
func NewOpenAIUseCase(logger log.Logger, client *openai.Client) *OpenAIUseCase {
	return &OpenAIUseCase{
		log:    log.NewHelper(log.With(logger, "module", "usecase/openai")),
		client: client,
	}
}
