package biz

import "github.com/go-kratos/kratos/v2/log"

type TGRepo interface {
}

type TGUseCase struct {
	repo TGRepo
	log  *log.Helper
}

func NewTGUseCase(repo TGRepo, logger log.Logger) *TGUseCase {
	return &TGUseCase{repo: repo, log: log.NewHelper(logger)}
}
