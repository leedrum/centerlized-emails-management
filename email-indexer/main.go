package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/lib/pq"

	dbpkg "github.com/leedrum/centerlized-emails-management/email-indexer/internal/db"
	"github.com/leedrum/centerlized-emails-management/email-indexer/internal/indexer"
	"github.com/leedrum/centerlized-emails-management/email-indexer/kafka"
)

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("env %s required", k)
	}
	return v
}

func main() {
	// --- Load environment ---
	pgDSN := mustEnv("DATABASE_URL")
	kafkaBrokers := mustEnv("KAFKA_BROKERS") // comma-separated list
	consumeTopic := getEnv("KAFKA_TOPIC", "email_indexer.reindex")
	dlqTopic := getEnv("KAFKA_DLQ_TOPIC", "email_indexer.dead_letter")
	consumerGroupID := getEnv("KAFKA_CONSUMER_GROUP", "email-indexer-group")

	esAddr := mustEnv("ELASTICSEARCH_ADDR")

	// --- Connect PostgreSQL ---
	sqlDB, err := sql.Open("postgres", pgDSN)
	if err != nil {
		log.Fatalf("pg connect: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("pg ping: %v", err)
	}
	log.Println("[OK] Connected PostgreSQL")

	q := dbpkg.New(sqlDB)

	// --- Connect Elasticsearch ---
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esAddr},
	})
	if err != nil {
		log.Fatalf("create es client: %v", err)
	}
	if res, err := esClient.Info(); err != nil {
		log.Fatalf("es info: %v", err)
	} else {
		res.Body.Close()
	}
	log.Println("[OK] Connected Elasticsearch")

	// --- Create indexing service ---
	svc := indexer.NewService(sqlDB, q, esClient)

	// --- Kafka config ---
	brokers := strings.Split(kafkaBrokers, ",")

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	cfg.Consumer.Return.Errors = true

	// --- DLQ Producer ---
	producerCfg := sarama.NewConfig()
	producerCfg.Producer.Return.Successes = true
	producerCfg.Version = sarama.V2_8_0_0

	dlqProducer, err := sarama.NewSyncProducer(brokers, producerCfg)
	if err != nil {
		log.Fatalf("create DLQ producer: %v", err)
	}
	defer dlqProducer.Close()

	// --- Create consumer group ---
	consumerGroup, err := sarama.NewConsumerGroup(brokers, consumerGroupID, cfg)
	if err != nil {
		log.Fatalf("create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Pass DLQ producer + DLQ topic to Kafka handler
	consumer := &kafka.KafkaConsumer{
		Svc:        svc,
		DLQ:        dlqProducer,
		DLQTopic:   dlqTopic,
		MaxRetries: 3,
	}

	ctx, cancel := context.WithCancel(context.Background())

	// --- Run Kafka consumer loop ---
	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{consumeTopic}, consumer); err != nil {
				log.Printf("Kafka consume error: %v", err)
				time.Sleep(time.Second)
			}
		}
	}()

	log.Printf("[READY] Listening Kafka topic: %s", consumeTopic)

	// --- Wait for shutdown ---
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("shutting down...")

	cancel()
	time.Sleep(time.Second)
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
