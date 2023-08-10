package biz

import (
	"github.com/go-kratos/kratos/v2/log"
)

type OpenAIRepo interface {
}

type OpenAIUseCase struct {
	log  *log.Helper
	repo OpenAIRepo
}

// NewOpenAIUseCase .
func NewOpenAIUseCase(logger log.Logger, repo OpenAIRepo) *OpenAIUseCase {
	return &OpenAIUseCase{
		log:  log.NewHelper(log.With(logger, "module", "usecase/openai")),
		repo: repo,
	}
}
