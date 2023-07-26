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
	Create(ctx context.Context, calendar *Calendar) error
	Update(ctx context.Context, calendar *Calendar) error
	Delete(ctx context.Context, calendar *Calendar) error
	Get(ctx context.Context, calendar *Calendar) (*Calendar, error)
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
	same := cals.Same(dbCals)
	for _, calendar := range same {
		if err := uc.repo.Update(ctx, &Calendar{
			ID:          calendar.ID,
			Description: calendar.Description,
			Summary:     calendar.Summary,
			UserID:      userID,
		}); err != nil {
			return err
		}
	}
	outdated := cals.Diff(dbCals)
	for _, calendar := range outdated {
		if err := uc.repo.Delete(ctx, &Calendar{
			ID:          calendar.ID,
			Description: calendar.Description,
			Summary:     calendar.Summary,
			UserID:      userID,
		}); err != nil {
			return err
		}
	}
	newness := dbCals.Diff(cals)
	for _, calendar := range newness {
		if err := uc.repo.Create(ctx, &Calendar{
			ID:          calendar.ID,
			Description: calendar.Description,
			Summary:     calendar.Summary,
			UserID:      userID,
		}); err != nil {
			return err
		}
	}
	return nil
}

// List lists calendars from database.
func (uc *CalendarUseCase) List(ctx context.Context, userID string) (Calendars, error) {
	uc.log.Debugf("calendar use case: list calendars")
	return uc.repo.List(ctx, userID)
}
