//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/kdimtricp/aical/internal/biz"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/kdimtricp/aical/internal/data"
	"github.com/kdimtricp/aical/internal/server"
	"github.com/kdimtricp/aical/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Google, *conf.OpenAI, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
