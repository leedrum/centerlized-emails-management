package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/leedrum/centerlized-emails-management/email-indexer/internal/indexer"
)

// Kafka message structure example:
// { "tenant_id": "acme", "email_id": 123 }
type EmailIndexMessage struct {
	TenantID string `json:"tenant_id"`
	EmailID  int64  `json:"email_id"`
}

type KafkaConsumer struct {
	Svc        *indexer.Service
	DLQ        sarama.SyncProducer
	DLQTopic   string
	MaxRetries int
}

func (c *KafkaConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *KafkaConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		// Process + retry + DLQ
		if err := c.handleMessage(msg); err != nil {
			log.Printf("[DLQ] Failed after retries, sending to DLQ: %v", err)
			c.sendToDLQ(msg, err)
		}

		// Mark offset regardless of success or DLQ
		session.MarkMessage(msg, "")
	}

	return nil
}

// -----------------------------------------------------------------------------
// Core message handling with retry & DLQ
// -----------------------------------------------------------------------------
func (c *KafkaConsumer) handleMessage(msg *sarama.ConsumerMessage) error {

	var payload EmailIndexMessage
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		return err
	}

	var err error

	for attempt := 1; attempt <= c.MaxRetries; attempt++ {

		err = c.Svc.IndexEmailByID(context.Background(), payload.TenantID, payload.EmailID)
		if err == nil {
			// success
			log.Printf("[OK] Indexed email %d for tenant %s", payload.EmailID, payload.TenantID)
			return nil
		}

		log.Printf("[RETRY %d/%d] tenant=%s email=%d err=%v",
			attempt, c.MaxRetries, payload.TenantID, payload.EmailID, err)

		time.Sleep(time.Second * 2) // backoff
	}

	// Fail after max retries → DLQ
	return err
}

// -----------------------------------------------------------------------------
// DLQ publisher
// -----------------------------------------------------------------------------
func (c *KafkaConsumer) sendToDLQ(msg *sarama.ConsumerMessage, processingErr error) {

	raw := map[string]interface{}{
		"original_topic": msg.Topic,
		"partition":      msg.Partition,
		"offset":         msg.Offset,
		"payload":        string(msg.Value),
		"error":          processingErr.Error(),
		"timestamp":      time.Now().Unix(),
	}

	b, _ := json.Marshal(raw)

	dlqMsg := &sarama.ProducerMessage{
		Topic: c.DLQTopic,
		Value: sarama.ByteEncoder(b),
	}

	_, _, err := c.DLQ.SendMessage(dlqMsg)
	if err != nil {
		log.Printf("[DLQ ERROR] Failed to publish to DLQ: %v", err)
	} else {
		log.Printf("[DLQ] Sent failed message → %s", c.DLQTopic)
	}
}
