package data

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
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

// NewGoogle .
func NewGoogle(c *conf.Google, logger log.Logger) (*Google, func(), error) {
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
func (g *Google) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// GetToken .
func (g *Google) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("code is empty")
	}
	return g.config.Exchange(ctx, code)
}

func (g *Google) RefreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token is empty")
	}
	token := &oauth2.Token{RefreshToken: refreshToken}
	return g.config.TokenSource(ctx, token).Token()
}

// GetUserInfo .
func (g *Google) GetUserInfo(ctx context.Context, token *oauth2.Token) (*User, error) {
	if token.AccessToken == "" {
		return nil, errors.New("access token is empty")
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
	return &User{
		ID:    userInfo.Id,
		Name:  userInfo.Name,
		Email: userInfo.Email,
	}, nil
}

// GetCalendar .
func (g *Google) CalendarInfo(ctx context.Context, token *oauth2.Token) (*Calendar, error) {
	if token.AccessToken == "" {
		return nil, errors.New("access token is empty")
	}
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := calendarAPI.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	calendar, err := srv.Calendars.Get("primary").Do()
	if err != nil {
		return nil, err
	}
	return &Calendar{
		ID:          calendar.Id,
		Summary:     calendar.Summary,
		Description: calendar.Description,
	}, nil
}
