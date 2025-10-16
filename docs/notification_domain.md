# Notification Domain Documentation

The Notification Domain is responsible for handling all outgoing communications, primarily email notifications. It ensures reliable delivery and provides a mechanism for tracking sent emails.

## Key Components

*   **Entities (`domain/entities`)**:
    *   `EmailSend`: Represents a record of an email that has been sent or is queued to be sent, including details like recipient, subject, body, and status.

*   **Repositories (`domain/repositories`, `internal/notification/infrastructure`)**:
    *   `EmailSendRepository`: Interface for persisting and retrieving `EmailSend` records.
    *   `postgres_repository.go`: Concrete implementation of `EmailSendRepository` using PostgreSQL.

*   **Application Services (`internal/notification/application`)**:
    *   `EmailService`: Orchestrates the email sending process. It prepares email content, records the `EmailSend` entity, and dispatches the email sending task (e.g., to a NATS queue).
    *   `send.go`: Contains the core logic for preparing and sending different types of emails (e.g., welcome, magic link).

*   **Domain Interfaces (`internal/notification/domain`)**:
    *   `Notification`: Defines the structure of a generic notification.
    *   `Repository`: Generic repository interface (implemented by `EmailSendRepository`).
    *   `Sender`: Interface for sending emails (e.g., `SMTPSender`).

*   **Infrastructure (`internal/notification/infrastructure`)**:
    *   `smtp_sender.go`: Concrete implementation of the `Sender` interface using SMTP to dispatch emails.
    *   `templates/`: Directory containing email templates (e.g., `magic_link.html`, `welcome.html`, `email.txt`).

*   **Worker (`internal/notification/worker`)**:
    *   `email_consumer.go`: A background worker that listens to a NATS queue for email sending requests. Upon receiving a request, it retrieves the email details, uses the `SMTPSender` to send the email, and updates the `EmailSend` record status.

## Notification Workflow Summary

```mermaid
graph TD
    A[Application Layer (e.g., User Registration)] --> B[EmailService.Send(EmailType, Data)]
    B --> C[EmailService: Create EmailSend Record]
    C --> D[EmailSendRepository.Save]
    D --> E[PostgreSQL]
    C --> F[Publish EmailSend Event to NATS]
    F --> G[NATS Messaging System]
    G --> H[EmailConsumer Worker]
    H --> I[EmailConsumer: Retrieve EmailSend Record]
    I --> J[SMTPSender.SendEmail]
    J --> K[External SMTP Server]
    K --> L[Recipient's Inbox]
    J --> M[EmailSendRepository.UpdateStatus (Sent/Failed)]
    M --> E
```

**Bullet Summary of Core Notification Workflow:**

*   **Initiation**: An action in the application (e.g., user registration, password reset) triggers the `EmailService` to send a specific type of email.
*   **Email Preparation & Recording**: The `EmailService` prepares the email content using templates, creates an `EmailSend` entity with a `PENDING` status, and saves it to the PostgreSQL database via `EmailSendRepository`.
*   **Event Publishing**: The `EmailService` then publishes an event (containing the `EmailSend` ID or full details) to a NATS queue, signaling that an email needs to be sent.
*   **Worker Consumption**: The `EmailConsumer` worker, subscribed to the NATS queue, receives the event.
*   **Email Sending**: The `EmailConsumer` retrieves the `EmailSend` record, uses the `SMTPSender` to dispatch the email through an external SMTP server.
*   **Status Update**: After attempting to send, the `EmailConsumer` updates the `EmailSend` record in PostgreSQL to `SENT` or `FAILED`, along with any relevant details (e.g., error messages).