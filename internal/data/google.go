package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"github.com/kdimtricp/aical/internal/conf"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendarAPI "google.golang.org/api/calendar/v3"
	oauth2API "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"time"
)

// GoogleRepo .
type GoogleRepo struct {
	config *oauth2.Config
}

// NewGoogleService .
func NewGoogleRepo(c *conf.Google, logger log.Logger) (biz.GoogleRepo, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the google resources")
	}
	return &GoogleRepo{
		config: &oauth2.Config{
			ClientID:     c.Client.Id,
			ClientSecret: c.Client.Secret,
			RedirectURL:  c.RedirectUrl,
			Scopes: []string{
				calendarAPI.CalendarScope,
				calendarAPI.CalendarEventsScope,
				oauth2API.UserinfoEmailScope,
				oauth2API.UserinfoProfileScope,
			},
			Endpoint: google.Endpoint,
		},
	}, cleanup, nil
}

// AuthCodeURL returns the url to redirect to google oauth2
func (g *GoogleRepo) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// TokenExchange returns a new oauth2 token from "Auth code"
func (g *GoogleRepo) TokenExchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

// TokenSource returns a new oauth2 token from "Refresh token"
func (g *GoogleRepo) TokenSource(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	t := &oauth2.Token{RefreshToken: refreshToken}
	return g.config.TokenSource(ctx, t).Token()
}

// UserInfo creates a new user from GoogleRepo oauth2
func (g *GoogleRepo) UserInfo(ctx context.Context, token *oauth2.Token) (*biz.User, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := oauth2API.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	ui, err := srv.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}
	return &biz.User{
		GoogleID: ui.Id,
		Email:    ui.Email,
		Name:     ui.Name,
	}, nil
}

// ListUserCalendars lists calendars from google calendar
func (g *GoogleRepo) ListUserCalendars(ctx context.Context, token *oauth2.Token) ([]*biz.Calendar, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	cals, err := srv.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	var calendars []*biz.Calendar
	for _, cal := range cals.Items {
		calendars = append(calendars, &biz.Calendar{
			GoogleID: cal.Id,
			Summary:  cal.Summary,
		})
	}
	return calendars, nil
}

// marshalEvent converts a biz.Event to a calendarAPI.Event
func marshalGoogleEvent(event *biz.Event) *calendarAPI.Event {
	return &calendarAPI.Event{
		Id:       event.GoogleID,
		Summary:  event.Summary,
		Location: event.Location,
		Start:    &calendarAPI.EventDateTime{DateTime: event.StartTime.Format(time.RFC3339)},
		End:      &calendarAPI.EventDateTime{DateTime: event.EndTime.Format(time.RFC3339)},
	}
}

// unmarshalGoogleEvent converts a calendarAPI.Event to a biz.Event
func unmarshalGoogleEvent(event *calendarAPI.Event) *biz.Event {
	var e biz.Event
	updated, err := time.Parse(time.RFC3339, event.Updated)
	if err == nil {
		e.UpdatedAt = updated
	}
	if startDate, err := time.Parse("2006-01-02", event.Start.Date); err == nil {
		e.StartTime = startDate
		e.IsAllDay = true
	}
	if startTime, err := time.Parse(time.RFC3339, event.Start.DateTime); err == nil {
		e.StartTime = startTime
		e.IsAllDay = false
	}
	if endDate, err := time.Parse("2006-01-02", event.End.Date); err == nil {
		e.EndTime = endDate
		e.IsAllDay = true
	}
	if endTime, err := time.Parse(time.RFC3339, event.End.DateTime); err == nil {
		e.EndTime = endTime
		e.IsAllDay = false
	}
	e.GoogleID = event.Id
	e.Summary = event.Summary
	e.Location = event.Location
	return &e
}

// ListEvents lists events from google calendar
func (g *GoogleRepo) ListCalendarEvents(ctx context.Context, token *oauth2.Token, calendarID string, opts *biz.GoogleListEventsOption) ([]*biz.Event, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	call := srv.Events.List(calendarID).Context(ctx)
	if opts != nil {
		call = opts.ListEventsCallWithOpts(ctx, call)
	}
	events, err := call.Do()
	if err != nil {
		return nil, err
	}
	var eventsList []*calendarAPI.Event
	// If event is recurring, then get all individual recurring events
	for _, event := range events.Items {
		if event.Recurrence != nil {
			call := srv.Events.Instances(calendarID, event.Id).Context(ctx)
			if opts != nil {
				call = opts.ListEventsInstancesCallWithOpts(ctx, call)
			}
			recurrenceEvents, err := call.Do()
			if err != nil {
				return nil, err
			}
			for _, recurrenceEvent := range recurrenceEvents.Items {
				eventsList = append(eventsList, recurrenceEvent)
			}
		} else {
			eventsList = append(eventsList, event)
		}
	}
	var bizEvents []*biz.Event
	for _, event := range eventsList {
		bizEvents = append(bizEvents, unmarshalGoogleEvent(event))
	}
	return bizEvents, nil
}

// CreateEvent creates a new event in google calendar
func (g *GoogleRepo) CreateCalendarEvent(ctx context.Context, token *oauth2.Token, event *biz.Event, calendarID string) (*biz.Event, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	e, err := srv.Events.Insert(calendarID, marshalGoogleEvent(event)).Do()
	if err != nil {
		return nil, err
	}
	return unmarshalGoogleEvent(e), nil
}

// UpdateEvent updates an event in google calendar
func (g *GoogleRepo) UpdateCalendarEvent(ctx context.Context, token *oauth2.Token, event *biz.Event, calendarID string) (*biz.Event, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	e, err := srv.Events.Update(calendarID, event.GoogleID, marshalGoogleEvent(event)).Do()
	if err != nil {
		return nil, err
	}
	return unmarshalGoogleEvent(e), nil
}

// DeleteEvent deletes an event in google calendar
func (g *GoogleRepo) DeleteCalendarEvent(ctx context.Context, token *oauth2.Token, event *biz.Event, calendarID string) error {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	return srv.Events.Delete(calendarID, event.GoogleID).Do()
}

// GetEvent gets an event in google calendar
func (g *GoogleRepo) GetCalendarEvent(ctx context.Context, token *oauth2.Token, event *biz.Event, calendarID string) (*biz.Event, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	e, err := srv.Events.Get(calendarID, event.GoogleID).Do()
	if err != nil {
		return nil, err
	}
	return unmarshalGoogleEvent(e), nil
}
