package domain

type NotificationChannel string

const (
	ChannelEmail    NotificationChannel = "email"
	ChannelSMS      NotificationChannel = "sms"
	ChannelTelegram NotificationChannel = "telegram"
)

type SendNotificationRequest struct {
	Channel   NotificationChannel `json:"channel"`
	Recipient string              `json:"recipient"`
	Subject   string              `json:"subject,omitempty"`
	Message   string              `json:"message"`
}

type SendNotificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
