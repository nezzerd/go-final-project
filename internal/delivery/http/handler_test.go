package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hotel-booking-system/internal/delivery/domain"
	"hotel-booking-system/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) SendNotification(req *domain.SendNotificationRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func TestDeliveryHandler_SendNotification(t *testing.T) {
	logger.Init("info")

	t.Run("success", func(t *testing.T) {
		mockNotifier := new(MockNotifier)
		mockNotifier.On("SendNotification", mock.AnythingOfType("*domain.SendNotificationRequest")).Return(nil)

		handler := NewDeliveryHandler(mockNotifier)

		reqBody := domain.SendNotificationRequest{
			Channel:   domain.ChannelEmail,
			Recipient: "test@example.com",
			Subject:   "Test Subject",
			Message:   "test message",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/notifications/send", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.SendNotification(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockNotifier.AssertExpectations(t)

		var response domain.SendNotificationResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response.Success)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mockNotifier := new(MockNotifier)
		handler := NewDeliveryHandler(mockNotifier)

		req := httptest.NewRequest("POST", "/api/notifications/send", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.SendNotification(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockNotifier.AssertNotCalled(t, "SendNotification")
	})

	t.Run("notifier error", func(t *testing.T) {
		mockNotifier := new(MockNotifier)
		mockNotifier.On("SendNotification", mock.AnythingOfType("*domain.SendNotificationRequest")).Return(assert.AnError)

		handler := NewDeliveryHandler(mockNotifier)

		reqBody := domain.SendNotificationRequest{
			Channel:   domain.ChannelEmail,
			Recipient: "test@example.com",
			Message:   "test message",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/notifications/send", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.SendNotification(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockNotifier.AssertExpectations(t)
	})
}

func TestNewDeliveryHandler(t *testing.T) {
	mockNotifier := new(MockNotifier)
	handler := NewDeliveryHandler(mockNotifier)
	assert.NotNil(t, handler)
	assert.Equal(t, mockNotifier, handler.notifier)
}
