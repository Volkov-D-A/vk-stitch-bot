package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/services"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/logs"
)

var (
	errSecretMismatch = errors.New("secrets from api and local do not match")
	errWrongGroupId   = errors.New("VK group ID from VK api is invalid")
)

//CallbackHandler base struct for callback handler
type CallbackHandler struct {
	services *services.Services
	logger   *logs.Logger
	config   *config.Config
}

//NewCallbackHandler return a new callback handler
func NewCallbackHandler(services *services.Services, logger *logs.Logger, config *config.Config) *CallbackHandler {
	return &CallbackHandler{
		services: services,
		logger:   logger,
		config:   config,
	}
}

//Post base handler for callback requests from VK api
func (cb *CallbackHandler) Post(w http.ResponseWriter, r *http.Request) {

	//Check type of event
	req := models.CallbackRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		cb.logger.Errorf("error while decoding JSON: %v", err)
		return
	}
	//Check group id in callback request from api
	if req.GroupId.String() != cb.config.VKGroupID {
		cb.logger.Error(errWrongGroupId)
		return
	}
	//Check secret in callback request from api
	if req.Secret != cb.config.VKCallbackSecret {
		cb.logger.Error(errSecretMismatch)
		return
	}
	// Switch action while type of event is
	switch req.EventType {
	case "confirmation":
		if err = cb.handleConfirmationEvent(w); err != nil {
			cb.logger.Errorf("error while handling confirmation event: %v", err)
		}
		return
	case "message_deny":
		if err = cb.handleMessageDenyEvent(&req); err != nil {
			cb.logger.Errorf("error while handling message_deny event: %v", err)
		}
	case "message_allow":
		if err = cb.handleMessageAllowEvent(&req); err != nil {
			cb.logger.Errorf("error while handling message_allow event: %v", err)
		}
	case "message_new":
		if err = cb.handleMessageNewEvent(&req); err != nil {
			cb.logger.Errorf("error while handling message_new event: %v", err)
		}
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		cb.logger.Errorf("error while sending response to VK api: %v", err)
	}
}

//InitRoutes initializing routes for callback handler
func (cb *CallbackHandler) InitRoutes() *http.ServeMux {
	//Route handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", cb.Post)
	return mux
}

//handleConfirmationEvent handle confirmation event from callback api
func (cb *CallbackHandler) handleConfirmationEvent(w http.ResponseWriter) error {
	err := cb.services.SendConfirmationResponse(w)
	if err != nil {
		return fmt.Errorf("error while sending confirmation response: %v", err)
	}
	return nil
}

//handleMessageAllowEvent handle message_allow event from callback api
func (cb *CallbackHandler) handleMessageAllowEvent(req *models.CallbackRequest) error {
	ma := models.MessageAllow{}
	if err := json.Unmarshal(req.EventObject, &ma); err != nil {
		return fmt.Errorf("error while unmarshaling 'message_allow' event object: %v", err)
	}
	if err := cb.services.Messaging.AddRecipient(&models.MessageRecipient{Id: ma.UserId}); err != nil {
		return fmt.Errorf("error while adding recipient %v", err)
	}
	return nil
}

//handleMessageDenyEvent handle message_deny event from callback api
func (cb *CallbackHandler) handleMessageDenyEvent(req *models.CallbackRequest) error {
	md := models.MessageDeny{}
	if err := json.Unmarshal(req.EventObject, &md); err != nil {
		return fmt.Errorf("error while unmarshaling 'message_deny' event object: %v", err)
	}
	if err := cb.services.Messaging.DeleteRecipient(&models.MessageRecipient{Id: md.UserId}); err != nil {
		return fmt.Errorf("error while deleting recipient %v", err)
	}
	return nil
}

//handleNewMessageEvent handle new_message callback request and if contains target string sending reply message and bot keyboard
func (cb *CallbackHandler) handleMessageNewEvent(req *models.CallbackRequest) error {
	ms := models.TypeMessageNew{}
	err := json.Unmarshal(req.EventObject, &ms)
	if err != nil {
		return fmt.Errorf("error while unmarshaling 'message_new' event object: %v", err)
	}
	//Case 1: Interesting product -> sending bot keyboard
	if strings.Contains(ms.Message.MessageText, "Меня заинтересовал этот товар.") || strings.Contains(ms.Message.MessageText, "Меня заинтересовал данный товар.") && ms.MessageFromId == 50126581 {
		if err = cb.services.Keyboard.SendProductKeyboard(&models.MessageRecipient{Id: ms.MessageFromId}); err != nil {
			return fmt.Errorf("error sending product keyboard: %v", err)
		}
	}
	//Case 2: Service message -> do mailing message to recipients
	if strings.Contains(ms.Message.MessageText, "@@@!@@@") && ms.MessageFromId == 22478488 {
		// Mailing
	}
	return nil
}
