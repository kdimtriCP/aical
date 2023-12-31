package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/internal/biz"
	"gorm.io/gorm"
	"time"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
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
	History    []*eventHistory
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

type events []*Event

func (es events) biz() []*biz.Event {
	events := make([]*biz.Event, len(es))
	for i, event := range es {
		events[i] = event.biz()
	}
	return events
}

type eventRepo struct {
	data *Data
	log  *log.Helper
}

func NewEventRepo(data *Data, logger log.Logger) biz.EventRepo {
	return &eventRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *eventRepo) Create(_ context.Context, event *biz.Event) (*biz.Event, error) {
	r.log.Debugf("CreateAll Event: %v", event)
	e := marshalEvent(event)
	tx := r.data.db.Begin()
	if err := tx.Create(&e).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Create(&eventHistory{
		EventID:    e.ID,
		CalendarID: e.CalendarID,
		ChangeType: biz.CREATED,
		ChangeTime: time.Now(),
		NewEvent:   *event,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return e.biz(), nil
}

func (r *eventRepo) Get(_ context.Context, event *biz.Event) (*biz.Event, error) {
	r.log.Debugf("Get Event: %v", event)
	e := marshalEvent(event)
	if err := r.data.db.Where(&e).First(&e).Error; err != nil {
		return nil, err
	}
	return e.biz(), nil
}

func (r *eventRepo) Update(_ context.Context, event *biz.Event) (*biz.Event, error) {
	r.log.Debugf("Update Event: %v", event)
	e := marshalEvent(event)
	pe := &Event{}
	tx := r.data.db.Begin()
	if err := tx.Where(&Event{ID: e.ID}).First(&pe).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	bpe := pe.biz()
	if err := tx.Model(&e).Updates(&e).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Create(&eventHistory{
		EventID:    e.ID,
		CalendarID: e.CalendarID,
		ChangeType: biz.UPDATED,
		ChangeTime: time.Now(),
		PrevEvent:  *bpe,
		NewEvent:   *event,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return e.biz(), nil
}

func (r *eventRepo) Delete(_ context.Context, event *biz.Event) error {
	r.log.Debugf("Delete Event: %v", event)
	e := marshalEvent(event)
	tx := r.data.db.Begin()
	if err := tx.Where(&Event{ID: e.ID}).Delete(&Event{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&eventHistory{
		EventID:    e.ID,
		CalendarID: e.CalendarID,
		ChangeType: biz.DELETED,
		ChangeTime: time.Now(),
		PrevEvent:  *event,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *eventRepo) List(_ context.Context, calendarID uuid.UUID) ([]*biz.Event, error) {
	r.log.Debugf("List events: %v", calendarID)
	var events events
	if err := r.data.db.Where("calendar_id = ?", calendarID).Find(&events).Error; err != nil {
		return nil, err
	}
	return events.biz(), nil
}
