package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPRequestsTotal(t *testing.T) {
	assert.NotNil(t, HTTPRequestsTotal)
	HTTPRequestsTotal.WithLabelValues("GET", "/test", "200").Inc()
}

func TestHTTPRequestDuration(t *testing.T) {
	assert.NotNil(t, HTTPRequestDuration)
	HTTPRequestDuration.WithLabelValues("GET", "/test").Observe(0.5)
}

func TestGRPCRequestsTotal(t *testing.T) {
	assert.NotNil(t, GRPCRequestsTotal)
	GRPCRequestsTotal.WithLabelValues("GetHotel", "success").Inc()
}

func TestKafkaMessagesProduced(t *testing.T) {
	assert.NotNil(t, KafkaMessagesProduced)
	KafkaMessagesProduced.Inc()
}

func TestKafkaMessagesConsumed(t *testing.T) {
	assert.NotNil(t, KafkaMessagesConsumed)
	KafkaMessagesConsumed.Inc()
}
