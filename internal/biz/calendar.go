package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type Calendar struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	GoogleID string
	Summary  string
}

// String is the string representation of the Calendar struct.
func (c *Calendar) String() string {
	return fmt.Sprintf("GoogleCalendarID=%s, Summary=%s", c.GoogleID, c.Summary)
}

type CalendarRepo interface {
	Create(ctx context.Context, calendar *Calendar) error
	Update(ctx context.Context, calendar *Calendar) error
	Delete(ctx context.Context, calendar *Calendar) error
	Get(ctx context.Context, calendar *Calendar) (*Calendar, error)
	List(ctx context.Context, userID uuid.UUID) ([]*Calendar, error)
}

type CalendarUseCase struct {
	db  CalendarRepo
	log *log.Helper
}

func NewCalendarUseCase(repo CalendarRepo, logger log.Logger) *CalendarUseCase {
	return &CalendarUseCase{
		db:  repo,
		log: log.NewHelper(log.With(logger, "caller", "biz.calendar.usecase")),
	}
}

// findOrCreate finds or creates a calendar in the database.
func (uc *CalendarUseCase) findOrCreate(ctx context.Context, calendar *Calendar) (*Calendar, error) {
	uc.log.Debugf("calendar use case: find or create calendar %s", calendar.ID)
	c, err := uc.db.Get(ctx, calendar)
	if err != nil {
		if err := uc.db.Create(ctx, calendar); err != nil {
			return nil, err
		}
		c, err := uc.db.Get(ctx, calendar)
		return c, err
	}
	return c, nil
}

// ListUserCalendars List lists calendars from database.
func (uc *CalendarUseCase) ListUserCalendars(ctx context.Context, userID uuid.UUID) ([]*Calendar, error) {
	uc.log.Debugf("calendar use case: list calendars")
	return uc.db.List(ctx, userID)
}

// Sync syncs down calendars. It will take incoming calendars and compare them to the ones in the database.
// If the calendar exists in the database, it will update it. If it doesn't exist, it will create it.
// If the calendar exists in the database but not in the incoming calendars, it will delete it.
func (uc *CalendarUseCase) Sync(ctx context.Context, userID uuid.UUID, calendars []*Calendar) error {
	uc.log.Debugf("calendar use case: sync calendars for user %s", userID)
	// Get calendars from database
	dbCalendars, err := uc.db.List(ctx, userID)
	if err != nil {
		return err
	}
	// Create a map of calendars from the database
	dbCalendarsMap := make(map[string]*Calendar)
	for _, c := range dbCalendars {
		dbCalendarsMap[c.GoogleID] = c
	}
	// Create a map of calendars from the incoming calendars
	incomingCalendarsMap := make(map[string]*Calendar)
	for _, c := range calendars {
		incomingCalendarsMap[c.GoogleID] = c
	}
	// Compare the two maps
	for _, c := range dbCalendars {
		// If the calendar exists in the database, update it
		if _, ok := incomingCalendarsMap[c.GoogleID]; ok {
			if err := uc.db.Update(ctx, c); err != nil {
				return err
			}
		} else {
			// If the calendar exists in the database but not in the incoming calendars, delete it
			if err := uc.db.Delete(ctx, c); err != nil {
				return err
			}
		}
	}
	// If the calendar doesn't exist in the database, create it
	for _, c := range calendars {
		if _, ok := dbCalendarsMap[c.GoogleID]; !ok {
			if err := uc.db.Create(ctx, &Calendar{
				UserID:   userID,
				GoogleID: c.GoogleID,
				Summary:  c.Summary,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// get gets a calendar from the database.
func (uc *CalendarUseCase) get(ctx context.Context, calendar *Calendar) (*Calendar, error) {
	uc.log.Debugf("calendar use case: get calendar %s", calendar.ID)
	return uc.db.Get(ctx, calendar)
}
