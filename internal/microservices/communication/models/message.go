package models

type SMSMessage struct {
	PhoneNumber string `json:"phone"`
	Body        string `json:"message"`
	Originator  string `json:"originator"`
}

type EmailMessage struct {
	EmailAddress string
	Body         string
	Originator   string
}
