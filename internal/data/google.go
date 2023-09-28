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

// googleRepo .
type googleRepo struct {
	config *oauth2.Config
}

func NewGoogleRepo(c *conf.Google, logger log.Logger) (biz.GoogleRepo, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the google resources")
	}
	return &googleRepo{
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
func (g *googleRepo) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// TokenExchange returns a new oauth2 token from "Auth code"
func (g *googleRepo) TokenExchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

// TokenSource returns a new oauth2 token from "Refresh token"
func (g *googleRepo) TokenSource(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	t := &oauth2.Token{RefreshToken: refreshToken}
	return g.config.TokenSource(ctx, t).Token()
}

// UserInfo creates a new user from googleRepo oauth2
func (g *googleRepo) UserInfo(ctx context.Context, token *oauth2.Token) (*biz.User, error) {
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
func (g *googleRepo) ListUserCalendars(ctx context.Context, token *oauth2.Token) ([]*biz.Calendar, error) {
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

func (g *googleRepo) ListCalendarEvents(ctx context.Context, token *oauth2.Token, calendarID string, opts *biz.GoogleListEventsOption) ([]*biz.Event, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	call := srv.Events.List(calendarID).Context(ctx)
	if opts != nil {
		call = opts.ListEventsCallWithOpts(call)
	}
	events, err := call.Do()
	if err != nil {
		return nil, err
	}
	var eventsList []*calendarAPI.Event
	// If Event is recurring, then get all individual recurring events
	for _, event := range events.Items {
		if event.Recurrence != nil {
			call := srv.Events.Instances(calendarID, event.Id).Context(ctx)
			if opts != nil {
				call = opts.ListEventsInstancesCallWithOpts(call)
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

func (g *googleRepo) CreateCalendarEvent(ctx context.Context, token *oauth2.Token, event *biz.Event, calendarID string) (*biz.Event, error) {
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

func (g *googleRepo) UpdateCalendarEvent(ctx context.Context, token *oauth2.Token, event *biz.Event, calendarID string) (*biz.Event, error) {
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

func (g *googleRepo) DeleteCalendarEvent(ctx context.Context, token *oauth2.Token, event *biz.Event, calendarID string) error {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	return srv.Events.Delete(calendarID, event.GoogleID).Do()
}

func (g *googleRepo) GetCalendarEvent(ctx context.Context, token *oauth2.Token, event *biz.Event, calendarID string) (*biz.Event, error) {
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

// CreateNewCalendar creates a new calendar in google calendar
func (g *googleRepo) CreateNewCalendar(ctx context.Context, token *oauth2.Token, calendarName string) (*biz.Calendar, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	calendar := &calendarAPI.Calendar{
		Summary: calendarName,
	}
	c, err := srv.Calendars.Insert(calendar).Do()
	if err != nil {
		return nil, err
	}
	return &biz.Calendar{
		GoogleID: c.Id,
		Summary:  c.Summary,
	}, nil
}
