package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init("info")
	assert.NotNil(t, Log)
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	assert.NotNil(t, logger)
}

func TestInitWithInvalidLevel(t *testing.T) {
	Init("invalid")
	assert.NotNil(t, Log)
}

func TestGetLoggerWithoutInit(t *testing.T) {
	Log = nil
	logger := GetLogger()
	assert.NotNil(t, logger)
}
