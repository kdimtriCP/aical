package data

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"time"
)

type Auth struct {
	Id string
}

type AuthRepo struct {
	data *Data
	ggl  *Google
	log  *log.Helper
}

func NewAuthRepo(data *Data, ggl *Google, oai *OpenAI, logger log.Logger) biz.AuthRepo {
	return &AuthRepo{
		data: data,
		ggl:  ggl,
		log:  log.NewHelper(logger),
	}
}

func (ar *AuthRepo) Login(ctx context.Context) error {
	ar.log.Info("Login data")
	return nil
}

func (ar *AuthRepo) Callback(ctx context.Context) error {
	ar.log.Info("Callback data")
	return nil
}

func (ar *AuthRepo) GetURL(ctx context.Context) (string, error) {
	ar.log.Info("GetURL data")
	stateString := randomState()
	ar.data.cache.Set(stateString, 0, time.Second*300)
	url := ar.ggl.config.AuthCodeURL(stateString)
	ar.log.Info("GetURL data: %s", url)
	return url, nil
}

func randomState() string {
	state := make([]byte, 16)
	rand.Read(state)
	return base64.URLEncoding.EncodeToString(state)
}
