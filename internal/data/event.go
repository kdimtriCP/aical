package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/internal/biz"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:varchar(255);primaryKey;default:gen_random_uuid()"`
	CalendarID uuid.UUID
	GoogleID   string
	Title      string
	Location   string
	StartTime  time.Time
	EndTime    time.Time
	IsUsed     bool
	IsAllDay   bool
}

func (e *Event) biz() *biz.Event {
	return &biz.Event{
		ID:         e.ID,
		GoogleID:   e.GoogleID,
		CalendarID: e.CalendarID,
		Summary:    e.Title,
		Location:   e.Location,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
		StartTime:  e.StartTime,
		EndTime:    e.EndTime,
		IsAllDay:   e.IsAllDay,
	}
}

func marshalEvent(event *biz.Event) *Event {
	return &Event{
		ID:         event.ID,
		GoogleID:   event.GoogleID,
		CalendarID: event.CalendarID,
		Title:      event.Summary,
		Location:   event.Location,
		StartTime:  event.StartTime,
		EndTime:    event.EndTime,
		IsAllDay:   event.IsAllDay,
	}
}

type Events []*Event

func (es Events) biz() []*biz.Event {
	events := make([]*biz.Event, len(es))
	for i, event := range es {
		events[i] = event.biz()
	}
	return events
}

type EventRepo struct {
	data *Data
	log  *log.Helper
}

func NewEventRepo(data *Data, logger log.Logger) biz.EventRepo {
	return &EventRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *EventRepo) Create(ctx context.Context, event *biz.Event) (*biz.Event, error) {
	r.log.Debugf("CreateAll event: %v", event)
	e := marshalEvent(event)
	if err := r.data.db.Create(&e).Error; err != nil {
		return nil, err
	}
	return e.biz(), nil
}

func (r *EventRepo) Get(ctx context.Context, event *biz.Event) (*biz.Event, error) {
	r.log.Debugf("Get event: %v", event)
	e := marshalEvent(event)
	if err := r.data.db.Where(&e).First(&e).Error; err != nil {
		return nil, err
	}
	return e.biz(), nil
}

func (r *EventRepo) Update(ctx context.Context, event *biz.Event) (*biz.Event, error) {
	r.log.Debugf("Update event: %v", event)
	e := marshalEvent(event)
	if err := r.data.db.Model(&e).Updates(&e).Error; err != nil {
		return nil, err
	}
	return e.biz(), nil
}

func (r *EventRepo) Delete(ctx context.Context, event *biz.Event) error {
	r.log.Debugf("Delete event: %v", event)
	e := marshalEvent(event)
	return r.data.db.Delete(&e).Error
}

func (r *EventRepo) List(ctx context.Context, calendarID uuid.UUID) ([]*biz.Event, error) {
	r.log.Debugf("List events: %v", calendarID)
	var events Events
	if err := r.data.db.Where("calendar_id = ?", calendarID).Find(&events).Error; err != nil {
		return nil, err
	}
	return events.biz(), nil
}
