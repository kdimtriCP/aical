package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/kdimtricp/aical/api/auth/v1"
	"github.com/kdimtricp/aical/internal/biz"
)

type Auth struct {
	Id string
}

type AuthRepo struct {
	data *Data
	ggl  *Google
	log  *log.Helper
}

func NewAuthRepo(data *Data, ggl *Google, logger log.Logger) biz.AuthRepo {
	return &AuthRepo{
		data: data,
		ggl:  ggl,
		log:  log.NewHelper(logger),
	}
}

func (ar *AuthRepo) Auth(ctx context.Context, ba *biz.Auth) *biz.Auth {
	ar.log.Debug("Auth data")
	ar.data.cache.Set(ba.State, ba.State, biz.STATE_KEY_DURATION)
	url := ar.ggl.config.AuthCodeURL(ba.State)
	ar.log.Debug("CallbackURL data: %s", url)
	return &biz.Auth{
		State: ba.State,
		URL:   url,
	}
}

func (ar *AuthRepo) Callback(ctx context.Context, ba *biz.Auth) (*biz.Auth, error) {
	ar.log.Debug("Callback data")
	cachedState := ar.data.cache.Get(ba.State)
	if cachedState == nil {
		ar.log.Error("CallbackStateCheck data: state not found")
		return nil, pb.ErrorStateNotFound("state not found: %s", ba.State)
	}
	if cachedState.Val() != ba.State {
		ar.log.Error("CallbackStateCheck data: state not match")
		return nil, pb.ErrorStateNotMatch("state not match: req[%s] check [%s]", ba.State, cachedState.Val())
	}
	token, err := ar.ggl.config.Exchange(ctx, ba.Code)
	if err != nil {
		ar.log.Error("CallbackToken data: %s", err)
		return nil, err
	}
	ar.log.Debug("CallbackToken data: %s", token.AccessToken)
	return &biz.Auth{
		State: ba.State,
		Code:  ba.Code,
		URL:   ba.URL,
		Token: token,
	}, nil
}
