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
	log   *log.Helper
	gg    *Google
	ai    *OpenAI
	uuc   *biz.UserUseCase
	cuc   *biz.CalendarUseCase
	euc   *biz.EventUseCase
}

var Jobs = map[string]func(){}

const (
	SYNC_LOOP_TIMEOUT = 10 * time.Minute
	GEN_LOOP_TIMEOUT  = 10 * time.Minute
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
	Jobs[s.c.Jobs[0].Name] = s.SyncLoop
	Jobs[s.c.Jobs[1].Name] = s.GenerateLoop
}

// SyncLoop .
func (s *CronService) SyncLoop() {
	if !s.mutex.TryLock() {
		s.log.Debugf("cron job:sync loop: another job is running")
		return
	}
	defer s.mutex.Unlock()
	syncStart := time.Now()
	s.log.Debugf("cron job:sync loop: started at %v", syncStart.Format(time.RFC3339))
	defer func() {
		s.log.Debugf("cron job:sync loop: finished at %v, took %v", time.Now().Format(time.RFC3339), time.Since(syncStart))
	}()
	ctx, cancel := context.WithTimeout(context.Background(), SYNC_LOOP_TIMEOUT)
	defer cancel()
	users, err := s.uuc.List(ctx)
	if err != nil {
		s.log.Errorf("cron job:sync loop: list users failed: %v", err)
		return
	}
	for _, user := range users {
		s.log.Debugf("cron job:sync loop: sync calendars for user: %v", user)
		token, err := s.gg.TokenSource(ctx, user.RefreshToken)
		if err != nil {
			s.log.Errorf("cron job:sync loop: get token failed: %v", err)
			return
		}
		// List google calendars
		cals, err := s.gg.ListCalendars(ctx, token)
		if err != nil {
			s.log.Errorf("cron job:sync loop: list google calendars failed: %v", err)
			return
		}
		// Sync google calendars with database
		if err := s.cuc.Sync(ctx, user.ID, cals); err != nil {
			s.log.Errorf("cron job:sync loop: sync calendars failed: %v", err)
			return
		}
		for _, cal := range cals {
			s.log.Debugf("cron job:sync loop: sync events for calendar: %v", cal)
			// List two weeks of Google calendar events
			// starting monday of this week and ending sunday of next week
			evs, err := s.gg.ListEvents(ctx, token, cal.ID, &GoogleListEventsOption{
				TimeMin: time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1).Format(time.RFC3339),
				TimeMax: time.Now().AddDate(0, 0, 14-int(time.Now().Weekday())).Format(time.RFC3339),
			})
			if err != nil {
				s.log.Errorf("cron job:sync loop: list google events failed: %v", err)
				return
			}
			// Sync google calendar events with database
			if err := s.euc.Sync(ctx, cal.ID, evs); err != nil {
				s.log.Errorf("cron job:sync loop: sync events failed: %v", err)
				return
			}
		}
	}
}

// GenerateLoop .
func (s *CronService) GenerateLoop() {
	if !s.mutex.TryLock() {
		s.log.Debugf("cron job:generate loop: another job is running")
		return
	}
	defer s.mutex.Unlock()
	genStart := time.Now()
	s.log.Debugf("cron job:generate loop: started at %v", genStart.Format(time.RFC3339))
	defer func() {
		s.log.Debugf("cron job:generate loop: finished at %v, took %v", time.Now().Format(time.RFC3339), time.Since(genStart))
	}()
	ctx, cancel := context.WithTimeout(context.Background(), GEN_LOOP_TIMEOUT)
	defer cancel()
	// List db users
	users, err := s.uuc.List(ctx)
	if err != nil {
		s.log.Errorf("cron job:generate loop: list users failed: %v", err)
		return
	}
	for _, user := range users {
		// List db calendars
		cals, err := s.cuc.List(ctx, user.ID)
		if err != nil {
			s.log.Errorf("cron job:generate loop: list google calendars failed: %v", err)
			return
		}
		var events biz.Events
		for _, cal := range cals {
			// List db events
			// not older than monday of this week
			// not newer than sunday of next week
			evs, err := s.euc.List(ctx, cal.ID, &biz.ListEventsOptions{
				IsUsed:   true,
				StartMin: time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1),
				StartMax: time.Now().AddDate(0, 0, 14-int(time.Now().Weekday())),
			})
			if err != nil {
				s.log.Errorf("cron job:generate loop: list google events failed: %v", err)
				return
			}
			events = append(events, evs...)
		}
		if len(events.FilterUsed()) <= 0 {
			s.log.Debugf("cron job:generate loop: no new events for user: %v", user)
			return
		}
		s.log.Debugf("cron job:generate loop: generate events for user: %v", user)
		// Generate events
		aievs, err := s.ai.GenerateNextWeekEvents(ctx, events)
		if err != nil {
			s.log.Errorf("cron job:generate loop: generate events failed: %v", err)
			return
		}
		// CreateAll new events in google calendar
		token, err := s.gg.TokenSource(ctx, user.RefreshToken)
		if err != nil {
			s.log.Errorf("cron job:sync loop: get token failed: %v", err)
			return
		}
		nevs, err := s.gg.CreateEvents(ctx, token, aievs)
		if err != nil {
			s.log.Errorf("cron job:generate loop: create events failed: %v", err)
			return
		}
		// CreateAll new events in database
		if err := s.euc.CreateAll(ctx, nevs); err != nil {
			s.log.Errorf("cron job:generate loop: create events failed: %v", err)
			return
		}
		// Mark all events as used
		events = append(events, nevs...)
		if err := s.euc.MarkAllUsed(ctx, events); err != nil {
			s.log.Errorf("cron job:generate loop: mark events as used failed: %v", err)
			return
		}
	}
}
