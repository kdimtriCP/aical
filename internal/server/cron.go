package server

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/kdimtricp/aical/internal/service"
	"github.com/robfig/cron/v3"
)

type CronServer struct {
	c   *conf.Cron
	crn *cron.Cron
	log *log.Helper
}

func NewCronServer(c *conf.Cron, logger log.Logger, srvs *service.CronService) (*CronServer, error) {
	srvs.Init()
	s := &CronServer{
		c:   c,
		crn: cron.New(),
		log: log.NewHelper(log.With(logger, "module", "server/cron")),
	}
	for _, job := range c.Jobs {
		j, ok := service.Jobs[job.Name]
		if !ok {
			s.log.Errorf("cron job: %s not found", job)
			return nil, errors.New("cron job not found")
		}
		if id, err := s.crn.AddFunc(job.Schedule, j); err != nil {
			s.log.Errorf("cron job: %s add failed: %v", job, err)
			return nil, err
		} else {
			s.log.Debugf("cron job: %s added, id: %d", job, id)
		}
	}
	return s, nil
}

func (s *CronServer) Start(_ context.Context) error {
	s.log.Debug("cron server: started")
	s.crn.Start()
	return nil
}

func (s *CronServer) Stop(_ context.Context) error {
	s.log.Debug("cron server: stopped")
	s.crn.Stop()
	return nil
}
