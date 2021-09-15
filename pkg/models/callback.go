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
	MessageDate    int    `json:"date"`
	MessageFromId  int    `json:"from_id"`
	MessageId      int    `json:"id"`
	MessageText    string `json:"text"`
	MessagePayload string `json:"payload"`
}

//Client a part for TypeMessageNew struct
type Client struct {
	ClientKeyboard       bool `json:"keyboard"`
	ClientInlineKeyboard bool `json:"inline_keyboard"`
}

//MessageAllow struct for handle "message_allow" callback requests from VK api
type MessageAllow struct {
	UserId int `json:"user_id"`
}

//MessageDeny struct for handle "message_deny" callback requests from VK api
type MessageDeny struct {
	UserId int `json:"user_id"`
}

type CallbackServerItem struct {
	ServerId     int    `json:"id"`
	ServerUrl    string `json:"url"`
	ServerStatus string `json:"status"`
	ServerTitle  string `json:"title"`
}
