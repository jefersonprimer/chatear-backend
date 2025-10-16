package domain

import "context"

type Sender interface {
	Send(ctx context.Context, emailSend *EmailSend) error
}
