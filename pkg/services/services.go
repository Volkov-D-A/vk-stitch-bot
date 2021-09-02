package services

import (
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/repository"
)

type Messaging interface {
	AddRecipient(rec *models.MessageRecipient) error
	DeleteRecipient(rec *models.MessageRecipient) error
	SendMessage() error
}

type Services struct {
	Messaging
}

func NewService(repos *repository.Repository) *Services {
	return &Services{
		Messaging: NewMessagingService(repos),
	}
}
