package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/leedrum/centerlized-emails-management/mail-fetcher/db"
	sqlc "github.com/leedrum/centerlized-emails-management/mail-fetcher/db/sqlc"

	"github.com/IBM/sarama"
)

func StartTokenUpdateConsumer(brokers []string, topic string, query *sqlc.Queries) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Failed to get partitions: %v", err)
	}

	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Failed to start partition consumer: %v", err)
		}
		defer pc.Close()

		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				var evt db.TokenUpdateEvent
				if err := json.Unmarshal(msg.Value, &evt); err != nil {
					log.Printf("Invalid token event: %v", err)
					continue
				}

				if err := db.UpsertToken(query, evt); err != nil {
					log.Printf("Failed to update token in DB: %v", err)
				} else {
					log.Printf("Updated token for %s (%s)", evt.UserEmail, evt.TenantName)
				}
			}
		}(pc)
	}

	<-context.Background().Done()
}
