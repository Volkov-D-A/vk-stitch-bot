package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"

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
	var srv []models.CallbackServerItem
	val := url.Values{}
	val.Add("group_id", cs.config.VKGroupID)
	result, err := cs.repos.SendRequest(&val, "groups.getCallbackServers", "items")
	if err != nil {
		return false, fmt.Errorf("error while sending 'getCallbackServers' request: %v", err)
	}
	js, _ := json.Marshal(result)
	if err = json.Unmarshal(js, &srv); err != nil {
		return false, fmt.Errorf("error while unmarshaling 'getCallbackServers' request: %v", err)
	}
	var validId []int
	for _, val := range srv {
		if val.ServerStatus != "ok" || val.ServerUrl != cs.config.CallbackUrl {
			if err := cs.RemoveCallbackServer(strconv.Itoa(val.ServerId)); err != nil {
				return false, fmt.Errorf("error while removing callback server: %v", err)
			}
		} else {
			validId = append(validId, val.ServerId)
		}
	}
	switch len(validId) {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		for i := 0; i < len(validId)-1; i++ {
			if err := cs.RemoveCallbackServer(strconv.Itoa(validId[i])); err != nil {
				return false, fmt.Errorf("error while removing callback server: %v", err)
			}
		}
		return true, nil
	}
}

func (cs *CallbackSetupService) RemoveCallbackServer(id string) error {
	val := url.Values{}
	val.Add("group_id", cs.config.VKGroupID)
	val.Add("server_id", id)
	_, err := cs.repos.SendRequest(&val, "groups.deleteCallbackServer", "")
	if err != nil {
		return fmt.Errorf("error while sending 'deleteCallbackServer' request: %v", err)
	}
	return nil
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
	if err = cs.SetupCallbackService(fmt.Sprintf("%.0f", result.(float64))); err != nil {
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
