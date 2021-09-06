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
}

type Request interface {
	SendRequest(q *url.Values, method string, expectedResult string) (interface{}, error)
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
