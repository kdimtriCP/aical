package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"time"
)

type ChangeTypeEnum string

const (
	CREATED ChangeTypeEnum = "CREATED"
	UPDATED ChangeTypeEnum = "UPDATED"
	DELETED ChangeTypeEnum = "DELETED"
)

type EventHistory struct {
	ID         uuid.UUID      `json:"history_id,omitempty"`
	EventID    uuid.UUID      `json:"event_id,omitempty"`
	CalendarID uuid.UUID      `json:"calendar_id,omitempty"`
	ChangeType ChangeTypeEnum `json:"change_type_enum,omitempty"`
	ChangeTime time.Time      `json:"change_time,omitempty"`
	PrevEvent  Event          `json:"prev_event"`
	NewEvent   Event          `json:"new_event"`
}

// changeDescription returns a string representation of the change
func (e *EventHistory) changeDescription() string {
	switch e.ChangeType {
	case CREATED:
		return fmt.Sprintf("New event with ID %s was created: %s",
			e.EventID, e.NewEvent.String())
	case UPDATED:
		return fmt.Sprintf("Event with ID %s was updated: "+
			"previous event: %s. "+
			"updated event: %s.",
			e.EventID,
			e.PrevEvent.String(), e.NewEvent.String(),
		)
	case DELETED:
		return fmt.Sprintf("Event with ID %s was deleted: %s",
			e.EventID, e.PrevEvent.String())
	default:
		return fmt.Sprintf("Unknown change type: %s", e.ChangeType)
	}
}

type EventHistoryRepo interface {
	ListCalendarEventHistory(ctx context.Context, calendarID uuid.UUID) ([]*EventHistory, error)
	DeleteCalendarEventHistory(ctx context.Context, calendarID uuid.UUID) error
}

type EventHistoryUseCase struct {
	db  EventHistoryRepo
	log *log.Helper
}

func NewEventHistoryUseCase(repo EventHistoryRepo, logger log.Logger) *EventHistoryUseCase {
	return &EventHistoryUseCase{
		db:  repo,
		log: log.NewHelper(logger),
	}
}

func (uc *EventHistoryUseCase) ListCalendarEventHistory(ctx context.Context, calendarID uuid.UUID) ([]*EventHistory, error) {
	uc.log.Debugf("list events for calendar %s", calendarID)
	return uc.db.ListCalendarEventHistory(ctx, calendarID)
}

func (uc *EventHistoryUseCase) DeleteCalendarEventHistory(ctx context.Context, calendarID uuid.UUID) error {
	uc.log.Debugf("delete events for calendar %s", calendarID)
	return uc.db.DeleteCalendarEventHistory(ctx, calendarID)
}
