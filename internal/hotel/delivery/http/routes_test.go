package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetupRoutes(t *testing.T) {
	mockUC := new(MockHotelUseCase)
	handler := NewHotelHandler(mockUC)

	mockUC.On("GetHotels", mock.Anything, mock.Anything, mock.Anything).Return([]interface{}{}, nil)

	r := SetupRoutes(handler)
	assert.NotNil(t, r)
}
