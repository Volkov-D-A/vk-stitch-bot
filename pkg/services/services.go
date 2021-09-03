package services

import (
	"net/http"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/logs"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/repository"
)

//Messaging uses for sending messages to VK group members
type Messaging interface {
	AddRecipient(req *models.MessageRecipient) error
	DeleteRecipient(req *models.MessageRecipient) error
	SendMessage() error
}

//CallbackSetup uses for setup and confirm parameters VK callback server
type CallbackSetup interface {
	SendConfirmationResponse(w http.ResponseWriter) error
	SetCallbackUrl() error
	GetConfirmationCode() (string, error)
	SetupCallbackService(srvId string) error
}

type Services struct {
	Messaging
	CallbackSetup
}

func NewService(repos *repository.Repository, logger *logs.Logger, config *config.Config) *Services {
	return &Services{
		Messaging:     NewMessagingService(repos, logger, config),
		CallbackSetup: NewCallbackSetupService(config),
	}
}
