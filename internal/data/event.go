package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"time"
)

type Event struct {
	ID         string `gorm:"type:varchar(255);primaryKey"`
	Title      string `gorm:"type:varchar(255)"`
	Location   string `gorm:"type:varchar(255)"`
	StartTime  string `gorm:"type:varchar(255)"`
	EndTime    string `gorm:"type:varchar(255)"`
	CalendarID string `gorm:"type:varchar(255)"`
}

func (e *Event) biz() *biz.Event {
	startTime, err := time.Parse(time.RFC3339, e.StartTime)
	if err != nil {
		panic(err)
	}
	endTime, err := time.Parse(time.RFC3339, e.EndTime)
	if err != nil {
		panic(err)
	}
	return &biz.Event{
		ID:         e.ID,
		CalendarID: e.CalendarID,
		Title:      e.Title,
		Location:   e.Location,
		StartTime:  startTime,
		EndTime:    endTime,
	}
}

// parseCalendar returns data event from biz event
func parseEvent(event *biz.Event) *Event {
	return &Event{
		ID:         event.ID,
		CalendarID: event.CalendarID,
		Title:      event.Title,
		Location:   event.Location,
		StartTime:  event.StartTime.Format(time.RFC3339),
		EndTime:    event.EndTime.Format(time.RFC3339),
	}
}

type Events []*Event

func (es Events) biz() biz.Events {
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

func (r *EventRepo) Create(ctx context.Context, event *biz.Event) error {
	r.log.Debugf("Create event: %v", event)
	e := parseEvent(event)
	return r.data.db.Create(&e).Error
}

func (r *EventRepo) Get(ctx context.Context, event *biz.Event) (*biz.Event, error) {
	r.log.Debugf("Get event: %v", event)
	e := parseEvent(event)
	err := r.data.db.Where(&e).First(&e).Error
	if err != nil {
		return nil, err
	}
	return e.biz(), nil
}

func (r *EventRepo) Update(ctx context.Context, event *biz.Event) error {
	r.log.Debugf("Update event: %v", event)
	e := parseEvent(event)
	return r.data.db.Model(&e).Updates(&e).Error
}

func (r *EventRepo) Delete(ctx context.Context, event *biz.Event) error {
	r.log.Debugf("Delete event: %v", event)
	e := parseEvent(event)
	return r.data.db.Delete(&e).Error
}

func (r *EventRepo) List(ctx context.Context, calendarID string) (biz.Events, error) {
	r.log.Debugf("List events: %v", calendarID)
	var events Events
	err := r.data.db.Where("calendar_id = ?", calendarID).Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events.biz(), nil
}
