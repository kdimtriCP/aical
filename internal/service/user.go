package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"net/http"

	pb "github.com/kdimtricp/aical/api/user/v1"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	log *log.Helper
	uc  *biz.UserUseCase
	gg  *Google
}

func NewUserService(
	logger log.Logger,
	uc *biz.UserUseCase,
	gg *Google,
) *UserService {
	return &UserService{
		log: log.NewHelper(logger),
		uc:  uc,
		gg:  gg,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	s.log.Debugf("create user code: %s", req.Code)
	if req.Code == "" {
		return nil, errors.New(http.StatusBadRequest, "code is empty", "code is empty")
	}
	u, err := s.gg.UserRegistration(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	user, err := s.uc.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserReply{
		UserID: user.ID,
	}, nil
}
