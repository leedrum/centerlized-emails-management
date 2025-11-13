package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	topic         string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	p := &Producer{
		asyncProducer: producer,
		topic:         topic,
	}

	// handle async errors/successes
	go func() {
		for {
			select {
			case err := <-p.asyncProducer.Errors():
				log.Printf("[Kafka] Failed to send message: %v", err)
			case msg := <-p.asyncProducer.Successes():
				log.Printf("[Kafka] Sent message to %s [partition:%d offset:%d]", msg.Topic, msg.Partition, msg.Offset)
			}
		}
	}()

	return p, nil
}

func (p *Producer) Publish(key string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(payload),
	}

	p.asyncProducer.Input() <- message
	return nil
}

func (p *Producer) Close() error {
	return p.asyncProducer.Close()
}
