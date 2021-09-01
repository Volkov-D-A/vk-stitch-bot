package repository

import (
	"fmt"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/db/pgdb"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
)

//RecipientRepository base struct for the repository
type RecipientRepository struct {
	db *pgdb.DB
}

//New creates a new RecipientRepository object
func New(db *pgdb.DB) *RecipientRepository {
	return &RecipientRepository{
		db: db,
	}
}

func (repo *RecipientRepository) Add(rec *models.Recipient) error {
	_, err := repo.db.Model(rec).Insert()
	if err != nil {
		return fmt.Errorf("SQL insert failed: %v", err)
	}
	return nil
}

func (repo *RecipientRepository) Delete(rec *models.Recipient) error {
	_, err := repo.db.Model(rec).Delete()
	if err != nil {
		return fmt.Errorf("SQL deletion failed: %v", err)
	}
	return nil
}
