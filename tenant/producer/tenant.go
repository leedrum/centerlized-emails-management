package producer

import (
	"github.com/IBM/sarama"
	"github.com/leedrum/centerlized-emails-management/tenant/models"
)

type TenantProducer struct {
	Topic        string
	SyncProducer sarama.SyncProducer
}

func (tp *TenantProducer) Publish(event models.TenantMessage) error {
	return publishToKafka(event, tp.SyncProducer, tp.Topic)
}

func publishToKafka(event models.TenantMessage, producer sarama.SyncProducer, topic string) error {
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
