package kafka

import (
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	producer *kafka.Producer
	logger   *logrus.Logger
	mu       sync.Mutex
}

func NewProducer(brokers string) (*Producer, error) {
	config := kafka.ConfigMap{
		// Broker connection
		"bootstrap.servers": brokers,

		// Zero data loss guarantees
		"acks":                                  "all", // Wait for all replicas to acknowledge
		"retries":                               10,    // Retry failed sends
		"max.in.flight.requests.per.connection": 5,     // Maintain ordering

		// Performance optimizations
		"linger.ms":        100,      // Wait up to 100ms to batch messages
		"batch.size":       16384,    // Batch size in bytes (16KB)
		"compression.type": "snappy", // Compress messages

		// Timeout settings
		"socket.timeout.ms":       30000,  // Socket timeout
		"request.timeout.ms":      30000,  // Request timeout
		"delivery.timeout.ms":     120000, // 2 minute total timeout
		"socket.keepalive.enable": true,   // Keep connections alive

		// Idempotence for exactly-once semantics
		"enable.idempotence": true, // Prevent duplicates
	}

	producer, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		logger:   logrus.New(),
		mu:       sync.Mutex{},
	}, nil
}

func (p *Producer) SendMessage(topic string, message []byte) error {
	deliveryChan := make(chan kafka.Event, 1)

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}, deliveryChan)

	if err != nil {
		return err
	}

	// Wait for delivery confirmation with timeout
	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return m.TopicPartition.Error
		}
		return nil

	case <-time.After(30 * time.Second):
		return kafka.NewError(kafka.ErrMsgTimedOut, "delivery confirmation timeout", false)
	}
}

func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Flush all pending messages with 30 second timeout
	p.producer.Flush(30 * 1000)
	p.producer.Close()
	return nil
}
