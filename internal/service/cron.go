package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"github.com/kdimtricp/aical/internal/conf"
	"time"
)

type CronService struct {
	c        *conf.Cron
	log      *log.Helper
	uuc      *biz.UserUseCase
	cuc      *biz.CalendarUseCase
	euc      *biz.EventUseCase
	ehuc     *biz.EventHistoryUseCase
	guc      *biz.GoogleUseCase
	aiuc     *biz.OpenAIUseCase
	lastSync time.Time
}

var Jobs = map[string]func(){}

const (
	SYNC_LOOP_TIMEOUT = 10 * time.Minute
)

func NewCronService(
	c *conf.Cron,
	logger log.Logger,
	uuc *biz.UserUseCase,
	cuc *biz.CalendarUseCase,
	euc *biz.EventUseCase,
	ehuc *biz.EventHistoryUseCase,
	guc *biz.GoogleUseCase,
	aiuc *biz.OpenAIUseCase,
) *CronService {
	return &CronService{
		c:    c,
		log:  log.NewHelper(log.With(logger, "module", "service/cron")),
		uuc:  uuc,
		cuc:  cuc,
		euc:  euc,
		ehuc: ehuc,
		guc:  guc,
		aiuc: aiuc,
	}
}

// Init initializes the cron service.
func (s *CronService) Init() {
	if len(s.c.Jobs) == 0 {
		return
	}
	Jobs[s.c.Jobs[0].Name] = s.SyncLoop
}

// SyncLoop .
func (s *CronService) SyncLoop() {
	syncStart := time.Now()
	s.log.Debugf("cron job:sync loop: start at %s", syncStart.Format(time.RFC3339))
	defer func() {
		s.lastSync = time.Now()
		s.log.Debugf("cron job:sync loop: end at %s, duration: %s", s.lastSync.Format(time.RFC3339), time.Since(syncStart))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), SYNC_LOOP_TIMEOUT)
	defer cancel()

	// List users from database
	users, err := s.uuc.List(ctx)
	if err != nil {
		s.log.Errorf("cron job:sync loop: list users failed: %v", err)
		return
	}
	for _, user := range users {
		token, err := s.guc.TokenSource(ctx, user.RefreshToken)
		if err != nil {
			s.log.Errorf("cron job:sync loop: get token failed: %v", err)
			return
		}
		ctx = biz.SetToken(ctx, token)
		if err := s.SyncUserCalendars(ctx, user); err != nil {
			s.log.Errorf("cron job:sync loop: sync user events failed: %v", err)
			return
		}
		calendars, err := s.cuc.ListUserCalendars(ctx, user.ID)
		if err != nil {
			s.log.Errorf("cron job:sync loop: list calendars failed: %v", err)
			return
		}
		for _, calendar := range calendars {
			if err := s.SyncCalendarEvents(ctx, calendar); err != nil {
				s.log.Errorf("cron job:sync loop: sync calendar events failed: %v", err)
				return
			}
			/*
				if err := s.GenerateCalendarEvents(ctx, calendar); err != nil {
					s.log.Errorf("cron job:sync loop: generate calendar events failed: %v", err)
					return
				}
			*/
		}
	}
	return
}

// SyncUserCalendars
func (s *CronService) SyncUserCalendars(ctx context.Context, user *biz.User) error {
	s.log.Debugf("cron job:sync loop: sync calendars for user: %v", user)
	// Get token from context
	token := biz.GetToken(ctx)
	if token == nil {
		s.log.Errorf("cron job:sync loop: token not found in context")
		return fmt.Errorf("token not found in context")
	}
	// Sync calendars
	calendars, err := s.guc.ListUserCalendars(ctx, token)
	if err != nil {
		s.log.Errorf("cron job:sync loop: list user calendars failed: %v", err)
		return err
	}
	if err := s.cuc.Sync(ctx, user.ID, calendars); err != nil {
		s.log.Errorf("cron job:sync loop: sync calendars failed: %v", err)
		return err
	}
	return nil
}

// SyncCalendarEvents
func (s *CronService) SyncCalendarEvents(ctx context.Context, calendar *biz.Calendar) error {
	s.log.Debugf("cron job:sync loop: sync events for calendar: %v", calendar)
	// Get token from context
	token := biz.GetToken(ctx)
	if token == nil {
		s.log.Errorf("cron job:sync loop: token not found in context")
		return fmt.Errorf("token not found in context")
	}
	events, err := s.guc.ListCalendarEvents(ctx, token, calendar.GoogleID, &biz.GoogleListEventsOption{
		TimeMin: time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1).Format(time.RFC3339), // this week
		TimeMax: time.Now().AddDate(0, 0, 14-int(time.Now().Weekday())).Format(time.RFC3339), // next week
	})
	if err != nil {
		s.log.Errorf("cron job:sync loop: list calendar events failed: %v", err)
		return err
	}
	if err := s.euc.Sync(ctx, calendar.ID, events); err != nil {
		s.log.Errorf("cron job:sync loop: sync events failed: %v", err)
		return err
	}
	return nil
}

// GenerateCalendarEvents
func (s *CronService) GenerateCalendarEvents(ctx context.Context, calendar *biz.Calendar) error {
	s.log.Debugf("cron job:sync loop: generate events for calendar: %v", calendar)
	// Get token from user refresh token

	// Get events from database
	events, err := s.euc.List(ctx, calendar.ID)
	if err != nil {
		s.log.Errorf("cron job:sync loop: list events failed: %v", err)
		return err
	}
	// Generate events
	if err := s.aiuc.GenerateCalendarEvents(ctx, calendar, events); err != nil {
		s.log.Errorf("cron job:sync loop: generate events failed: %v", err)
		return err
	}
	return nil
}
