package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTracer(t *testing.T) {
	tp, err := InitTracer("test-service", "http://localhost:14268/api/traces")
	assert.NoError(t, err)
	assert.NotNil(t, tp)

	err = Shutdown(context.Background(), tp)
	assert.NoError(t, err)
}

func TestShutdown_NilProvider(t *testing.T) {
	err := Shutdown(context.Background(), nil)
	assert.NoError(t, err)
}
