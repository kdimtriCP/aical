package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type Auth struct {
	Id string
}

type AuthRepo interface {
	Login(context.Context) error
	Callback(ctx context.Context) error
	GetURL(ctx context.Context) (string, error)
}

type AuthUsecase struct {
	repo AuthRepo
	log  *log.Helper
}

func NewAuthUsecase(repo AuthRepo, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *AuthUsecase) Login(ctx context.Context) error {
	uc.log.Info("Login biz")
	return uc.repo.Login(ctx)
}

func (uc *AuthUsecase) Callback(ctx context.Context) error {
	uc.log.Info("Callback biz")
	return uc.repo.Callback(ctx)
}

func (uc *AuthUsecase) GetURL(ctx context.Context) (string, error) {
	uc.log.Info("GetURL biz")
	return uc.repo.GetURL(ctx)
}
