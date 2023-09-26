package server

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	authpb "github.com/kdimtricp/aical/api/auth/v1"
	chatpb "github.com/kdimtricp/aical/api/chat/v1"
	userpb "github.com/kdimtricp/aical/api/user/v1"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/kdimtricp/aical/internal/service"
	shttp "net/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, logger log.Logger,
	auth *service.AuthService,
	user *service.UserService,
	chat *service.ChatService,
) *http.Server {
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
	chatpb.RegisterChatHTTPServer(srv, chat)
	authpb.RegisterAuthServiceHTTPServer(srv, auth)
	userpb.RegisterUserServiceHTTPServer(srv, user)
	srv.HandleFunc("/", func(w shttp.ResponseWriter, r *shttp.Request) {
		shttp.Redirect(w, r, "login", shttp.StatusTemporaryRedirect)
	})
	return srv
}

const GG_CALENDAR_URL = "https://calendar.google.com/calendar/u/0/r"
const USER_URL_PATH = "/user"

// ResponseFunc redirects State request to url generated from oauth2config
// and Callback request to root url.
func ResponseFunc(w http.ResponseWriter, r *http.Request, i interface{}) error {
	switch v := i.(type) {
	case *authpb.LoginReply:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(shttp.StatusOK)
		if _, err := w.Write([]byte(v.LoginPage)); err != nil {
			panic(err)
		}
	case *authpb.AuthReply:
		shttp.Redirect(w, r, v.Url, shttp.StatusTemporaryRedirect)
	case *authpb.CallbackReply:
		redirectURL := fmt.Sprintf("%s?code=%s&tgid=%d", USER_URL_PATH, v.Code, v.UserID)
		shttp.Redirect(w, r, redirectURL, shttp.StatusTemporaryRedirect)
	case *userpb.CreateUserReply:
		shttp.Redirect(w, r, GG_CALENDAR_URL, shttp.StatusTemporaryRedirect)
	case *chatpb.UserChatResponse:
		if _, err := w.Write([]byte(v.Answer)); err != nil {
			panic(err)
		}
	}
	return nil
}
