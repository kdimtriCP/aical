package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/oauth2"
	calendarAPI "google.golang.org/api/calendar/v3"
)

type GoogleRepo interface {
	AuthCodeURL(state string) string
	TokenExchange(ctx context.Context, code string) (*oauth2.Token, error)
	TokenSource(ctx context.Context, refreshToken string) (*oauth2.Token, error)
	UserInfo(ctx context.Context, token *oauth2.Token) (*User, error)
	ListUserCalendars(ctx context.Context, token *oauth2.Token) ([]*Calendar, error)
	CreateNewCalendar(ctx context.Context, token *oauth2.Token, calendarName string) (*Calendar, error)
	CreateCalendarEvent(ctx context.Context, token *oauth2.Token, event *Event, calendarID string) (*Event, error)
	UpdateCalendarEvent(ctx context.Context, token *oauth2.Token, event *Event, calendarID string) (*Event, error)
	GetCalendarEvent(ctx context.Context, token *oauth2.Token, event *Event, calendarID string) (*Event, error)
	DeleteCalendarEvent(ctx context.Context, token *oauth2.Token, event *Event, calendarID string) error
	ListCalendarEvents(ctx context.Context, token *oauth2.Token, calendarID string, opts *GoogleListEventsOption) ([]*Event, error)
}

type GoogleUseCase struct {
	repo GoogleRepo
	log  *log.Helper
}

func NewGoogleUseCase(repo GoogleRepo, logger log.Logger) *GoogleUseCase {
	return &GoogleUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// AuthCodeURL returns the URL to OAuth 2.0 provider's consent page
func (uc *GoogleUseCase) AuthCodeURL(state string) string {
	uc.log.Debugf("AuthCodeURL")
	return uc.repo.AuthCodeURL(state)
}

// TokenExchange exchanges an authorization code for a token
func (uc *GoogleUseCase) TokenExchange(ctx context.Context, code string) (*oauth2.Token, error) {
	uc.log.Debugf("TokenExchange code: %s", code)
	return uc.repo.TokenExchange(ctx, code)
}

// TokenSource returns a token source
func (uc *GoogleUseCase) TokenSource(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	uc.log.Debugf("TokenSource refreshToken: %s", refreshToken)
	return uc.repo.TokenSource(ctx, refreshToken)
}

// UserInfo creates user in database
func (uc *GoogleUseCase) UserInfo(ctx context.Context, token *oauth2.Token) (*User, error) {
	uc.log.Debugf("create user code: %s", token)
	return uc.repo.UserInfo(ctx, token)
}

// ListUserCalendars lists user calendars
func (uc *GoogleUseCase) ListUserCalendars(ctx context.Context, token *oauth2.Token) ([]*Calendar, error) {
	uc.log.Debugf("ListUserCalendars")
	return uc.repo.ListUserCalendars(ctx, token)
}

// GoogleListEventsOption is the option for list events
type GoogleListEventsOption struct {
	TimeMin           string
	TimeMax           string
	UpdatedMin        string
	MaxResults        int64
	OrderByUpdateTime bool
}

// ListEventsCallWithOpts returns a call to list events
func (o *GoogleListEventsOption) ListEventsCallWithOpts(call *calendarAPI.EventsListCall) *calendarAPI.EventsListCall {
	if o.TimeMin != "" {
		call = call.TimeMin(o.TimeMin)
	}
	if o.TimeMax != "" {
		call = call.TimeMax(o.TimeMax)
	}
	if o.UpdatedMin != "" {
		call = call.UpdatedMin(o.UpdatedMin)
	}
	if o.MaxResults > 0 {
		call = call.MaxResults(o.MaxResults)
	}
	if o.OrderByUpdateTime {
		call = call.OrderBy("updated")
	}
	return call
}

// ListEventsInstancesCallWithOpts returns a call to list events instances
func (o *GoogleListEventsOption) ListEventsInstancesCallWithOpts(call *calendarAPI.EventsInstancesCall) *calendarAPI.EventsInstancesCall {
	if o.TimeMin != "" {
		call = call.TimeMin(o.TimeMin)
	}
	if o.TimeMax != "" {
		call = call.TimeMax(o.TimeMax)
	}
	if o.MaxResults > 0 {
		call = call.MaxResults(o.MaxResults)
	}
	return call
}

// ListCalendarEvents lists calendar events with options
func (uc *GoogleUseCase) ListCalendarEvents(ctx context.Context, token *oauth2.Token, calendarID string, opts *GoogleListEventsOption) ([]*Event, error) {
	uc.log.Debugf("ListCalendarEvents calendarID: %s", calendarID)
	return uc.repo.ListCalendarEvents(ctx, token, calendarID, opts)
}
