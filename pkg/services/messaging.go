package services

import (
	"fmt"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/logs"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/repository"
)

type MessagingService struct {
	repos  repository.Messaging
	logger *logs.Logger
	config *config.Config
}

func NewMessagingService(repos repository.Messaging, logger *logs.Logger, config *config.Config) *MessagingService {
	return &MessagingService{
		repos:  repos,
		logger: logger,
		config: config,
	}
}

func (ms *MessagingService) AddRecipient(rec *models.MessageRecipient) error {
	err := ms.repos.AddRecipient(rec)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MessagingService) DeleteRecipient(rec *models.MessageRecipient) error {
	err := ms.repos.DeleteRecipient(rec)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MessagingService) SendMessage() error {
	list, err := ms.repos.GelAllRecipients()
	if err != nil {
		return err
	}
	//TODO: wrote code to sending meessages
	fmt.Println(list)
	return nil
}
