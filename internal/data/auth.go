package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/kdimtricp/aical/api/auth/v1"
	"github.com/kdimtricp/aical/internal/biz"
	"time"
)

type AuthRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &AuthRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ar *AuthRepo) SetState(ctx context.Context, state string, duration time.Duration) (string, error) {
	ar.log.Debug("State data")
	if err := ar.data.cache.Set(state, state, duration).Err(); err != nil {
		return "", err
	}
	return state, nil
}

func (ar *AuthRepo) CheckState(ctx context.Context, state string) error {
	ar.log.Debug("Callback data")
	cachedState := ar.data.cache.Get(state)
	if cachedState.Err() != nil {
		return cachedState.Err()
	}
	if cachedState == nil {
		ar.log.Error("CallbackStateCheck data: state not found")
		return pb.ErrorStateNotFound("state not found: %s", state)
	}
	if cachedState.Val() != state {
		ar.log.Error("CallbackStateCheck data: state not match")
		return pb.ErrorStateNotMatch("state not match: req[%s] check [%s]", state, cachedState.Val())
	}
	return nil
}
