package db

import (
	"context"
	"time"

	sqlc "github.com/leedrum/centerlized-emails-management/mail-fetcher/db/sqlc"
)

type TokenUpdateEvent struct {
	TenantName   string
	UserEmail    string
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

func UpsertToken(q *sqlc.Queries, evt TokenUpdateEvent) error {
	ctx := context.Background()
	return q.UpsertToken(ctx, sqlc.UpsertTokenParams{
		TenantName:   evt.TenantName,
		Email:        evt.UserEmail,
		ClientID:     evt.ClientID,
		ClientSecret: evt.ClientSecret,
		AccessToken:  evt.AccessToken,
		RefreshToken: evt.RefreshToken,
		TokenExpiry:  evt.Expiry,
	})
}
