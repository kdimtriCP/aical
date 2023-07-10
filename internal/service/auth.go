package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"

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
	s.log.Info("Login request")
	url, err := s.uc.GetURL(ctx)
	if err != nil {
		return nil, err
	}
	s.log.Info("Login url: %s", url)
	return &pb.LoginResponse{
		Url: url,
	}, nil
}
func (s *AuthService) Callback(ctx context.Context, req *pb.CallbackRequest) (*pb.CallbackResponse, error) {
	return &pb.CallbackResponse{}, nil
}
