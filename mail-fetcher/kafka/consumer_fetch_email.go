package kafka

import (
	"log"
	"time"

	"github.com/leedrum/centerlized-emails-management/mail-fetcher/config"
	"github.com/leedrum/centerlized-emails-management/mail-fetcher/gmail"
	"github.com/leedrum/centerlized-emails-management/mail-fetcher/s3"
	"golang.org/x/oauth2"
)

func StartFetchingEmailConsumer(tenant string, cfg *config.Config, tokenOAuth2 oauth2.Token) {
	log.Printf("[Tenant %s] Worker started", tenant)

	gmailSvc := gmail.NewGmailClient(
		cfg.GmailClientID,
		cfg.GmailClientSecret,
		tokenOAuth2,
	)
	fetcher := gmail.NewFetcher(gmailSvc)
	uploader := s3.NewS3Uploader(cfg.S3Bucket)
	producer, err := NewProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	if err != nil {
		log.Printf("Failed to create producer: %v", tenant, err)
		return
	}

	for {
		files, err := fetcher.FetchLatestEmails("me", 10)
		if err != nil {
			log.Printf("[Tenant %s] Error fetching emails: %v", tenant, err)
			time.Sleep(1 * time.Minute)
			continue
		}

		for _, file := range files {
			s3Key := tenant + "/emails/" + file
			if err := uploader.UploadFile(file, s3Key); err != nil {
				log.Printf("[Tenant %s] Upload failed: %v", tenant, err)
				continue
			}
			PublishFetchedEmailProducer(producer, tenant, s3Key, file)
		}

		time.Sleep(5 * time.Minute)
	}
}
