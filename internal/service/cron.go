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
	log   *log.Helper
	gg    *Google
}

var Jobs = map[string]func(){}

const EVENTS_LOOP_TIMEOUT = 40 * time.Second

func NewCronService(
	c *conf.Cron,
	logger log.Logger,
	uuc *biz.UserUseCase,
	cuc *biz.CalendarUseCase,
	gg *Google,
) *CronService {
	return &CronService{
		c:   c,
		log: log.NewHelper(log.With(logger, "module", "service/cron")),
		uuc: uuc,
		cuc: cuc,
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
	s.log.Debugf("cron job: events loop started")
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

		}
	}

}
