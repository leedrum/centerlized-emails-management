package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) *Producer {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	return &Producer{producer: p}
}

func (p *Producer) Publish(topic string, msg interface{}) error {
	data, _ := json.Marshal(msg)
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
