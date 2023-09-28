package biz

import (
	"context"
	"golang.org/x/oauth2"
)

//goland:noinspection ALL,GoUnnecessarilyExportedIdentifiers
const (
	TOKEN_KEY = "token"
)

// SetToken returns context with token
func SetToken(ctx context.Context, token *oauth2.Token) context.Context {
	return context.WithValue(ctx, TOKEN_KEY, token)
}

// GetToken returns token from context
func GetToken(ctx context.Context) *oauth2.Token {
	return ctx.Value(TOKEN_KEY).(*oauth2.Token)
}
