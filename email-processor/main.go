package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/leedrum/centerlized-emails-management/email-processor/internal/config"
	"github.com/leedrum/centerlized-emails-management/email-processor/internal/db"
	"github.com/leedrum/centerlized-emails-management/email-processor/internal/kafka"

	"github.com/jhillyerd/enmime"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioClient  *minio.Client
	pgDB         = db.NewPostgresDB()
	kafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
)

func initMinio() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("failed to init minio: %v", err)
	}
	minioClient = client
}

func resolveThread(ctx context.Context, db *sql.DB, tenant string, messageID, inReplyTo, refs, subject string) (int64, error) {
	tableEmails := fmt.Sprintf(`"%s".emails`, tenant)
	tableThreads := fmt.Sprintf(`"%s".threads`, tenant)

	// 1. If replying to another email → use its thread
	if inReplyTo != "" {
		var threadID int64
		err := db.QueryRowContext(ctx,
			fmt.Sprintf(`SELECT thread_id FROM %s WHERE message_id = $1 LIMIT 1;`, tableEmails),
			inReplyTo,
		).Scan(&threadID)

		if err == nil {
			return threadID, nil
		}
	}

	// 2. Try all references from oldest → newest
	if refs != "" {
		ids := strings.Fields(refs)
		for _, ref := range ids {
			var threadID int64
			err := db.QueryRowContext(ctx,
				fmt.Sprintf(`SELECT thread_id FROM %s WHERE message_id = $1 LIMIT 1;`, tableEmails),
				ref,
			).Scan(&threadID)

			if err == nil {
				return threadID, nil
			}
		}
	}

	// 3. No parent found → create a new thread
	var newThreadID int64
	err := db.QueryRowContext(ctx,
		fmt.Sprintf(`INSERT INTO %s (subject) VALUES ($1) RETURNING id;`, tableThreads),
		subject,
	).Scan(&newThreadID)

	if err != nil {
		return 0, fmt.Errorf("create thread: %v", err)
	}

	return newThreadID, nil
}

func handleEmailMessage(ctx context.Context, data map[string]interface{}) error {
	tenant := data["tenant_id"].(string)
	s3Key := data["s3_key"].(string)
	file := data["file"].(string)

	// Ensure tenant schema
	if err := db.EnsureTenantSchema(pgDB, tenant); err != nil {
		return err
	}

	// Download EML file from S3/MinIO
	tmpFile := "/tmp/" + filepath.Base(file)
	err := minioClient.FGetObject(ctx, "emails", s3Key, tmpFile, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer os.Remove(tmpFile)

	r, err := os.Open(tmpFile)
	if err != nil {
		return err
	}
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return fmt.Errorf("failed to parse eml: %v", err)
	}

	m, _ := mail.ReadMessage(bytes.NewReader(env.Root.Content))
	from, _ := mail.ParseAddress(m.Header.Get("From"))
	to, _ := mail.ParseAddress(m.Header.Get("To"))
	subject := env.GetHeader("Subject")
	messageID := env.GetHeader("Message-Id")
	inReplyTo := env.GetHeader("In-Reply-To")
	refs := env.GetHeader("References")

	// Resolve thread
	threadID, err := resolveThread(
		ctx,
		pgDB,
		tenant,
		messageID,
		inReplyTo,
		refs,
		subject,
	)
	if err != nil {
		return fmt.Errorf("thread resolution error: %v", err)
	}

	attachmentsInfo := []map[string]string{}

	for _, a := range env.Attachments {
		objectName := fmt.Sprintf("%s/%d_%s", tenant, time.Now().UnixNano(), a.FileName)
		_, err := minioClient.PutObject(ctx, "attachments", objectName,
			bytes.NewReader(a.Content), int64(len(a.Content)),
			minio.PutObjectOptions{ContentType: a.ContentType})
		if err != nil {
			log.Printf("failed to upload attachment: %v", err)
			continue
		}
		attachmentsInfo = append(attachmentsInfo, map[string]string{
			"name": a.FileName,
			"s3":   objectName,
		})
	}

	attachmentsJSON, _ := json.Marshal(attachmentsInfo)

	// Insert into tenant schema
	table := fmt.Sprintf(`"%s".emails`, tenant)
	query := fmt.Sprintf(`
		INSERT INTO %s
		(subject, sender_name, sender_email, receiver_name, receiver_email,
		raw_body, text_body, attachments, sent_at, received_at, s3_key,
		message_id, in_reply_to, references, thread_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id;
	`, table)

	var id int64
	err = pgDB.QueryRowContext(ctx, query,
		subject,
		from.Name, from.Address,
		to.Name, to.Address,
		env.HTML,
		env.Text,
		string(attachmentsJSON),
		time.Now(),
		time.Now(),
		s3Key,
		messageID,
		inReplyTo,
		refs,
		threadID,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("insert email: %v", err)
	}

	// Produce message for indexer
	prod := kafka.NewProducer(kafkaBrokers)
	return prod.Publish("email_indexer.new_email", map[string]interface{}{
		"tenant_id": tenant,
		"email_id":  id,
	})
}

func main() {
	initMinio()

	handler := func(ctx context.Context, msg map[string]interface{}) error {
		log.Printf("Processing email for tenant=%v", msg["tenant_id"])
		return handleEmailMessage(ctx, msg)
	}

	consumer := kafka.NewConsumer(kafkaBrokers, "email-processor-group", handler)
	log.Println("🚀 Email Processor started...")

	ctx := context.Background()
	cfg := config.LoadConfig()
	topics := []string{"email_processor.fetched_email"}
	if cfg.KafkaTopic != "" {
		topics = strings.Split(cfg.KafkaTopic, ",") // scale-up for specific tenant
	}
	consumer.Start(ctx, topics)
}
