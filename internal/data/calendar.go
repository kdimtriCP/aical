package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"gorm.io/gorm"
)

type Calendar struct {
	gorm.Model
	ID          string `gorm:"type:varchar(255);primary_key" json:"id,omitempty"`
	Name        string `gorm:"type:varchar(255);not null" json:"name,omitempty"`
	Description string `gorm:"type:varchar(255);not null" json:"description,omitempty"`
	Summary     string `gorm:"type:varchar(255);not null" json:"summary,omitempty"`
	UserID      string `gorm:"type:varchar(255);not null" json:"user_id,omitempty"`
}

func (c *Calendar) biz() *biz.Calendar {
	return &biz.Calendar{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Summary:     c.Summary,
		UserID:      c.UserID,
	}
}

// parse returns data calendar from biz calendar
func (c *Calendar) parse(bc *biz.Calendar) {
	c.ID = bc.ID
	c.Name = bc.Name
	c.Description = bc.Description
	c.Summary = bc.Summary
	c.UserID = bc.UserID
}

type Calendars []*Calendar

func (cs Calendars) biz() biz.Calendars {
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

func (r *CalendarRepo) Create(ctx context.Context, calendar *biz.Calendar) (*biz.Calendar, error) {
	r.log.Debugf("Create calendar: %v", calendar)
	var c *Calendar
	// Check if calendar already exists
	tx := r.data.db.Where("id = ?", calendar.ID).First(&c)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	if tx.RowsAffected > 0 && c != nil {
		r.log.Infof("Calendar already exists: %v", calendar)
		return c.biz(), nil
	}
	c.parse(calendar)
	tx = r.data.db.Create(&c)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return c.biz(), nil
}

func (r *CalendarRepo) Get(ctx context.Context, calendar *biz.Calendar) (*biz.Calendar, error) {
	r.log.Debugf("Get calendar: %v", calendar)
	var c *Calendar
	tx := r.data.db.Where("id = ?", calendar.ID).First(&c)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return c.biz(), nil
}

func (r *CalendarRepo) Update(ctx context.Context, calendar *biz.Calendar) error {
	r.log.Debugf("Update calendar: %v", calendar)
	var c *Calendar
	c.parse(calendar)
	return r.data.db.Model(&Calendar{}).Where("id = ?", calendar.ID).Updates(&c).Error
}

func (r *CalendarRepo) Delete(ctx context.Context, calendar *biz.Calendar) error {
	r.log.Debugf("Delete calendar: %v", calendar)
	var c *Calendar
	c.parse(calendar)
	return r.data.db.Where("id = ?", calendar.ID).Delete(&c).Error
}

func (r *CalendarRepo) List(ctx context.Context, userID string) (biz.Calendars, error) {
	r.log.Debugf("List cs for user: %v", userID)
	var cs Calendars
	tx := r.data.db.Where("user_id = ?", userID).Find(&cs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return cs.biz(), nil
}
