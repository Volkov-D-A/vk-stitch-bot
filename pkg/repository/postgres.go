package repository

import (
	"context"
	"fmt"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/db/pg"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
)

type DataPostgres struct {
	db *pg.DB
}

func NewDataPostgres(db *pg.DB) *DataPostgres {
	return &DataPostgres{
		db: db,
	}
}

//AddRecipient adds messaging recipient into database
func (mp *DataPostgres) AddRecipient(rec *models.MessageRecipient) error {
	conn, err := mp.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	ct, err := conn.Exec(context.Background(), "INSERT INTO messages_recipients VALUES ($1)", rec.Id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("value not inserted")
	}
	return nil
}

//DeleteRecipient removes a messaging recipient from the database
func (mp *DataPostgres) DeleteRecipient(rec *models.MessageRecipient) error {
	conn, err := mp.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	ct, err := conn.Exec(context.Background(), "DELETE FROM messages_recipients WHERE recipient_vk_id=$1", rec.Id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("deletint item not found")
	}
	return nil
}

//GelAllRecipients returns list of all messaging recipients
func (mp *DataPostgres) GelAllRecipients() (*models.MessagingList, error) {
	result := models.MessagingList{}
	conn, err := mp.db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), "SELECT * FROM messages_recipients")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		result.List = append(result.List, id)
	}

	return &result, nil
}
