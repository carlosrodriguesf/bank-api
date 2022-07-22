package model

import (
	"context"
	"time"
)

type (
	Credentials struct {
		Document string `json:"document" validate:"required"`
		Secret   string `json:"secret" validate:"required"`
	}
	Session struct {
		Token     string    `json:"token"`
		Account   Account   `json:"account"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

func SetSessionOnContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, "session", session)
}

func GetSessionFromContext(ctx context.Context) *Session {
	session, ok := ctx.Value("session").(*Session)
	if !ok {
		return nil
	}
	return session
}
