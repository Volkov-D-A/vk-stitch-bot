package repository

import (
	"net/url"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/db/pg"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
)

type Data interface {
	AddRecipient(rec *models.MessageRecipient) error
	DeleteRecipient(rec *models.MessageRecipient) error
	GelAllRecipients() (*models.MessagingList, error)
	CountRecipients(user interface{}) (int, error)
}

type Request interface {
	SendRequest(q *url.Values, method string, expectedResult string) (interface{}, error)
	SendMessage(text string, keyboard interface{}, mr *models.MessageRecipient) error
	GetCallbackServerInfo() ([]models.CallbackServerItem, error)
	RemoveCallbackServer(id string) error
	SetCallbackUrl() (string, error)
	GetConfirmationCode() (string, error)
	SetupCallbackService(srvId string) error
	GetGroupUsers() ([]int, error)
	CheckAllowedMessages(id int) (bool, error)
}

type Repository struct {
	Data
	Request
}

func NewRepository(db *pg.DB, conf *config.Config) *Repository {
	return &Repository{
		Data:    NewDataPostgres(db),
		Request: NewRequestApiRepository(conf),
	}
}
