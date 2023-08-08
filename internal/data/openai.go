package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/kdimtricp/aical/pkg/openai"
)

type OpenAIRepo struct {
	log    *log.Helper
	client *openai.Client
}

func NewOpenAIRepo(c *conf.OpenAI, logger log.Logger) biz.OpenAIRepo {
	return &OpenAIRepo{
		log:    log.NewHelper(logger),
		client: openai.NewClient(c.Api.Key, c.Api.Model),
	}
}
