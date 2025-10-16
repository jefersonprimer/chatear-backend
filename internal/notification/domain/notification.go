package domain

import "time"

type EmailSend struct {
	ID           string
	Recipient    string
	Subject      string
	Body         string
	TemplateName string
	SentAt       time.Time
	CreatedAt    time.Time
	ErrorMessage string
	Status       string
}

type Notification struct {
	ID        string
	Type      string
	Recipient string
	Subject   string
	Body      string
	SentAt    time.Time
	CreatedAt time.Time
}
