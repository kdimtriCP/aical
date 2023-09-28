package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/kdimtricp/aical/api/auth/v1"
	"github.com/kdimtricp/aical/internal/biz"
	"time"
)

type authRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (ar *authRepo) SetState(_ context.Context, ad *biz.AuthData, duration time.Duration) (*biz.AuthData, error) {
	ar.log.Debug("State data")
	if err := ar.data.cache.Set(ad.State, ad.UserId, duration).Err(); err != nil {
		return nil, err
	}
	if err := ar.data.cache.Set(ad.UserId, ad.State, duration).Err(); err != nil {
		return nil, err
	}
	return ad, nil
}

func (ar *authRepo) CheckState(_ context.Context, ad *biz.AuthData) (*biz.AuthData, error) {
	ar.log.Debug("Callback data")
	tx := ar.data.cache.Get(ad.State)
	if tx == nil {
		ar.log.Error("CallbackStateCheck data: chat id not found")
		return nil, pb.ErrorStateNotFound("chat id not found: %s", ad.State)
	}
	if tx.Err() != nil {
		return nil, tx.Err()
	}
	chatID := tx.Val()

	tx = ar.data.cache.Get(chatID)
	if tx == nil {
		ar.log.Error("CallbackStateCheck data: state not found")
		return nil, pb.ErrorStateNotFound("state not found: %s", ad.State)
	}
	if tx.Err() != nil {
		return nil, tx.Err()
	}
	state := tx.Val()

	if ad.State != state {
		ar.log.Error("CallbackStateCheck data: state not match")
		return nil, pb.ErrorStateNotMatch("state not match: req[%s] check [%s]", ad.State, state)
	}
	return &biz.AuthData{
		State:  state,
		UserId: chatID,
	}, nil
}
