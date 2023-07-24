package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/kdimtricp/aical/api/auth/v1"
	"github.com/kdimtricp/aical/internal/biz"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	log *log.Helper
	uc  *biz.AuthUsecase
	gg  *Google
}

func NewAuthService(
	logger log.Logger,
	uc *biz.AuthUsecase,
	gg *Google,
) *AuthService {
	return &AuthService{
		log: log.NewHelper(logger),
		uc:  uc,
		gg:  gg,
	}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	const loginPage = `<html><body>
<a href="/auth/google/login">Login with Google</a>
</body></html>
`
	s.log.Debug("State request")
	return &pb.LoginReply{
		LoginPage: loginPage,
	}, nil
}

func (s *AuthService) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthReply, error) {
	s.log.Debug("State request")
	state, err := s.uc.SetState(ctx)
	if err != nil {
		return nil, err
	}
	url := s.gg.AuthCodeURL(state)
	s.log.Debug("State url: %s", url)
	return &pb.AuthReply{
		Url: url,
	}, nil
}

func (s *AuthService) Callback(ctx context.Context, req *pb.CallbackRequest) (*pb.CallbackReply, error) {
	s.log.Debug("Callback request")
	if err := s.uc.CheckState(ctx, req.State); err != nil {
		return nil, err
	}
	return &pb.CallbackReply{
		Code: req.Code,
	}, nil
}
