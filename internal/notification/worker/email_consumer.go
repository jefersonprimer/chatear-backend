package worker

import (

	"context"

	"encoding/json"

	"log"



	"github.com/nats-io/nats.go"

	"github.com/jefersonprimer/chatear-backend/internal/notification/application"

	"github.com/jefersonprimer/chatear-backend/shared/events"

)



type NatsEmailConsumer struct {

	conn        *nats.Conn

	emailSender *application.EmailSender

}



func NewNatsEmailConsumer(conn *nats.Conn, emailSender *application.EmailSender) (*NatsEmailConsumer, error) {



	return &NatsEmailConsumer{



		conn:        conn,



		emailSender: emailSender,



	}, nil



}



func (c *NatsEmailConsumer) Start(ctx context.Context) {

	_, err := c.conn.Subscribe("email.send", func(msg *nats.Msg) {

		c.handleEmailSend(ctx, msg)

	})

	if err != nil {

		log.Fatalf("Error subscribing to email.send subject: %v", err)

	}



	log.Println("NATS email consumer started")

	<-ctx.Done()

}



func (c *NatsEmailConsumer) handleEmailSend(ctx context.Context, msg *nats.Msg) {

	var request events.EmailSendRequest

	if err := json.Unmarshal(msg.Data, &request); err != nil {

		log.Printf("Error unmarshaling email send request: %v", err)

		return

	}



	log.Printf("Processing email send request for recipient: %s", request.Recipient)



	emailSend, err := c.emailSender.Send(ctx, request.Recipient, request.Subject, request.Body, request.TemplateName)

	if err != nil {

		log.Printf("Error sending email to %s: %v", request.Recipient, err)

		return

	}



	log.Printf("Email sent successfully to %s with ID: %s", request.Recipient, emailSend.ID)

}
