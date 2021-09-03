package models

import "encoding/json"

//CallbackRequest base struct for callback requests from VK api
type CallbackRequest struct {
	EventType   string          `json:"type"`
	EventObject json.RawMessage `json:"object"`
	GroupId     json.Number     `json:"group_id"`
	EventId     string          `json:"event_id"`
	Secret      string          `json:"secret"`
}

//TypeMessageNew struct for handle "message_new" callback requests from VK api
type TypeMessageNew struct {
	Message `json:"message"`
	Client  `json:"client_info"`
}

//Message a part for TypeMessageNew struct
type Message struct {
	MessageDate   int    `json:"date"`
	MessageFromId int    `json:"from_id"`
	MessageId     int    `json:"id"`
	MessageText   string `json:"text"`
}

//Client a part for TypeMessageNew struct
type Client struct {
	ClientKeyboard       bool `json:"keyboard"`
	ClientInlineKeyboard bool `json:"inline_keyboard"`
}
