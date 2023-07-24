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
type Calendars []*Calendar

type CalendarRepo struct {
	data   *Data
	google *Google
	log    *log.Helper
}

func NewCalendarRepo(data *Data, ggl *Google, logger log.Logger) biz.CalendarRepo {
	return &CalendarRepo{data: data, google: ggl, log: log.NewHelper(logger)}
}

func (calendar *Calendar) Biz() *biz.Calendar {
	return &biz.Calendar{
		ID:          calendar.ID,
		Name:        calendar.Name,
		Description: calendar.Description,
		Summary:     calendar.Summary,
		UserID:      calendar.UserID,
	}
}

func (calendars Calendars) Biz() []*biz.Calendar {
	list := make([]*biz.Calendar, len(calendars))
	for i, calendar := range calendars {
		list[i] = calendar.Biz()
	}
	return list
}

func (c *CalendarRepo) CreateCalendar(ctx context.Context, UserID string) (*biz.Calendar, error) {
	c.log.Debugf("create calendar userID: %s", UserID)
	user := &User{}
	tx := c.data.db.Where("id = ?", UserID).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	c.log.Debugf("create calendar for user: %s", user.Name)
	// Refreshing token and saving new if it changed
	token, err := c.google.RefreshToken(ctx, user.RefreshToken)
	if err != nil {
		return nil, err
	}
	if token.RefreshToken != user.RefreshToken {
		user.RefreshToken = token.RefreshToken
		tx = c.data.db.Save(user)
		if tx.Error != nil {
			return nil, tx.Error
		}
	}
	calendar, err := c.google.CalendarInfo(ctx, token, "primary")
	if err != nil {
		return nil, err
	}
	calendar.UserID = user.ID
	tx = c.data.db.Where("id = ?", calendar.ID).First(calendar)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	if tx.RowsAffected > 0 {
		c.log.Infof("calendar already exists: %s", calendar.Name)
		return calendar.Biz(), nil
	}
	tx = c.data.db.Create(calendar)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return calendar.Biz(), nil
}

// ListCalendars .
func (c *CalendarRepo) ListCalendars(ctx context.Context, UserID string) ([]*biz.Calendar, error) {
	c.log.Debugf("list calendars userID: %s", UserID)
	user := &User{}
	tx := c.data.db.Where("id = ?", UserID).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	c.log.Debugf("list calendars for user: %s", user.Name)
	// Refreshing token and saving new if it changed
	token, err := c.google.RefreshToken(ctx, user.RefreshToken)
	if err != nil {
		return nil, err
	}
	if token.RefreshToken != user.RefreshToken {
		user.RefreshToken = token.RefreshToken
		tx = c.data.db.Save(user)
		if tx.Error != nil {
			return nil, tx.Error
		}
	}
	calendars, err := c.google.ListCalendars(ctx, token)
	if err != nil {
		return nil, err
	}
	for _, calendar := range calendars {
		calendar.UserID = user.ID
		tx = c.data.db.Where("id = ?", calendar.ID).First(calendar)
		if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
			return nil, tx.Error
		}
		if tx.RowsAffected > 0 {
			c.log.Infof("calendar already exists: %s", calendar.Name)
			continue
		}
		tx = c.data.db.Create(calendar)
		if tx.Error != nil {
			return nil, tx.Error
		}
	}
	return calendars.Biz(), nil
}
