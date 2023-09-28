package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/internal/biz"
	"gorm.io/gorm"
	"time"
)

type eventHistory struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID    uuid.UUID
	CalendarID uuid.UUID
	ChangeType biz.ChangeTypeEnum // Тип изменения: CREATED, UPDATED, DELETED
	ChangeTime time.Time          // Время изменения
	PrevEvent  biz.Event          `gorm:"embedded;embeddedPrefix:prev_"`
	NewEvent   biz.Event          `gorm:"embedded;embeddedPrefix:new_"`
}

func (eh *eventHistory) biz() *biz.EventHistory {
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

type eventHistoryRepo struct {
	data *Data
	log  *log.Helper
}

func NewEventHistoryRepo(data *Data, logger log.Logger) biz.EventHistoryRepo {
	return &eventHistoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *eventHistoryRepo) ListCalendarEventHistory(_ context.Context, calendarID uuid.UUID) ([]*biz.EventHistory, error) {
	log.Debugf("List Event history: %v", calendarID)
	var eventHistories []*eventHistory
	var bizEventHistories []*biz.EventHistory
	if err := r.data.db.Where("calendar_id = ?", calendarID).Find(&eventHistories).Error; err != nil {
		return nil, err
	}
	for _, eventHistory := range eventHistories {
		bizEventHistories = append(bizEventHistories, eventHistory.biz())
	}
	return bizEventHistories, nil
}

func (r *eventHistoryRepo) DeleteCalendarEventHistory(_ context.Context, calendarID uuid.UUID) error {
	log.Debugf("Delete Event history: %v", calendarID)
	return r.data.db.Where("calendar_id = ?", calendarID).Delete(&eventHistory{}).Error
}
