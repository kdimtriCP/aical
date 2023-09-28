package biz

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
)

//goland:noinspection GoSnakeCaseUsage,GoUnnecessarilyExportedIdentifiers
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
	SetState(context.Context, *AuthData, time.Duration) (*AuthData, error)
	CheckState(context.Context, *AuthData) (*AuthData, error)
}

type AuthUsecase struct {
	repo AuthRepo
	log  *log.Helper
}

type AuthData struct {
	State  string
	UserId string
}

func NewAuthUsecase(repo AuthRepo, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *AuthUsecase) SetState(ctx context.Context, userID int64) (string, error) {
	uc.log.Debug("State biz")
	state := randomState()
	userIDstr := fmt.Sprintf("%d", userID)
	if userID == 0 {
		userIDstr = state
	}
	ad := &AuthData{
		State:  state,
		UserId: userIDstr,
	}
	ad, err := uc.repo.SetState(ctx, ad, STATE_KEY_DURATION)
	if err != nil {
		return "", err
	}
	return ad.State, nil
}

func (uc *AuthUsecase) CheckState(ctx context.Context, state string) (int64, error) {
	uc.log.Debug("Callback biz")
	ad, err := uc.repo.CheckState(ctx, &AuthData{
		State: state,
	})
	if err != nil {
		return 0, err
	}
	if ad.UserId == ad.State {
		return 0, nil
	}
	userid, err := strconv.ParseInt(ad.UserId, 10, 64)
	if err != nil {
		return 0, err
	}
	return userid, nil
}
