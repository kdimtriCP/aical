package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	ID         uuid.UUID `json:"id,omitempty"`
	CalendarID uuid.UUID `json:"calendar_id,omitempty"`
	GoogleID   string    `json:"google_id,omitempty"`
	Summary    string    `json:"title,omitempty"`
	Location   string    `json:"location,omitempty"`
	StartTime  time.Time `json:"start_time,omitempty"`
	EndTime    time.Time `json:"end_time,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	IsAllDay   bool      `json:"is_all_day,omitempty"`
}

// String .
func (e *Event) String() string {
	return fmt.Sprintf(
		"%s from %s to %s at %s",
		e.Summary,
		e.StartTime.Format(time.RFC3339),
		e.EndTime.Format(time.RFC3339),
		e.Location,
	)
}

// EventRepo .
type EventRepo interface {
	Get(ctx context.Context, event *Event) (*Event, error)
	Create(ctx context.Context, event *Event) (*Event, error)
	Update(ctx context.Context, event *Event) (*Event, error)
	Delete(ctx context.Context, event *Event) error
	List(ctx context.Context, calendarID uuid.UUID) ([]*Event, error)
}

type EventUseCase struct {
	db  EventRepo
	log *log.Helper
}

func NewEventUseCase(repo EventRepo, logger log.Logger) *EventUseCase {
	return &EventUseCase{
		db:  repo,
		log: log.NewHelper(logger),
	}
}

// List lists events from db
func (uc *EventUseCase) List(ctx context.Context, calendarID uuid.UUID) ([]*Event, error) {
	uc.log.Debugf("list events for calendar %s", calendarID)
	return uc.db.List(ctx, calendarID)
}

// Get gets an event from db
func (uc *EventUseCase) Get(ctx context.Context, event *Event) (*Event, error) {
	uc.log.Debugf("get event: %v", event)
	return uc.db.Get(ctx, event)
}

// Create creates a new event
func (uc *EventUseCase) Create(ctx context.Context, event *Event) (*Event, error) {
	uc.log.Debugf("Create event: %v", event)
	return uc.db.Create(ctx, event)
}

// Delete deletes an event
func (uc *EventUseCase) Delete(ctx context.Context, event *Event) error {
	uc.log.Debugf("Delete event: %v", event)
	return uc.db.Delete(ctx, event)
}

// Update updates an event
func (uc *EventUseCase) Update(ctx context.Context, event *Event) (*Event, error) {
	uc.log.Debugf("Update event: %v", event)
	return uc.db.Update(ctx, event)
}

// Sync syncs down database events with incoming events from Google calendar
//   - if event exists in db and not in Google, delete it
//   - if event exists in db and in Google, update it
//   - if event not in db and in Google, create it
func (uc *EventUseCase) Sync(ctx context.Context, calendarID uuid.UUID, events []*Event) error {
	uc.log.Debugf("Sync events for calendar %s", calendarID)
	// List events from db
	dbEvents, err := uc.db.List(ctx, calendarID)
	if err != nil {
		return err
	}
	// Create a map of db events
	dbEventsMap := make(map[string]*Event)
	for _, e := range dbEvents {
		dbEventsMap[e.GoogleID] = e
	}
	// Create a map of incoming events
	incomingEventsMap := make(map[string]*Event)
	for _, e := range events {
		incomingEventsMap[e.GoogleID] = e
	}
	for _, e := range dbEvents {
		// Delete events that are in db but not in Google
		if _, ok := incomingEventsMap[e.GoogleID]; !ok {
			uc.log.Debugf("Delete event %s", e)
			if err := uc.db.Delete(ctx, e); err != nil {
				return err
			}
		} else {
			// Update db events that are present in Google
			uc.log.Debugf("Update event %s", e)
			if _, err := uc.db.Update(ctx, e); err != nil {
				return err
			}
		}
	}
	// Create events that are in Google but not in db
	for _, e := range events {
		if _, ok := dbEventsMap[e.GoogleID]; !ok {
			uc.log.Debugf("Create event %s", e)
			if _, err := uc.db.Create(ctx, &Event{
				CalendarID: calendarID,
				GoogleID:   e.GoogleID,
				Summary:    e.Summary,
				Location:   e.Location,
				StartTime:  e.StartTime,
				EndTime:    e.EndTime,
				IsAllDay:   e.IsAllDay,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
