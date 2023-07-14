package data

import (
	"context"
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
	return g.config.AuthCodeURL(state)
}

// GetToken .
func (g *Google) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}

// GetUserInfo .
func (g *Google) GetUserInfo(ctx context.Context, token *oauth2.Token) (*oauth2API.Userinfo, error) {
	client := oauth2.NewClient(ctx, g.config.TokenSource(ctx, token))
	srv, err := oauth2API.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return srv.Userinfo.Get().Do()
}
