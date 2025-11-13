package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/leedrum/centerlized-emails-management/mail-fetcher/config"
	"github.com/leedrum/centerlized-emails-management/mail-fetcher/db"
	sqlc "github.com/leedrum/centerlized-emails-management/mail-fetcher/db/sqlc"
	"github.com/leedrum/centerlized-emails-management/mail-fetcher/kafka"
	"golang.org/x/oauth2"
)

func main() {
	cfg := config.LoadConfig()
	pg := db.NewPostgresDB()
	q := sqlc.New(pg)
	tenantList := strings.Split(cfg.Tenants, ",")
	// TODO: get all the tenant from db if the config tenants is 'all'
	for _, tenant := range tenantList {
		gmailCredentials, err := q.GetGmailCredentialsByTenantName(
			context.Background(),
			tenant,
		)

		if err != nil {
			log.Default().Println(fmt.Sprintf("Can not find the credential for tenant: %s", tenant))
		}

		for _, credential := range gmailCredentials {
			token := oauth2.Token{
				AccessToken:  credential.AccessToken,
				RefreshToken: credential.RefreshToken,
				Expiry:       credential.TokenExpiry,
			}

			go kafka.StartFetchingEmailConsumer(tenant, cfg, token)
		}

		go kafka.StartTokenUpdateConsumer(cfg.KafkaBroker, "mail_fetcher.token_update", q)

	}

	select {} // block forever
}
