package biz

import (
	"context"
	"golang.org/x/oauth2"
)

const (
	TOKEN_KEY   = "token"
	CHAT_ID_KEY = "chat_id"
)

// SetToken returns context with token
func SetToken(ctx context.Context, token *oauth2.Token) context.Context {
	return context.WithValue(ctx, TOKEN_KEY, token)
}

// GetToken returns token from context
func GetToken(ctx context.Context) *oauth2.Token {
	return ctx.Value(TOKEN_KEY).(*oauth2.Token)
}

// SetChatID returns context with chat id
func SetChatID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, CHAT_ID_KEY, id)
}

// GetChatID returns chat id from context
func GetChatID(ctx context.Context) int64 {
	return ctx.Value(CHAT_ID_KEY).(int64)
}
