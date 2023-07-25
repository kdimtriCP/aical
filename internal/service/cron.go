package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"github.com/kdimtricp/aical/internal/conf"
	"sync"
	"time"
)

type CronService struct {
	c     *conf.Cron
	mutex sync.Mutex
	uuc   *biz.UserUseCase
	cuc   *biz.CalendarUseCase
	euc   *biz.EventUseCase
	log   *log.Helper
	gg    *Google
}

var Jobs = map[string]func(){}

const (
	EVENTS_LOOP_TIMEOUT   = 10 * time.Minute
	EVENTS_MIN_START_TIME = -time.Hour * 24 * 7
	EVENTS_MAX_START_TIME = time.Hour * 24
)

func NewCronService(
	c *conf.Cron,
	logger log.Logger,
	uuc *biz.UserUseCase,
	cuc *biz.CalendarUseCase,
	euc *biz.EventUseCase,
	gg *Google,
) *CronService {
	return &CronService{
		c:   c,
		log: log.NewHelper(log.With(logger, "module", "service/cron")),
		uuc: uuc,
		cuc: cuc,
		euc: euc,
		gg:  gg,
	}
}

// Init initializes the cron service.
func (s *CronService) Init() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.c.Jobs) == 0 {
		return
	}
	Jobs[s.c.Jobs[0].Name] = s.EventsLoop
}

// EventsLoop .
func (s *CronService) EventsLoop() {
	syncStart := time.Now()
	s.log.Debugf("cron job: events loop started at %v", syncStart.Format(time.RFC3339))
	defer func() {
		s.log.Debugf("cron job: events loop finished at %v, took %v", time.Now().Format(time.RFC3339), time.Since(syncStart))
	}()
	ctx, cancel := context.WithTimeout(context.Background(), EVENTS_LOOP_TIMEOUT)
	defer cancel()
	users, err := s.uuc.List(ctx)
	if err != nil {
		s.log.Errorf("cron job: list users failed: %v", err)
		return
	}
	for _, user := range users {
		s.log.Debugf("cron job: sync calendars for user: %v", user)
		token, err := s.gg.TokenSource(ctx, user.RefreshToken)
		if err != nil {
			s.log.Errorf("cron job: get token failed: %v", err)
			return
		}
		// List google calendars
		cals, err := s.gg.ListCalendars(ctx, token)
		if err != nil {
			s.log.Errorf("cron job: list google calendars failed: %v", err)
			return
		}
		// Sync google calendars with database
		if err := s.cuc.Sync(ctx, user.ID, cals); err != nil {
			s.log.Errorf("cron job: sync calendars failed: %v", err)
			return
		}
		for _, cal := range cals {
			s.log.Debugf("cron job: sync events for calendar: %v", cal)
			// List a week of Google calendar events
			evs, err := s.gg.ListEvents(ctx, token, cal.ID, &GoogleListEventsOption{
				TimeMin: syncStart.Add(EVENTS_MIN_START_TIME),
				TimeMax: syncStart.Add(EVENTS_MAX_START_TIME),
			})
			if err != nil {
				s.log.Errorf("cron job: list google events failed: %v", err)
				return
			}
			// Sync google calendar events with database
			if err := s.euc.Sync(ctx, cal.ID, evs); err != nil {
				s.log.Errorf("cron job: sync events failed: %v", err)
				return
			}

		}
	}
}
