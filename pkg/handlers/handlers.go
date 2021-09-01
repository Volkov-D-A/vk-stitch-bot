package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/logs"
)

var (
	logger          = logs.Get()
	errWrongGroupId = errors.New("VK group ID from VK api is invalid")
)

//CallbackHandler base struct for callback handler
type CallbackHandler struct {
	recRepo models.RecipientService
}

//NewCallbackHandler return a new callback handler
func NewCallbackHandler(recRepo models.RecipientService) *CallbackHandler {
	return &CallbackHandler{
		recRepo: recRepo,
	}
}

//Post base handler for callback requests from VK api
func (cb *CallbackHandler) Post(w http.ResponseWriter, r *http.Request) {

	//Check type of event
	req := models.CallbackRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Errorf("error while decoding JSON: %v", err)
	}
	// Switch action while type of event is
	switch req.EventType {
	case "confirmation":
		sendConfirmationResponse(&req, w)
	case "message_new":
		handleNewMessageEvent(&req)
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Unexpected event type"))
		if err != nil {
			logger.Errorf("error while sending response to VK api: %v", err)
		}
	}
}

//handleNewMessageEvent handle new_message callback request and if contains target string sending reply message and bot keyboard
func handleNewMessageEvent(req *models.CallbackRequest) {
	ms := models.TypeMessageNew{}
	err := json.Unmarshal(req.EventObject, &ms)
	if err != nil {
		logger.Error(err)
	}
	if strings.Contains(ms.Message.MessageText, "Меня заинтересовал данный товар.") && ms.MessageFromId == 50126581 {
		// Send BOT keyboard command
		logger.Info("Requested info about product")
	}
}

//sendConfirmationResponse uses to confirm callback url on VK settings
func sendConfirmationResponse(req *models.CallbackRequest, w http.ResponseWriter) {
	// FIXME: group id need to place into config
	if req.GroupId != 66770381 {
		logger.Errorf("error while confirmation callback server: %v", errWrongGroupId)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Bad requests"))
		if err != nil {
			logger.Errorf("error while sending response: %v", err)
		}
	}
	//Send reply to VK api
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("7a165ee7"))
	if err != nil {
		logger.Errorf("error while sending response: %v", err)
	}
}
