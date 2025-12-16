package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRoutes(t *testing.T) {
	mockUC := new(MockBookingUseCase)
	handler := NewBookingHandler(mockUC)

	r := SetupRoutes(handler)
	assert.NotNil(t, r)
}
