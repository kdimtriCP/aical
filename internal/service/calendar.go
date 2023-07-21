package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
	"net/http"

	pb "github.com/kdimtricp/aical/api/calendar/v1"
)

type CalendarService struct {
	pb.UnimplementedCalendarServiceServer
	uc  *biz.CalendarUseCase
	log *log.Helper
}

func NewCalendarService(uc *biz.CalendarUseCase, logger log.Logger) *CalendarService {
	return &CalendarService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *CalendarService) CreateCalendar(ctx context.Context, req *pb.CreateCalendarRequest) (*pb.CreateCalendarReply, error) {
	s.log.Debugf("create calendar for userID: %s", req.UserID)
	if req.UserID == "" {
		return nil, errors.New(http.StatusBadRequest, "userID is empty", "userID is empty")
	}
	calendar, err := s.uc.CreateCalendar(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pb.CreateCalendarReply{
		CalendarID: calendar.ID,
	}, nil
}
