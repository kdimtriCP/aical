package biz

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/oauth2"
	"time"
)

const (
	STATE_KEY_DURATION = time.Second * 300
)

type Auth struct {
	State string
	URL   string
	Code  string
	Token *oauth2.Token
}

type AuthRepo interface {
	Auth(context.Context, *Auth) *Auth
	Callback(context.Context, *Auth) error
}

type AuthUsecase struct {
	repo AuthRepo
	log  *log.Helper
}

func NewAuthUsecase(repo AuthRepo, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *AuthUsecase) Auth(ctx context.Context) (*Auth, error) {
	uc.log.Debug("Auth biz")
	state, err := randomState()
	if err != nil {
		return nil, err
	}
	return uc.repo.Auth(ctx, &Auth{
		State: state,
	}), nil
}

func (uc *AuthUsecase) Callback(ctx context.Context, a *Auth) error {
	uc.log.Debug("Callback biz")
	if err := uc.repo.Callback(ctx, a); err != nil {
		return err
	}
	return nil
}

func randomState() (string, error) {
	state := make([]byte, 16)
	_, err := rand.Read(state)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(state), nil
}
