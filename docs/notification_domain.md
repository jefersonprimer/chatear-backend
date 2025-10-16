# Notification Domain

## Overview

The notification domain is responsible for sending notifications to users. Currently, it only supports sending emails, but it is designed to be extensible to support other notification channels like SMS or push notifications.

## Directory Structure

```
internal/notification/
├── application/
│   └── send.go
├── domain/
│   ├── notification.go
│   ├── repository.go
│   └── sender.go
├── infrastructure/
│   ├── postgres_repository.go
│   ├── smtp_sender.go
│   └── templates/
│       ├── email.txt
│       ├── welcome.html
│       └── magic_link.html
└── worker/
    └── email_consumer.go
```

## Domain

The domain layer contains the core business logic of the notification domain.

### Entities

*   `EmailSend`: Represents an email to be sent to a user with tracking information.
*   `Notification`: Legacy entity for general notifications (kept for backward compatibility).

### Interfaces

*   `Repository`: An interface for persisting email sends with methods for saving, retrieving by ID, and querying by recipient.
*   `Sender`: An interface for sending emails via various providers.

## Application

The application layer contains the use cases of the notification domain.

*   `EmailSender`: A use case for sending emails with error handling and status tracking.

### EmailSender Methods

*   `Send(ctx, recipient, subject, body, templateName)`: Sends an email and logs the result to the database.

## Infrastructure

The infrastructure layer contains the implementation of the domain interfaces.

*   `PostgresRepository`: A PostgreSQL implementation of the `Repository` interface that stores email sends in the `email_sends` table.
*   `SMTPSender`: An SMTP implementation of the `Sender` interface with template support.
*   `templates/`: A directory containing email templates in both plain text and HTML formats.

### Database Schema

The `email_sends` table stores all email sending attempts with the following fields:
- `id`: UUID primary key
- `recipient`: Email address of the recipient
- `subject`: Email subject
- `body`: Email body content
- `template_name`: Name of the template used (optional)
- `sent_at`: Timestamp when the email was sent
- `created_at`: Timestamp when the record was created
- `error_message`: Error message if sending failed (optional)
- `status`: Status of the email ('pending', 'sent', 'failed')

### SMTP Configuration

The SMTP sender uses the following environment variables:
- `SMTP_HOST`: SMTP server hostname
- `SMTP_PORT`: SMTP server port
- `SMTP_USER`: SMTP username
- `SMTP_PASS`: SMTP password
- `SMTP_FROM`: Sender email address

### Email Templates

The system supports multiple email templates:
- `email.txt`: Plain text template for general emails
- `welcome.html`: HTML template for welcome emails
- `magic_link.html`: HTML template for magic link authentication emails

Templates can be selected by specifying the `template_name` parameter when sending emails.

## Worker

The worker layer contains the NATS consumer that listens for `email.send` events and triggers the `EmailSender` use case.

### Event-Driven Architecture

The notification domain uses an event-driven architecture to send emails asynchronously. When another part of the system wants to send an email, it publishes an `email.send` event to the NATS server with the following JSON payload:

```json
{
  "recipient": "user@example.com",
  "subject": "Email Subject",
  "body": "Email body content",
  "template_name": "welcome"
}
```

The NATS consumer in the notification domain listens for these events and sends the email accordingly, logging the result to the database.

This decouples the email sending logic from the rest of the system, making it more resilient and scalable.

### Error Handling

The system includes comprehensive error handling:
- Failed email sends are logged with error messages
- Email status is tracked ('pending', 'sent', 'failed')
- All email attempts are stored in the database for auditing
- The system continues processing other emails even if one fails
