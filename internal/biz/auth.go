package biz

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

const (
	STATE_KEY_DURATION = time.Second * 300
)

func randomState() string {
	state := make([]byte, 16)
	_, err := rand.Read(state)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(state)
}

type AuthRepo interface {
	SetState(context.Context, string, time.Duration) (string, error)
	CheckState(context.Context, string) error
}

type AuthUsecase struct {
	repo AuthRepo
	log  *log.Helper
}

func NewAuthUsecase(repo AuthRepo, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *AuthUsecase) SetState(ctx context.Context) (string, error) {
	uc.log.Debug("State biz")
	return uc.repo.SetState(ctx, randomState(), STATE_KEY_DURATION)
}

func (uc *AuthUsecase) CheckState(ctx context.Context, state string) error {
	uc.log.Debug("Callback biz")
	return uc.repo.CheckState(ctx, state)
}
