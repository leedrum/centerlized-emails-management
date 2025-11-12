package producer

import (
	"github.com/IBM/sarama"
	"github.com/leedrum/centerlized-emails-management/tenant/config"
	"github.com/leedrum/centerlized-emails-management/tenant/models"
)

type Producer interface {
	Publish(topic string, payload string)
}

func NewKafkaSyncProducer(cfg config.Config) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = cfg.KafkaRetryMax

	return sarama.NewSyncProducer(cfg.Brokers, config)
}

func PublishToKafka(event models.TenantMessage, producer sarama.SyncProducer, topic string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.AggregateID),
		Value: sarama.ByteEncoder(event.Payload),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(event.EventType)},
		},
	}

	_, _, err := producer.SendMessage(msg)
	return err
}
