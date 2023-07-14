package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type UserInfo struct {
	Id    string `json:"id" form:"id" query:"id" validate:"required"`
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
	Code  string `json:"code" form:"code" query:"code"`
}

type UserRepo interface {
	CreateUser(ctx context.Context, code string) (string, error)
	GetUserById(ctx context.Context, id string) (*UserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (*UserInfo, error)
	ListUsers(ctx context.Context) ([]*UserInfo, error)
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, code string) (string, error) {
	log.Debug("create user code: %s", code)
	return uc.repo.CreateUser(ctx, code)
}
