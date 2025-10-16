package events

type EmailSendRequest struct {
	Recipient    string `json:"recipient"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	TemplateName string `json:"template_name,omitempty"`
}
