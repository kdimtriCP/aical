package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/kdimtricp/aical/api/auth/v1"
	"github.com/kdimtricp/aical/internal/biz"
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

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	const loginPage = `<html><body>
<a href="/auth/google/login">Login with Google</a>
</body></html>
`
	s.log.Debug("Auth request")
	return &pb.LoginReply{
		LoginPage: loginPage,
	}, nil
}

func (s *AuthService) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthReply, error) {
	s.log.Debug("Auth request")
	auth, err := s.uc.Auth(ctx)
	if err != nil {
		return nil, err
	}
	s.log.Debug("Auth url: %s", auth.URL)
	return &pb.AuthReply{
		Url: auth.URL,
	}, nil
}
func (s *AuthService) Callback(ctx context.Context, req *pb.CallbackRequest) (*pb.CallbackReply, error) {
	s.log.Debug("Callback request")
	if err := s.uc.Callback(ctx, &biz.Auth{
		State: req.State,
		Code:  req.Code,
	}); err != nil {
		return nil, err
	}
	return &pb.CallbackReply{
		Code: req.Code,
	}, nil
}
