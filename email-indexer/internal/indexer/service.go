package indexer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	dbpkg "github.com/leedrum/centerlized-emails-management/email-indexer/internal/db" // sqlc-generated package
)

type Service struct {
	db *sql.DB
	q  *dbpkg.Queries
	es *elasticsearch.Client
}

func NewService(dbConn *sql.DB, q *dbpkg.Queries, es *elasticsearch.Client) *Service {
	return &Service{
		db: dbConn,
		q:  q,
		es: es,
	}
}

// IndexEmailByID fetches email from tenant schema and indexes into ES
func (s *Service) IndexEmailByID(ctx context.Context, tenant string, emailID int64) error {
	// validate tenant string to avoid SQL injection
	if !isValidTenant(tenant) {
		return fmt.Errorf("invalid tenant name")
	}

	// start a transaction and set search_path
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// set search_path for this transaction
	if _, err := tx.ExecContext(ctx, fmt.Sprintf(`SET LOCAL search_path = "%s"`, tenant)); err != nil {
		return fmt.Errorf("set search_path: %w", err)
	}

	// use sqlc queries bound to this tx
	qtx := dbpkg.New(tx)

	// get email by id
	emailRec, err := qtx.GetEmailByID(ctx, emailID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("email id %d not found for tenant %s", emailID, tenant)
		}
		return err
	}

	var attachList []string

	if emailRec.Attachments.Valid && len(emailRec.Attachments.RawMessage) > 0 {
		if err := json.Unmarshal(emailRec.Attachments.RawMessage, &attachList); err != nil {
			log.Printf("failed to unmarshal attachments: %v", err)
			attachList = []string{}
		}
	} else {
		attachList = []string{}
	}

	// prepare document for ES
	doc := map[string]interface{}{
		"id":             emailRec.ID,
		"message_id":     emailRec.MessageID,
		"subject":        emailRec.Subject,
		"sender_name":    emailRec.SenderName,
		"sender_email":   emailRec.SenderEmail,
		"receiver_name":  emailRec.ReceiverName,
		"receiver_email": emailRec.ReceiverEmail,
		"raw_body":       emailRec.RawBody,
		"text_body":      emailRec.TextBody,
		"attachments":    attachList,
		"sent_at":        emailRec.SentAt,
		"received_at":    emailRec.ReceivedAt,
		"s3_key":         emailRec.S3Key,
		"thread_id":      emailRec.ThreadID,
	}

	// index name - sanitize tenant to alphanum + underscore
	idx := fmt.Sprintf("tenant_%s_emails", sanitizeTenant(tenant))

	// ensure index exists (create if not)
	if err := ensureIndex(ctx, s.es, idx); err != nil {
		return fmt.Errorf("ensure index: %w", err)
	}

	// index doc with ID = emailID
	body, _ := json.Marshal(doc)
	res, err := s.es.Index(
		idx,
		strings.NewReader(string(body)),
		s.es.Index.WithDocumentID(fmt.Sprintf("%d", emailRec.ID)),
		s.es.Index.WithContext(ctx),
		s.es.Index.WithRefresh("wait_for"), // optional
	)
	if err != nil {
		return fmt.Errorf("es index error: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("es index response error: %s", res.String())
	}

	// commit tx
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func isValidTenant(t string) bool {
	for _, r := range t {
		if !(r == '_' || (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '-') {
			return false
		}
	}
	return true
}

func sanitizeTenant(t string) string {
	// replace non-alnum with underscore and lower-case
	var b strings.Builder
	for _, r := range strings.ToLower(t) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}

// ensureIndex is a small helper that creates index if not exists.
// For production you should add mappings and settings.
func ensureIndex(ctx context.Context, es *elasticsearch.Client, index string) error {
	// check exists
	existsRes, err := es.Indices.Exists([]string{index})
	if err != nil {
		return err
	}
	defer existsRes.Body.Close()
	if existsRes.StatusCode == 200 {
		return nil
	}
	// create
	createRes, err := es.Indices.Create(index)
	if err != nil {
		return err
	}
	defer createRes.Body.Close()
	if createRes.IsError() {
		return fmt.Errorf("create index error: %s", createRes.String())
	}
	return nil
}
