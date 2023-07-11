package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"net/http"

	pb "github.com/kdimtricp/aical/api/auth/v1"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	uc  *biz.AuthUsecase
	log *log.Helper
}

func NewAuthService(uc *biz.AuthUsecase, logger log.Logger) *AuthService {
	return &AuthService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	const loginPage = `<html><body>
<a href="/auth/google/login">Login with Google</a>
</body></html>
`
	s.log.Debug("Auth request")
	return &pb.LoginResponse{
		LoginPage: loginPage,
	}, nil
}

func (s *AuthService) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.log.Debug("Auth request")
	auth, err := s.uc.Auth(ctx)
	if err != nil {
		return nil, err
	}
	s.log.Debug("Auth url: %s", auth.URL)
	return &pb.AuthResponse{
		Url: auth.URL,
	}, nil
}
func (s *AuthService) Callback(ctx context.Context, req *pb.CallbackRequest) (*pb.CallbackResponse, error) {
	s.log.Debug("Callback request")
	auth, err := s.uc.Callback(ctx, &biz.Auth{
		State: req.State,
		Code:  req.Code,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CallbackResponse{
		Token: auth.Token.AccessToken,
	}, nil
}

// AuthResponseFunc redirects Auth request to url generated from oauth2config
// and Callback request to root url.
func AuthResponseFunc(w http.ResponseWriter, r *http.Request, i interface{}) error {
	switch v := i.(type) {
	case *pb.LoginResponse:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(v.LoginPage)); err != nil {
			panic(err)
		}
	case *pb.AuthResponse:
		http.Redirect(w, r, v.Url, http.StatusTemporaryRedirect)
	case *pb.CallbackResponse:
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(v.Token)); err != nil {
			panic(err)
		}
	}
	return nil
}
