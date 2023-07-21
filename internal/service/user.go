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
	uc  *biz.UserUseCase
	log *log.Helper
}

func NewUserService(uc *biz.UserUseCase, logger log.Logger) *UserService {
	return &UserService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	s.log.Debugf("create user code: %s", req.Code)
	if req.Code == "" {
		return nil, errors.New(http.StatusBadRequest, "code is empty", "code is empty")
	}
	user, err := s.uc.CreateUser(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserReply{
		UserID: user.ID,
	}, nil
}
