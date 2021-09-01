package models

//Recipient receives message
type Recipient struct {
	Id int
}

//RecipientService methods for managing recipient
type RecipientService interface {
	Add(rec *Recipient) error
	Delete(rec *Recipient) error
	//Get() []*Recipient
}
