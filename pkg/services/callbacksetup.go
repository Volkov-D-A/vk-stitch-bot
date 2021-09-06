package services

import (
	"fmt"
	"net/http"
	"net/url"

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

func (cs *CallbackSetupService) SendConfirmationResponse(w http.ResponseWriter) error {
	code, err := cs.GetConfirmationCode()
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

func (cs *CallbackSetupService) SetCallbackUrl() error {
	val := url.Values{}
	val.Add("url", cs.config.CallbackUrl)
	val.Add("title", "VKBot")
	val.Add("group_id", cs.config.VKGroupID)
	val.Add("secret_key", cs.config.VKCallbackSecret)
	result, err := cs.repos.SendRequest(&val, "groups.addCallbackServer", "server_id")
	if err != nil {
		return fmt.Errorf("error on sending 'addCallbackServer' request: %v", err)
	}
	if err = cs.SetupCallbackService(fmt.Sprintf("%f.0", result)); err != nil {
		return fmt.Errorf("error while setting up callback service: %v", err)
	}
	return nil
}

func (cs *CallbackSetupService) SetupCallbackService(srvId string) error {
	val := url.Values{}
	val.Add("group_id", cs.config.VKGroupID)
	val.Add("server_id", srvId)
	val.Add("api_version", "5.131")
	val.Add("message_allow", "1")
	val.Add("message_deny", "1")
	val.Add("message_new", "1")
	_, err := cs.repos.SendRequest(&val, "groups.setCallbackSettings", "")
	if err != nil {
		return fmt.Errorf("error on sending 'setCallbackSettings' request: %v", err)
	}
	return nil
}

func (cs *CallbackSetupService) GetConfirmationCode() (string, error) {
	val := url.Values{}
	val.Add("group_id", cs.config.VKGroupID)
	result, err := cs.repos.SendRequest(&val, "groups.getCallbackConfirmationCode", "code")
	if err != nil {
		return "", fmt.Errorf("error on sending 'getCallbackConfirmationCode' request: %v", err)
	}
	return result.(string), nil
}
