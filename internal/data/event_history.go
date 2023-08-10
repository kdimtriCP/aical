package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/internal/biz"
	"gorm.io/gorm"
	"time"
)

type EventHistory struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID    uuid.UUID
	CalendarID uuid.UUID
	ChangeType biz.ChangeTypeEnum // Тип изменения: CREATED, UPDATED, DELETED
	ChangeTime time.Time          // Время изменения
	PrevEvent  biz.Event          `gorm:"embedded;embeddedPrefix:prev_"`
	NewEvent   biz.Event          `gorm:"embedded;embeddedPrefix:new_"`
}

func (eh *EventHistory) biz() *biz.EventHistory {
	return &biz.EventHistory{
		ID:         eh.ID,
		EventID:    eh.EventID,
		CalendarID: eh.CalendarID,
		ChangeType: eh.ChangeType,
		ChangeTime: eh.ChangeTime,
		PrevEvent:  eh.PrevEvent,
		NewEvent:   eh.NewEvent,
	}
}

func marshalEventHistory(eventHistory *biz.EventHistory) *EventHistory {
	return &EventHistory{
		ID:         eventHistory.ID,
		EventID:    eventHistory.EventID,
		CalendarID: eventHistory.CalendarID,
		ChangeType: eventHistory.ChangeType,
		ChangeTime: eventHistory.ChangeTime,
		PrevEvent:  eventHistory.PrevEvent,
		NewEvent:   eventHistory.NewEvent,
	}
}

type EventHistoryRepo struct {
	data *Data
	log  *log.Helper
}

func NewEventHistoryRepo(data *Data, logger log.Logger) biz.EventHistoryRepo {
	return &EventHistoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *EventHistoryRepo) ListCalendarEventHistory(ctx context.Context, calendarID uuid.UUID) ([]*biz.EventHistory, error) {
	log.Debugf("List event history: %v", calendarID)
	var eventHistories []*EventHistory
	var bizEventHistories []*biz.EventHistory
	if err := r.data.db.Where("calendar_id = ?", calendarID).Find(&eventHistories).Error; err != nil {
		return nil, err
	}
	for _, eventHistory := range eventHistories {
		bizEventHistories = append(bizEventHistories, eventHistory.biz())
	}
	return bizEventHistories, nil
}

func (r *EventHistoryRepo) DeleteCalendarEventHistory(ctx context.Context, calendarID uuid.UUID) error {
	log.Debugf("Delete event history: %v", calendarID)
	return r.data.db.Where("calendar_id = ?", calendarID).Delete(&EventHistory{}).Error
}
