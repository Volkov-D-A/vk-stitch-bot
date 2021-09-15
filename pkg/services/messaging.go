package services

import (
	"fmt"
	"time"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/repository"
)

type MessagingService struct {
	repos  *repository.Repository
	config *config.Config
}

func NewMessagingService(repos *repository.Repository, config *config.Config) *MessagingService {
	return &MessagingService{
		repos:  repos,
		config: config,
	}
}

func (ms *MessagingService) InitDatabase() error {
	//Check filled database
	cnt, err := ms.repos.CountRecipients(nil)
	if err != nil {
		return fmt.Errorf("error while counts database %v", err)
	}
	//if database empty
	if cnt == 0 {
		res, err := ms.repos.GetGroupUsers()
		if err != nil {
			return fmt.Errorf("error while getting group users %v", err)
		}
		for i := 0; i < len(res)-1; i++ {
			result, err := ms.repos.CheckAllowedMessages(res[i])
			if err != nil {
				return fmt.Errorf("error while checking allowed messages: %v", err)
			}
			if result {
				cnt, err := ms.repos.CountRecipients(res[i])
				if err != nil {
					return fmt.Errorf("error while checking presence recipient: %v", err)
				}
				if cnt == 0 {
					if err = ms.repos.AddRecipient(&models.MessageRecipient{Id: res[i]}); err != nil {
						return fmt.Errorf("error while adding recipient: %v", err)
					}
				}
			} else {
				cnt, err := ms.repos.CountRecipients(res[i])
				if err != nil {
					return fmt.Errorf("error while checking presence recipient: %v", err)
				}
				if cnt == 1 {
					if err = ms.repos.DeleteRecipient(&models.MessageRecipient{Id: res[i]}); err != nil {
						return fmt.Errorf("error while deleting recipient: %v", err)
					}
				}
			}
			if i%20 == 0 {
				time.Sleep(time.Second * 1)
			}
		}
	}
	return nil
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

func (ms *MessagingService) SendMultipleMessages() error {
	list, err := ms.repos.GelAllRecipients()
	if err != nil {
		return err
	}
	fmt.Println(list)
	return nil
}
