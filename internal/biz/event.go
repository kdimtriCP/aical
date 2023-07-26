package biz

import (
	"context"
	"fmt"
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
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
	IsUsed     bool      `json:"is_used,omitempty"`
	IsAllDay   bool      `json:"is_all_day,omitempty"`
}

// String .
func (e *Event) String() string {
	return fmt.Sprintf(
		"%s from %s to %s at %s",
		e.Title,
		e.StartTime.Format(time.RFC3339),
		e.EndTime.Format(time.RFC3339),
		e.Location,
	)
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

// FilterByStartTime filters events by start time
func (es Events) FilterByStartTime(startMin, startMax time.Time) Events {
	var filtered Events
	for _, event := range es {
		if startMin.Before(event.StartTime) && startMax.After(event.StartTime) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// FilterUsed filters out used events
func (es Events) FilterUsed() Events {
	var filtered Events
	for _, event := range es {
		if !event.IsUsed {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// EventRepo .
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
	// update same events in db only if google event update time is newer
	sameEvents := dbEvents.Same(events)
	for _, event := range sameEvents {
		e, err := uc.repo.Get(ctx, event)
		if err != nil {
			return err
		}
		if event.UpdatedAt.After(e.UpdatedAt) {
			if err := uc.repo.Update(ctx, event); err != nil {
				return err
			}
		}
	}
	// delete deleted events
	oldEvents := events.Diff(dbEvents)
	for _, event := range oldEvents {
		if err := uc.repo.Delete(ctx, event); err != nil {
			return err
		}
	}
	// Create new events
	newEvents := dbEvents.Diff(events)
	for _, event := range newEvents {
		if err := uc.repo.Create(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// ListEventsOptions .
type ListEventsOptions struct {
	StartMin time.Time
	StartMax time.Time
	IsUsed   bool
}

// List lists events from repo
func (uc *EventUseCase) List(ctx context.Context, calendarID string, opts *ListEventsOptions) (Events, error) {
	uc.log.Debugf("list events for calendar %s", calendarID)
	calendarEvents, err := uc.repo.List(ctx, calendarID)
	if err != nil {
		return nil, err
	}
	if opts == nil {
		return calendarEvents, nil
	}
	if !opts.StartMin.IsZero() && !opts.StartMax.IsZero() {
		calendarEvents = calendarEvents.FilterByStartTime(opts.StartMin, opts.StartMax)
	}
	if opts.IsUsed {
		calendarEvents = calendarEvents.FilterUsed()
	}
	return calendarEvents, nil
}

// CreateAll creates new events
func (uc *EventUseCase) CreateAll(ctx context.Context, events Events) error {
	for _, event := range events {
		if err := uc.Create(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Create creates a new event
func (uc *EventUseCase) Create(ctx context.Context, event *Event) error {
	uc.log.Debugf("Create event: %v", event)
	return uc.repo.Create(ctx, event)
}

// MarkAllUsed marks event as used
func (uc *EventUseCase) MarkAllUsed(ctx context.Context, events Events) error {
	uc.log.Debugf("mark %d events as used", len(events))
	for _, event := range events {
		if err := uc.MarkUsed(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// MarkUsed marks event as used
func (uc *EventUseCase) MarkUsed(ctx context.Context, event *Event) error {
	uc.log.Debugf("mark event: %v\nas used", event)
	// Update event as used
	event.IsUsed = true
	return uc.repo.Update(ctx, event)
}
