package main

import (
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/kdimtricp/aical/internal/server"
	"os"

	"github.com/kdimtricp/aical/internal/conf"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X 'main.version=0.0.1a' -X 'main.name=aical'"
var (
	// name is the name of the compiled software.
	name string
	// version is the version of the compiled software.
	version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, HTTPServer *http.Server, GRPCServer *grpc.Server, CronServer *server.CronServer, TGServer *server.TGServer) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(name),
		kratos.Version(version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			HTTPServer,
			GRPCServer,
			CronServer,
			TGServer,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(
		log.NewStdLogger(os.Stdout),
		"ts", log.Timestamp("2006.01.02T15:04:05"),
		"caller", log.DefaultCaller,
	)

	c := config.New(
		config.WithSource(
			env.NewSource("AICAL_"),
			file.NewSource(flagconf),
		),
	)
	defer func(c config.Config) {
		err := c.Close()
		if err != nil {
		}
	}(c)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Google, bc.Openai, bc.Cron, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
