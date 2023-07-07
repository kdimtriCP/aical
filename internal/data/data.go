package data

import (
	"aical/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGoogle, NewOpenAI)

// Data .
type Data struct {
	// TODO wrapped database client
}

// Google .
type Google struct {
}

// OpenAI .
type OpenAI struct {
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}

// NewGoogle .
func NewGoogle(c *conf.Google, logger log.Logger) (*Google, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the google resources")
	}
	return &Google{}, cleanup, nil
}

// NewOpenAI .
func NewOpenAI(c *conf.OpenAI, logger log.Logger) (*OpenAI, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the openai resources")
	}
	return &OpenAI{}, cleanup, nil
}
