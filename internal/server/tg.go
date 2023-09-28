package server

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/kdimtricp/aical/internal/service"
	"strings"
)

type TGServer struct {
	log  *log.Helper
	bot  *tgbotapi.BotAPI
	auth *service.AuthService
	chat *service.ChatService
}

func NewTGServer(c *conf.Server, logger log.Logger, _ *service.TGService, auth *service.AuthService, chat *service.ChatService) (*TGServer, error) {
	bot, err := tgbotapi.NewBotAPI(c.Tg.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = false

	return &TGServer{
		log:  log.NewHelper(log.With(logger, "module", "server/tgs")),
		bot:  bot,
		auth: auth,
		chat: chat,
	}, nil
}

func (s *TGServer) Start(ctx context.Context) error {
	s.log.Info("tgs server: started")
	uc := s.bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:         0,
		Limit:          0,
		Timeout:        60,
		AllowedUpdates: nil,
	})
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return nil
		// receive update from channel and then handle it
		case update := <-uc:
			s.handleUpdate(ctx, update)
		}
	}
}

func (s *TGServer) Stop(_ context.Context) error {
	s.log.Info("tgs server: stopped")
	return nil
}

func (s *TGServer) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		s.handleMessage(ctx, update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		s.handleButton(update.CallbackQuery)
		break
	}
}

func (s *TGServer) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	s.log.Infof("Message: %s", message.Text)
	if strings.HasPrefix(message.Text, "/login") {
		url, err := s.auth.AuthWithID(ctx, message.From.ID)
		if err != nil {
			s.log.Errorf("tg user auth error: %s", err.Error())
		}
		if _, err := s.bot.Send(tgbotapi.NewMessage(message.Chat.ID, url)); err != nil {
			s.log.Errorf("sending url for user login  error: %s", err.Error())
		}
	} else {
		answer, err := s.chat.TGChat(ctx, fmt.Sprintf("%d", message.From.ID), message.Text)
		if err != nil {
			s.log.Errorf("getting tg chat answer error: %s", err.Error())
		}
		if _, err := s.bot.Send(tgbotapi.NewMessage(message.Chat.ID, answer)); err != nil {
			s.log.Errorf("sending tg chat answer error: %s,", err.Error())
		}
	}
}

func (s *TGServer) handleButton(callback *tgbotapi.CallbackQuery) {
	s.log.Infof("Button: %s", callback.Data)
}

func (s *TGServer) handleCommand(command string) error {
	var err error

	switch command {
	case "/login":
		break
	default:
		s.log.Infof("Unknown command: %s", command)
	}
	return err
}
