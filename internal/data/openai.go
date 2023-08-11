package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/kdimtricp/aical/pkg/openai"
)

type OpenAIRepo struct {
	log          *log.Helper
	client       *openai.Client
	mctx         *openai.ChatCompletionMessageContext
	functions    []openai.ChatCompletionFunction
	functionsMap map[string]openai.ChatCompletionFunction
}

const INIT_PROMT = ""

func NewOpenAIRepo(c *conf.OpenAI, logger log.Logger) biz.OpenAIRepo {
	return &OpenAIRepo{
		log:          log.NewHelper(logger),
		client:       openai.NewClient(c.Api.Key, c.Api.Model),
		mctx:         openai.NewChatCompletionMessageContext("system", INIT_PROMT),
		functions:    initFunctions(),
		functionsMap: mapFunctions(),
	}
}

func initFunctions() []openai.ChatCompletionFunction {
	return []openai.ChatCompletionFunction{
		biz.CreateEventFunctionDescription(),
		biz.UpdateEventFunctionDescription(),
		biz.DeleteEventFunctionDescription(),
		biz.ListEventsFunctionDescription(),
		biz.ListCalendarsFunctionDescription(),
		biz.CurrentTimeFunctionDescription(),
	}
}

func mapFunctions() map[string]openai.ChatCompletionFunction {
	functions := initFunctions()
	m := make(map[string]openai.ChatCompletionFunction)
	for _, f := range functions {
		m[f.Name] = f
	}
	return m
}

// CreateRequest creates new openai request
func (r *OpenAIRepo) CreateRequest(input string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model:            r.client.Model(),
		Messages:         r.mctx.Add("user", input).Messages(),
		Functions:        r.functions,
		FunctionCall:     nil,
		Temperature:      0,
		TopP:             0,
		N:                0,
		Stream:           false,
		Stop:             nil,
		MaxTokens:        0,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		LogitBias:        nil,
		User:             "",
	}
}

// AddMessage adds message to openai request
