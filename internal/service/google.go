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

// GetAuthURL .
func (g *Google) AuthCodeURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// token .
func (g *Google) token(ctx context.Context, code string) (*oauth2.Token, error) {
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
	return token, nil
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

// GetUser .
func (g *Google) UserRegistration(ctx context.Context, code string) (*biz.User, error) {
	token, err := g.token(ctx, code)
	if err != nil {
		return nil, err
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

func (g *Google) CalendarInfo(ctx context.Context, token *oauth2.Token, name string) (*biz.Calendar, error) {
	if token.AccessToken == "" {
		return nil, errors.New("access token is empty")
	}
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	calendar, err := srv.Calendars.Get(name).Do()
	if err != nil {
		return nil, err
	}
	return &biz.Calendar{
		ID:          calendar.Id,
		Summary:     calendar.Summary,
		Description: calendar.Description,
	}, nil
}

// ListCalendars .
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

// SubscribeToCalendar starts watching for changes to a collection of events on a given calendar
func (g *Google) SubscribeToCalendar(ctx context.Context, token *oauth2.Token, calendarID string, webhookUrl string) error {
	if token.AccessToken == "" {
		return errors.New("access token is empty")
	}
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	channel := &calendarAPI.Channel{
		Address: webhookUrl,
		Id:      "unique-id",
		Type:    "web_hook",
	}
	_, err = srv.Events.Watch(calendarID, channel).Do()
	if err != nil {
		return err
	}
	return nil
}
