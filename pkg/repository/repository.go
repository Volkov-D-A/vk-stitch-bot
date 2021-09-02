package repository

import (
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/db/pg"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
)

type Messaging interface {
	AddRecipient(rec *models.MessageRecipient) error
	DeleteRecipient(rec *models.MessageRecipient) error
	GelAllRecipients() (*models.MessagingList, error)
}

type Repository struct {
	Messaging
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		Messaging: NewMessagingPostgres(db),
	}
}
