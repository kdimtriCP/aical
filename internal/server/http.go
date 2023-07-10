package server

import (
	"github.com/kdimtricp/aical/api/auth/v1"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/kdimtricp/aical/internal/service"
	slog "log"
	shttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, auth *service.AuthService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			logging.Server(logger),
			recovery.Recovery(),
		),
		http.ResponseEncoder(ResponseFunc),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterAuthServiceHTTPServer(srv, auth)
	srv.HandlePrefix("/", loginPageHandler())
	return srv
}

// ResponseFunc redirects Login request to url generated from oauth2config
func ResponseFunc(w http.ResponseWriter, r *http.Request, i interface{}) error {
	switch v := i.(type) {
	case *v1.LoginResponse:
		slog.Println("LoginResponse received")
		shttp.Redirect(w, r, v.Url, shttp.StatusTemporaryRedirect)
	}
	return nil
}

const loginPage = `<html><body>
<a href="/login">Log in with Google</a>
</body></html>
`

func loginPageHandler() shttp.HandlerFunc {
	return func(w shttp.ResponseWriter, r *shttp.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(shttp.StatusOK)
		w.Write([]byte(loginPage))
	}
}
