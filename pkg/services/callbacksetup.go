package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/repository"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"
)

type CallbackSetupService struct {
	repos  *repository.Repository
	config *config.Config
}

func NewCallbackSetupService(repos *repository.Repository, config *config.Config) *CallbackSetupService {
	return &CallbackSetupService{
		repos:  repos,
		config: config,
	}
}

func (cs *CallbackSetupService) CheckCallbackServerInfo() (bool, error) {
	srv, err := cs.repos.GetCallbackServerInfo()
	if err != nil {
		return false, fmt.Errorf("error getting callback server info: %v", err)
	}
	var validId []int
	for _, val := range srv {
		if val.ServerTitle == cs.config.Callback.Title {
			if val.ServerStatus != "ok" || val.ServerUrl != cs.config.Callback.URL {
				if err := cs.repos.RemoveCallbackServer(strconv.Itoa(val.ServerId)); err != nil {
					return false, fmt.Errorf("error while removing callback server: %v", err)
				}
			} else {
				validId = append(validId, val.ServerId)
			}
		}
	}
	switch len(validId) {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		for i := 0; i < len(validId)-1; i++ {
			if err := cs.repos.RemoveCallbackServer(strconv.Itoa(validId[i])); err != nil {
				return false, fmt.Errorf("error while removing callback server: %v", err)
			}
		}
		return true, nil
	}
}

func (cs *CallbackSetupService) SendConfirmationResponse(w http.ResponseWriter) error {
	code, err := cs.repos.GetConfirmationCode()
	if err != nil {
		return fmt.Errorf("error while getting confirmation code: %v", err)
	}
	//Send reply to VK api
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(code))
	if err != nil {
		return fmt.Errorf("error while sending response: %v", err)
	}
	return nil
}

func (cs *CallbackSetupService) SetupCallback() error {
	srvId, err := cs.repos.SetCallbackUrl()
	if err != nil {
		return fmt.Errorf("error while setting callback url: %v", err)
	}
	if err := cs.repos.SetupCallbackService(srvId); err != nil {
		return fmt.Errorf("error while setting callback service: %v", err)
	}
	return nil
}
