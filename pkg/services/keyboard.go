package services

import (
	"fmt"

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
	kb := models.Keyboard{
		Inline:  true,
		OneTime: false,
	}
	kb.Buttons = make([][]models.Button, 2)
	kb.Buttons[0] = append(kb.Buttons[0], models.Button{Action: models.Action{Type: "text", Label: "Хочу купить схему", Payload: models.Payload{Button: "wantBuyPattern"}}, Color: "primary"})
	kb.Buttons[1] = append(kb.Buttons[1], models.Button{Action: models.Action{Type: "text", Label: "Хочу купить набор", Payload: models.Payload{Button: "wantBuySet"}}, Color: "primary"})

	err := ks.repos.SendMessage("Что Вас интересует?", &kb, req)
	if err != nil {
		return fmt.Errorf("error while sending product keyboard: %v", err)
	}
	return nil
}

func (ks *KeyboardService) ReplyToKeyboard(pl *models.Payload, mr *models.MessageRecipient) error {
	switch pl.Button {
	case "wantBuyPattern":
		if err := ks.sendWantPatternReply(mr); err != nil {
			return fmt.Errorf("error while handling wantBuyPattern key: %v", err)
		}
	case "wantBuySet":
		if err := ks.sendWantSetReply(mr); err != nil {
			return fmt.Errorf("error while handling wantBuySet key: %v", err)
		}
	}

	return nil
}

func (ks *KeyboardService) sendWantSetReply(mr *models.MessageRecipient) error {
	text := "Если Вам нужен набор или нарезка мулине для схемы, Вам надо обратиться к Татьяне Арефьевой: https://vk.com/id118768758\nОна занимается формированием наборов по моим схемам.\n\nЕще можно обратиться к Инне Ушаковой: https://vk.com/id40636067 Она тоже нарезкой занимается."
	err := ks.repos.SendMessage(text, nil, mr)
	if err != nil {
		return fmt.Errorf("err while sending 'wantBuySet' reply: %v", err)
	}
	return nil
}

func (ks *KeyboardService) sendWantPatternReply(mr *models.MessageRecipient) error {
	text := "Вы можете мои схемы купить прямо на сайте: https://forstitch.ru/\n\nКак купить(инструкция): https://forstitch.ru/how-to-buy/\n\nА так же на сайте: https://www.stitch.su/patterns?designer=8&buy=1\n\nЖителям Украины, которые хотят купить мои схемы, наборы по ним, можно обращаться к Миле Вождь https://vk.com/id15980005\nhttps://www.instagram.com/mika__mila_katya/\ne-mail: redbest.calico@gmail.com"
	err := ks.repos.SendMessage(text, nil, mr)
	if err != nil {
		return fmt.Errorf("err while sending 'wantBuySet' reply: %v", err)
	}
	return nil
}
