package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"

	pb "github.com/kdimtricp/aical/api/chat/v1"
)

type ChatService struct {
	uc  *biz.ChatUseCase
	guc *biz.GoogleUseCase
	uuc *biz.UserUseCase
	log *log.Helper
	pb.UnimplementedChatServer
}

func NewChatService(uc *biz.ChatUseCase, guc *biz.GoogleUseCase, uuc *biz.UserUseCase, logger log.Logger) *ChatService {
	return &ChatService{
		uc:  uc,
		guc: guc,
		uuc: uuc,
		log: log.NewHelper(logger),
	}
}

func (s *ChatService) UserChat(ctx context.Context, req *pb.UserChatRequest) (*pb.UserChatResponse, error) {
	s.log.Debugf("UserChat request: %v", req)
	user, err := s.uuc.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	token, err := s.guc.TokenSource(ctx, user.RefreshToken)
	if err != nil {
		s.log.Errorf("cron job:sync loop: get token failed: %v", err)
		return nil, err
	}
	ctx = biz.SetToken(ctx, token)
	answer, err := s.uc.UserChat(ctx, req.Question)
	if err != nil {
		return nil, err
	}
	r := &pb.UserChatResponse{Answer: answer}
	s.log.Debugf("UserChat reply: %v", r)
	return r, nil
}
