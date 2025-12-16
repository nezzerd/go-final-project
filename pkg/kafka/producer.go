package kafka

import (
	"context"
	"encoding/json"

	"hotel-booking-system/pkg/logger"
	"hotel-booking-system/pkg/metrics"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) SendMessage(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		logger.GetLogger().WithError(err).Error("failed to send kafka message")
		return err
	}

	metrics.KafkaMessagesProduced.Inc()
	logger.GetLogger().WithField("key", key).Info("kafka message sent")
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
