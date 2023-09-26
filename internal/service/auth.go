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
	gg  *biz.GoogleUseCase
}

func NewAuthService(
	logger log.Logger,
	uc *biz.AuthUsecase,
	gg *biz.GoogleUseCase,
) *AuthService {
	return &AuthService{
		log: log.NewHelper(logger),
		uc:  uc,
		gg:  gg,
	}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	const loginPage = `<html><body>
<a href="/auth/google/login">Login with GoogleRepo</a>
</body></html>
`
	s.log.Debug("Login request: %v", req)
	return &pb.LoginReply{
		LoginPage: loginPage,
	}, nil
}

func (s *AuthService) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthReply, error) {
	s.log.Debug("Auth request: %v", req)
	state, err := s.uc.SetState(ctx, 0)
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
	s.log.Debug("Callback request: %v", req)
	userID, err := s.uc.CheckState(ctx, req.State)
	if err != nil {
		return nil, err
	}
	return &pb.CallbackReply{
		Code:   req.Code,
		UserID: userID,
	}, nil
}

func (s *AuthService) AuthWithID(ctx context.Context, id int64) (string, error) {
	s.log.Debug("Auth with id: \"%d\" request", id)
	state, err := s.uc.SetState(ctx, id)
	if err != nil {
		return "", err
	}
	url := s.gg.AuthCodeURL(state)
	s.log.Debug("State url: %s", url)
	return url, nil
}
