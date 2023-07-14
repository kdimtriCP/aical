package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/kdimtricp/aical/api/auth/v1"
	"github.com/kdimtricp/aical/internal/biz"
)

type AuthRepo struct {
	data   *Data
	google *Google
	log    *log.Helper
}

func NewAuthRepo(data *Data, ggl *Google, logger log.Logger) biz.AuthRepo {
	return &AuthRepo{
		data:   data,
		google: ggl,
		log:    log.NewHelper(logger),
	}
}

func (ar *AuthRepo) Auth(ctx context.Context, ba *biz.Auth) *biz.Auth {
	ar.log.Debug("Auth data")
	ar.data.cache.Set(ba.State, ba.State, biz.STATE_KEY_DURATION)
	url := ar.google.GetAuthURL(ba.State)
	ar.log.Debug("CallbackURL data: %s", url)
	return &biz.Auth{
		URL: url,
	}
}

func (ar *AuthRepo) Callback(ctx context.Context, ba *biz.Auth) error {
	ar.log.Debug("Callback data")
	cachedState := ar.data.cache.Get(ba.State)
	if cachedState == nil {
		ar.log.Error("CallbackStateCheck data: state not found")
		return pb.ErrorStateNotFound("state not found: %s", ba.State)
	}
	if cachedState.Val() != ba.State {
		ar.log.Error("CallbackStateCheck data: state not match")
		return pb.ErrorStateNotMatch("state not match: req[%s] check [%s]", ba.State, cachedState.Val())
	}
	return nil
}
