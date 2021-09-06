package models

type MessageRecipient struct {
	Id int
}

type MessagingList struct {
	List []int
}

type VkMessage struct {
	Keyboard    string
	PeerId      string
	MessageText string
}
