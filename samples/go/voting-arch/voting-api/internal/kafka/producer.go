package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers string) (*Producer, error) {
	dialer := &kafka.Dialer{
		Timeout:   10 * 1000000000,
		DualStack: true,
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokers},
		Dialer:  dialer,
	})

	return &Producer{writer: writer}, nil
}

func (p *Producer) SendMessage(topic string, message []byte) error {
	err := p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Value: message,
		},
	)
	return err
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
