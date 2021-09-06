package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/repository"

	"github.com/Volkov-D-A/vk-stitch-bot/pkg/config"
	"github.com/Volkov-D-A/vk-stitch-bot/pkg/models"
)

type KeyboardService struct {
	repos  *repository.Repository
	config *config.Config
}

func NewKeyboardService(repos *repository.Repository, config *config.Config) *KeyboardService {
	return &KeyboardService{
		repos:  repos,
		config: config,
	}
}

func (ks *KeyboardService) SendProductKeyboard(req *models.MessageRecipient) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	val := url.Values{}
	kb := models.Keyboard{
		Inline:  true,
		OneTime: false,
	}
	kb.Buttons = make([][]models.Button, 1)
	kb.Buttons[0] = append(kb.Buttons[0], models.Button{Action: models.Action{Type: "text", Label: "Хочу купить схему", Payload: []string{"wantBuyPattern"}}})
	kb.Buttons[0] = append(kb.Buttons[0], models.Button{Action: models.Action{Type: "text", Label: "Хочу купить набор", Payload: []string{"wantBuySet"}}})
	js, err := json.Marshal(kb)
	if err != nil {
		return fmt.Errorf("error marshalling keyboard: %v", err)
	}
	val.Add("keyboard", string(js))
	val.Add("peer_id", strconv.Itoa(req.Id))
	val.Add("message", "Что Вас интересует?")
	val.Add("random_id", strconv.Itoa(r.Intn(1000000000)))
	_, err = ks.repos.SendRequest(&val, "messages.send", "")
	if err != nil {
		return fmt.Errorf("error while sending product keyboard request: %v", err)
	}
	return nil
}
