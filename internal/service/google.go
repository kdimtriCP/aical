package service

import (
	"context"
	"errors"
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

// Google .
type Google struct {
	config *oauth2.Config
}

// NewGoogleService .
func NewGoogleService(c *conf.Google, logger log.Logger) (*Google, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the google resources")
	}
	return &Google{
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
func (g *Google) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// UserRegistration creates a new user from google oauth2
func (g *Google) UserRegistration(ctx context.Context, code string) (*biz.User, error) {
	if code == "" {
		return nil, errors.New("code is empty")
	}
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	if token.RefreshToken == "" {
		return nil, errors.New("refresh token is empty")
	}
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := oauth2API.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	userInfo, err := srv.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}
	return &biz.User{
		ID:           userInfo.Id,
		Name:         userInfo.Name,
		Email:        userInfo.Email,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (g *Google) TokenSource(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	if refreshToken == "" {
		return nil, errors.New("bad request")
	}
	t := &oauth2.Token{RefreshToken: refreshToken}
	token, err := g.config.TokenSource(ctx, t).Token()
	if err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, errors.New("access token is empty")
	}
	if token.RefreshToken == "" {
		return nil, errors.New("refresh token is empty")
	}
	return token, nil
}

// ListCalendars lists calendars from google calendar
func (g *Google) ListCalendars(ctx context.Context, token *oauth2.Token) (biz.Calendars, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	calendarList, err := srv.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	var calendars []*biz.Calendar
	for _, calendar := range calendarList.Items {
		calendars = append(calendars, &biz.Calendar{
			ID:          calendar.Id,
			Summary:     calendar.Summary,
			Description: calendar.Description,
		})
	}
	return calendars, nil
}

type GoogleListEventsOption struct {
	TimeMin           string
	TimeMax           string
	UpdatedMin        string
	MaxResults        int64
	OrderByUpdateTime bool
}

// listEventsCallWithOpts returns a call to list events
func (o *GoogleListEventsOption) listEventsCallWithOpts(ctx context.Context, call *calendarAPI.EventsListCall) *calendarAPI.EventsListCall {
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

// listEventsInstancesCallWithOpts returns a call to list events instances
func (o *GoogleListEventsOption) listEventsInstancesCallWithOpts(ctx context.Context, call *calendarAPI.EventsInstancesCall) *calendarAPI.EventsInstancesCall {
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

// ListEvents lists events from google calendar
func (g *Google) ListEvents(ctx context.Context, token *oauth2.Token, calendarID string, opts *GoogleListEventsOption) (biz.Events, error) {
	var eventsList []*biz.Event
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	call := srv.Events.List(calendarID).Context(ctx)
	if opts != nil {
		call = opts.listEventsCallWithOpts(ctx, call)
	}
	events, err := call.Do()
	if err != nil {
		return nil, err
	}
	// If event is recurring, then get all individual recurring events
	for _, event := range events.Items {
		if event.Recurrence != nil {
			call := srv.Events.Instances(calendarID, event.Id).Context(ctx)
			if opts != nil {
				call = opts.listEventsInstancesCallWithOpts(ctx, call)
			}
			recurrenceEvents, err := call.Do()
			if err != nil {
				return nil, err
			}
			for _, recurrenceEvent := range recurrenceEvents.Items {
				eventsList = append(eventsList, toBizEvent(recurrenceEvent, calendarID))
			}
		} else {
			eventsList = append(eventsList, toBizEvent(event, calendarID))
		}
	}
	return eventsList, nil
}

// toBizEvent converts a google calendar event to biz event
func toBizEvent(event *calendarAPI.Event, calendarID string) *biz.Event {
	var e biz.Event
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
	e.ID = event.Id
	e.Title = event.Summary
	e.Location = event.Location
	e.CalendarID = calendarID
	return &e
}

// QuickAddEvent creates a new event in google calendar from quickAdd string
func (g *Google) QuickAddEvent(ctx context.Context, token *oauth2.Token, calendarID string, quickAdd string) (*biz.Event, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	event, err := srv.Events.QuickAdd(calendarID, quickAdd).Do()
	if err != nil {
		return nil, err
	}
	return toBizEvent(event, calendarID), nil
}

// CreateEvent creates a new event in google calendar
func (g *Google) CreateEvent(ctx context.Context, token *oauth2.Token, event *biz.Event) (*biz.Event, error) {
	calendarID := event.CalendarID
	if calendarID == "" {
		calendarID = "primary"
	}
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	e := &calendarAPI.Event{
		Summary:   event.Title,
		Location:  event.Location,
		Start:     &calendarAPI.EventDateTime{DateTime: event.StartTime.Format(time.RFC3339)},
		End:       &calendarAPI.EventDateTime{DateTime: event.EndTime.Format(time.RFC3339)},
		Attendees: []*calendarAPI.EventAttendee{},
	}
	ev, err := srv.Events.Insert(calendarID, e).Do()
	if err != nil {
		return nil, err
	}
	return toBizEvent(ev, calendarID), nil
}

// CreateEvents creates new events in google calendar
func (g *Google) CreateEvents(ctx context.Context, token *oauth2.Token, events []*biz.Event) ([]*biz.Event, error) {
	var eventsList []*biz.Event
	for _, event := range events {
		e, err := g.CreateEvent(ctx, token, event)
		if err != nil {
			return nil, err
		}
		eventsList = append(eventsList, e)
	}
	return eventsList, nil
}
