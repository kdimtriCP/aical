package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Event struct {
	ID         string    `json:"id,omitempty"`
	CalendarID string    `json:"calendar_id,omitempty"`
	Title      string    `json:"title,omitempty"`
	Location   string    `json:"location,omitempty"`
	StartTime  time.Time `json:"start_time,omitempty"`
	EndTime    time.Time `json:"end_time,omitempty"`
}

type Events []*Event

// Diff compares events with input events and return the difference
func (es Events) Diff(events Events) Events {
	var diff Events
	for _, event := range events {
		if !es.Contains(event) {
			diff = append(diff, event)
		}
	}
	return diff
}

// Same compares events with input events and return the same
func (es Events) Same(events Events) Events {
	var same Events
	for _, event := range events {
		if es.Contains(event) {
			same = append(same, event)
		}
	}
	return same
}

// Contains checks if events contains input event
func (es Events) Contains(event *Event) bool {
	for _, e := range es {
		if e.ID == event.ID {
			return true
		}
	}
	return false
}

type EventRepo interface {
	Create(ctx context.Context, event *Event) error
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, event *Event) error
	Get(ctx context.Context, event *Event) (*Event, error)
	List(ctx context.Context, calendarID string) (Events, error)
}

type EventUseCase struct {
	repo EventRepo
	log  *log.Helper
}

func NewEventUseCase(repo EventRepo, logger log.Logger) *EventUseCase {
	return &EventUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Sync syncs calendar events with input events
func (uc *EventUseCase) Sync(ctx context.Context, calendarID string, events Events) error {
	uc.log.Debugf("sync events for calendar %s", calendarID)
	// get events from repo
	dbEvents, err := uc.repo.List(ctx, calendarID)
	if err != nil {
		return err
	}
	// update same events
	sameEvents := dbEvents.Same(events)
	for _, event := range sameEvents {
		if err := uc.repo.Update(ctx, event); err != nil {
			return err
		}
	}
	// delete deleted events
	oldEvents := events.Diff(dbEvents)
	for _, event := range oldEvents {
		if err := uc.repo.Delete(ctx, event); err != nil {
			return err
		}
	}
	// create new events
	newEvents := dbEvents.Diff(events)
	for _, event := range newEvents {
		if err := uc.repo.Create(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
