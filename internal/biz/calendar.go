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
type Calendars []*Calendar

// Diff compares calendars with input calendars and return the difference
func (cs Calendars) Diff(calendars Calendars) Calendars {
	var diff Calendars
	for _, calendar := range calendars {
		if !cs.Contains(calendar) {
			diff = append(diff, calendar)
		}
	}
	return diff
}

// Same compares calendars with input calendars and return the same
func (cs Calendars) Same(calendars Calendars) Calendars {
	var same Calendars
	for _, calendar := range calendars {
		if cs.Contains(calendar) {
			same = append(same, calendar)
		}
	}
	return same
}

// Contains checks if calendars contains input calendar
func (cs Calendars) Contains(calendar *Calendar) bool {
	for _, c := range cs {
		if c.ID == calendar.ID {
			return true
		}
	}
	return false
}

type CalendarRepo interface {
	Create(ctx context.Context, calendar *Calendar) (*Calendar, error)
	Get(ctx context.Context, calendar *Calendar) (*Calendar, error)
	Update(ctx context.Context, calendar *Calendar) error
	Delete(ctx context.Context, calendar *Calendar) error
	List(ctx context.Context, userID string) (Calendars, error)
}

type CalendarUseCase struct {
	repo CalendarRepo
	log  *log.Helper
}

func NewCalendarUseCase(repo CalendarRepo, logger log.Logger) *CalendarUseCase {
	return &CalendarUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Sync syncs calendars from Google Calendar API.
func (uc *CalendarUseCase) Sync(ctx context.Context, userID string, cals Calendars) error {
	uc.log.Debugf("calendar use case: sync calendars")
	dbCals, err := uc.repo.List(ctx, userID)
	if err != nil {
		return err
	}
	newness := cals.Diff(dbCals)
	for _, calendar := range newness {
		_, err := uc.repo.Create(ctx, calendar)
		if err != nil {
			return err
		}
	}
	outdated := dbCals.Diff(cals)
	for _, calendar := range outdated {
		err := uc.repo.Delete(ctx, calendar)
		if err != nil {
			return err
		}
	}
	same := cals.Same(dbCals)
	for _, calendar := range same {
		err := uc.repo.Update(ctx, calendar)
		if err != nil {
			return err
		}
	}
	return nil
}
