package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID           string `json:"id" form:"id" query:"id" validate:"required"`
	Name         string `json:"name" form:"name" query:"name"`
	Email        string `json:"email" form:"email" query:"email"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" query:"refresh_token"`
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, user *User) (*User, error)
	List(ctx context.Context) ([]*User, error)
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(logger)}
}

// Create creates user in database
func (uc *UserUseCase) Create(ctx context.Context, user *User) error {
	uc.log.Debugf("Create user code: %v", user)
	return uc.repo.Create(ctx, user)
}

// Get gets user from database
func (uc *UserUseCase) Get(ctx context.Context, user *User) (*User, error) {
	uc.log.Debugf("get user: %v", user)
	return uc.repo.Get(ctx, user)
}

// ListUsers lists all users from database
func (uc *UserUseCase) List(ctx context.Context) ([]*User, error) {
	uc.log.Debugf("list users")
	return uc.repo.List(ctx)
}
