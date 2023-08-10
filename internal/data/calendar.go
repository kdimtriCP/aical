package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/internal/biz"
	"gorm.io/gorm"
)

type Calendar struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID
	GoogleID      string
	Summary       string
	Events        []*Event
	EventsHistory []*EventHistory
}

func (c *Calendar) biz() *biz.Calendar {
	return &biz.Calendar{
		ID:       c.ID,
		GoogleID: c.GoogleID,
		Summary:  c.Summary,
		UserID:   c.UserID,
	}
}

// marshalCalendar returns data calendar from biz calendar
func marshalCalendar(bc *biz.Calendar) *Calendar {
	return &Calendar{
		ID:       bc.ID,
		GoogleID: bc.GoogleID,
		Summary:  bc.Summary,
		UserID:   bc.UserID,
	}
}

type Calendars []*Calendar

func (cs Calendars) biz() []*biz.Calendar {
	calendars := make([]*biz.Calendar, len(cs))
	for i, calendar := range cs {
		calendars[i] = calendar.biz()
	}
	return calendars
}

type CalendarRepo struct {
	data *Data
	log  *log.Helper
}

func NewCalendarRepo(data *Data, logger log.Logger) biz.CalendarRepo {
	return &CalendarRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *CalendarRepo) Create(ctx context.Context, calendar *biz.Calendar) error {
	r.log.Debugf("CreateAll calendar: %v", calendar)
	c := marshalCalendar(calendar)
	return r.data.db.Create(&c).Error
}

func (r *CalendarRepo) Get(ctx context.Context, calendar *biz.Calendar) (*biz.Calendar, error) {
	r.log.Debugf("Get calendar: %v", calendar)
	c := marshalCalendar(calendar)
	tx := r.data.db.Where(&c).First(&c)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return c.biz(), nil
}

func (r *CalendarRepo) Update(ctx context.Context, calendar *biz.Calendar) error {
	r.log.Debugf("Update calendar: %v", calendar)
	c := marshalCalendar(calendar)
	return r.data.db.Model(&c).Updates(&c).Error
}

func (r *CalendarRepo) Delete(ctx context.Context, calendar *biz.Calendar) error {
	r.log.Debugf("Delete calendar: %v", calendar)
	c := marshalCalendar(calendar)
	return r.data.db.Where(&c).Delete(&c).Error
}

func (r *CalendarRepo) List(ctx context.Context, userID uuid.UUID) ([]*biz.Calendar, error) {
	r.log.Debugf("List cs for user: %v", userID)
	var cs Calendars
	tx := r.data.db.Where("user_id = ?", userID).Find(&cs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return cs.biz(), nil
}
