package kafka

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProducer(t *testing.T) {
	producer := NewProducer([]string{"localhost:9092"}, "test-topic")
	assert.NotNil(t, producer)
	assert.NotNil(t, producer.writer)
}

func TestProducer_SendMessage_SerializationError(t *testing.T) {
	producer := NewProducer([]string{"localhost:9092"}, "test-topic")

	type InvalidType struct {
		Channel chan int
	}

	err := producer.SendMessage(context.Background(), "key", InvalidType{Channel: make(chan int)})
	assert.Error(t, err)
}
