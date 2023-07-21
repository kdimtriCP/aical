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
}

var Jobs = map[string]func(){}

const EVENTS_LOOP_TIMEOUT = 40 * time.Second

func NewCronService(c *conf.Cron, logger log.Logger, uuc *biz.UserUseCase, cuc *biz.CalendarUseCase) *CronService {
	return &CronService{
		c:   c,
		uuc: uuc,
		cuc: cuc,
		log: log.NewHelper(log.With(logger, "module", "service/cron")),
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
	users, err := s.uuc.ListUsers(ctx)
	if err != nil {
		s.log.Errorf("cron job: list users failed: %v", err)
		return
	}
	for _, user := range users {
		s.log.Debugf("cron job: user id: %s", user.ID)
	}

}
