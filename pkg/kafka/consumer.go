package kafka

import (
	"context"
	"encoding/json"

	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/metrics"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

func (c *Consumer) ReadMessage(ctx context.Context, handler func([]byte) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			logger.GetLogger().WithError(err).Error("failed to read kafka message")
			return err
		}

		metrics.KafkaMessagesConsumed.Inc()
		logger.GetLogger().WithField("offset", msg.Offset).Info("kafka message received")

		if err := handler(msg.Value); err != nil {
			logger.GetLogger().WithError(err).Error("failed to handle kafka message")
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func UnmarshalMessage(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
