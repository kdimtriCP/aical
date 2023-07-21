package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID    string `json:"id" form:"id" query:"id" validate:"required"`
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

type UserRepo interface {
	CreateUser(ctx context.Context, code string) (*User, error)
	GetUserById(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, code string) (*User, error) {
	uc.log.Debugf("create user code: %s", code)
	return uc.repo.CreateUser(ctx, code)
}

func (uc *UserUseCase) GetUserById(ctx context.Context, id string) (*User, error) {
	uc.log.Debugf("get user id: %s", id)
	return uc.repo.GetUserById(ctx, id)
}

// ListUsers .
func (uc *UserUseCase) ListUsers(ctx context.Context) ([]*User, error) {
	return uc.repo.ListUsers(ctx)
}
