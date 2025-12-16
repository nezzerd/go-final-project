package service

import (
	"fmt"
	"hotel-booking-system/internal/delivery/domain"
	"hotel-booking-system/pkg/logger"

	tele "gopkg.in/telebot.v3"
)

type Notifier interface {
	SendNotification(req *domain.SendNotificationRequest) error
}

type DeliveryService struct {
	telegramBot *tele.Bot
}

func NewDeliveryService(telegramToken string) (*DeliveryService, error) {
	var bot *tele.Bot
	var err error

	if telegramToken != "" {
		bot, err = tele.NewBot(tele.Settings{
			Token: telegramToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create telegram bot: %w", err)
		}
	}

	return &DeliveryService{
		telegramBot: bot,
	}, nil
}

func (ds *DeliveryService) SendNotification(req *domain.SendNotificationRequest) error {
	switch req.Channel {
	case domain.ChannelEmail:
		return ds.sendEmail(req)
	case domain.ChannelSMS:
		return ds.sendSMS(req)
	case domain.ChannelTelegram:
		return ds.sendTelegram(req)
	default:
		return fmt.Errorf("unsupported channel: %s", req.Channel)
	}
}

func (ds *DeliveryService) sendEmail(req *domain.SendNotificationRequest) error {
	logger.GetLogger().WithFields(map[string]interface{}{
		"channel":   "email",
		"recipient": req.Recipient,
		"subject":   req.Subject,
	}).Info("sending email notification")

	return nil
}

func (ds *DeliveryService) sendSMS(req *domain.SendNotificationRequest) error {
	logger.GetLogger().WithFields(map[string]interface{}{
		"channel":   "sms",
		"recipient": req.Recipient,
	}).Info("sending SMS notification")

	return nil
}

func (ds *DeliveryService) sendTelegram(req *domain.SendNotificationRequest) error {
	if ds.telegramBot == nil {
		return fmt.Errorf("telegram bot not configured")
	}

	chatID, err := parseTelegramChatID(req.Recipient)
	if err != nil {
		return fmt.Errorf("invalid telegram chat ID: %w", err)
	}

	message := req.Message
	if req.Subject != "" {
		message = fmt.Sprintf("*%s*\n\n%s", req.Subject, req.Message)
	}

	_, err = ds.telegramBot.Send(&tele.User{ID: chatID}, message, &tele.SendOptions{
		ParseMode: tele.ModeMarkdown,
	})
	if err != nil {
		logger.GetLogger().WithError(err).Error("failed to send telegram notification")
		return err
	}

	logger.GetLogger().WithFields(map[string]interface{}{
		"channel":   "telegram",
		"recipient": req.Recipient,
	}).Info("telegram notification sent")

	return nil
}

func parseTelegramChatID(recipient string) (int64, error) {
	var chatID int64
	_, err := fmt.Sscanf(recipient, "%d", &chatID)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}
