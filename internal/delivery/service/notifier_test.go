package service

import (
	"testing"

	"hotel-booking-system/internal/delivery/domain"
	"hotel-booking-system/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDeliveryService(t *testing.T) {
	logger.Init("info")

	t.Run("success with telegram token", func(t *testing.T) {
		service, err := NewDeliveryService("test-token")
		assert.Error(t, err)
		assert.Nil(t, service)
	})

	t.Run("success without telegram token", func(t *testing.T) {
		service, err := NewDeliveryService("")
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})
}

func TestDeliveryService_SendNotification(t *testing.T) {
	logger.Init("info")

	t.Run("unsupported channel", func(t *testing.T) {
		service, err := NewDeliveryService("")
		require.NoError(t, err)

		req := &domain.SendNotificationRequest{
			Channel:   "unsupported",
			Recipient: "test@example.com",
			Message:   "test message",
		}

		err = service.SendNotification(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported channel")
	})

	t.Run("email channel", func(t *testing.T) {
		service, err := NewDeliveryService("")
		require.NoError(t, err)

		req := &domain.SendNotificationRequest{
			Channel:   domain.ChannelEmail,
			Recipient: "test@example.com",
			Subject:   "Test Subject",
			Message:   "test message",
		}

		err = service.SendNotification(req)
		assert.NoError(t, err)
	})

	t.Run("sms channel", func(t *testing.T) {
		service, err := NewDeliveryService("")
		require.NoError(t, err)

		req := &domain.SendNotificationRequest{
			Channel:   domain.ChannelSMS,
			Recipient: "+1234567890",
			Message:   "test message",
		}

		err = service.SendNotification(req)
		assert.NoError(t, err)
	})

	t.Run("telegram channel without bot", func(t *testing.T) {
		service, err := NewDeliveryService("")
		require.NoError(t, err)

		req := &domain.SendNotificationRequest{
			Channel:   domain.ChannelTelegram,
			Recipient: "123456789",
			Message:   "test message",
		}

		err = service.SendNotification(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "telegram bot not configured")
	})
}

func TestParseTelegramChatID(t *testing.T) {
	tests := []struct {
		name      string
		recipient string
		want      int64
		wantErr   bool
	}{
		{
			name:      "valid chat ID",
			recipient: "123456789",
			want:      123456789,
			wantErr:   false,
		},
		{
			name:      "invalid chat ID",
			recipient: "invalid",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "empty string",
			recipient: "",
			want:      0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTelegramChatID(tt.recipient)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
