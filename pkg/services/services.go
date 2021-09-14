package services

import (
	"net/http"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/repository"
)

//Messaging uses for sending messages to VK group members
type Messaging interface {
	AddRecipient(req *models.MessageRecipient) error
	DeleteRecipient(req *models.MessageRecipient) error
	SendMultipleMessages() error
	InitDatabase() error
}

//CallbackSetup uses for setup and confirm parameters VK callback server
type CallbackSetup interface {
	SendConfirmationResponse(w http.ResponseWriter) error
	SetupCallback() error
	CheckCallbackServerInfo() (bool, error)
}

type Keyboard interface {
	SendProductKeyboard(req *models.MessageRecipient) error
	ReplyToKeyboard(pl *models.Payload, mr *models.MessageRecipient) error
}

type Services struct {
	Messaging
	CallbackSetup
	Keyboard
}

func NewService(repos *repository.Repository, config *config.Config) *Services {
	return &Services{
		Messaging:     NewMessagingService(repos, config),
		CallbackSetup: NewCallbackSetupService(repos, config),
		Keyboard:      NewKeyboardService(repos, config),
	}
}
