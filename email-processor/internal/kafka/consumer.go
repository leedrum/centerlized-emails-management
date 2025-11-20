package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type MessageHandler func(ctx context.Context, msg map[string]interface{}) error

type Consumer struct {
	group   sarama.ConsumerGroup
	handler MessageHandler
}

func NewConsumer(brokers []string, groupID string, handler MessageHandler) *Consumer {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_5_0_0
	group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		log.Fatalf("failed to create consumer group: %v", err)
	}
	return &Consumer{group: group, handler: handler}
}

func (c *Consumer) Start(ctx context.Context, topics []string) {
	for {
		err := c.group.Consume(ctx, topics, c)
		if err != nil {
			log.Printf("consumer error: %v", err)
		}
	}
}

// sarama.ConsumerGroupHandler implementation
func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Value, &data); err == nil {
			c.handler(context.Background(), data)
		} else {
			log.Printf("unmarshal error: %v", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
