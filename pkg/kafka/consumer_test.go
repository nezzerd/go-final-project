package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsumer(t *testing.T) {
	consumer := NewConsumer([]string{"localhost:9092"}, "test-topic", "test-group")
	assert.NotNil(t, consumer)
	assert.NotNil(t, consumer.reader)
}

func TestUnmarshalMessage(t *testing.T) {
	data := []byte(`{"test": "value"}`)
	var result map[string]string
	err := UnmarshalMessage(data, &result)
	assert.NoError(t, err)
	assert.Equal(t, "value", result["test"])
}

func TestUnmarshalMessage_InvalidJSON(t *testing.T) {
	data := []byte(`invalid json`)
	var result map[string]string
	err := UnmarshalMessage(data, &result)
	assert.Error(t, err)
}
