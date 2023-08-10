package service

import (
	"context"
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
) *CronService {
	return &CronService{
		c:    c,
		log:  log.NewHelper(log.With(logger, "module", "service/cron")),
		uuc:  uuc,
		cuc:  cuc,
		euc:  euc,
		ehuc: ehuc,
		guc:  guc,
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
		s.log.Debugf("cron job:sync loop: sync calendars for user: %v", user)
		// Get token from user refresh token
		token, err := s.guc.TokenSource(ctx, user.RefreshToken)
		if err != nil {
			s.log.Errorf("cron job:sync loop: get token failed: %v", err)
			return
		}
		// Sync calendars
		calendars, err := s.guc.ListUserCalendars(ctx, token)
		if err != nil {
			s.log.Errorf("cron job:sync loop: list user calendars failed: %v", err)
			return
		}
		if err := s.cuc.Sync(ctx, user.ID, calendars); err != nil {
			s.log.Errorf("cron job:sync loop: sync calendars failed: %v", err)
			return
		}
		var changes []*biz.EventHistory
		// Sync events
		thisWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
		nextWeek := time.Now().AddDate(0, 0, 14-int(time.Now().Weekday()))
		for _, calendar := range calendars {
			events, err := s.guc.ListCalendarEvents(ctx, token, calendar.GoogleID, &biz.GoogleListEventsOption{
				TimeMin: thisWeek.Format(time.RFC3339),
				TimeMax: nextWeek.Format(time.RFC3339),
			})
			if err != nil {
				s.log.Errorf("cron job:sync loop: list calendar events failed: %v", err)
				return
			}
			calendar, err := s.cuc.Get(ctx, calendar)
			if err != nil {
				s.log.Errorf("cron job:sync loop: get calendar failed: %v", err)
				return
			}
			if err := s.euc.Sync(ctx, calendar.ID, events); err != nil {
				s.log.Errorf("cron job:sync loop: sync events failed: %v", err)
				return
			}
			if eh, err := s.ehuc.ListCalendarEventHistory(ctx, calendar.ID); err != nil {
				s.log.Errorf("cron job:sync loop: list calendar event history failed: %v", err)
				return
			} else {
				changes = append(changes, eh...)
			}
		}
		if len(changes) > 0 {
			// TODO: Send changes to assistant
			//s.aiuc.

			// Delete event history after assistant is done
			for _, calendar := range calendars {
				if err := s.ehuc.DeleteCalendarEventHistory(ctx, calendar.ID); err != nil {
					s.log.Errorf("cron job:sync loop: delete calendar event history failed: %v", err)
					return
				}
			}
		}
	}
}
