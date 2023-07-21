package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type Calendar struct {
	ID          string
	Name        string
	Description string
	Summary     string
	UserID      string
}

type CalendarRepo interface {
	CreateCalendar(ctx context.Context, userID string) (*Calendar, error)
}

type CalendarUseCase struct {
	repo CalendarRepo
	log  *log.Helper
}

func NewCalendarUseCase(repo CalendarRepo, logger log.Logger) *CalendarUseCase {
	return &CalendarUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *CalendarUseCase) CreateCalendar(ctx context.Context, userID string) (*Calendar, error) {
	uc.log.Debugf("create calendar for userID: %s", userID)
	return uc.repo.CreateCalendar(ctx, userID)
}
