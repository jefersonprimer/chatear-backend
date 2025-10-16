# Worker Architecture

This document outlines the architecture and design principles for background workers within the Chatear backend.

## Purpose

Workers are responsible for handling asynchronous tasks, processing events, and executing long-running operations that do not need to be part of the immediate request-response cycle of the main API. This approach helps in:

*   **Improving API Responsiveness**: Offloading heavy computations or I/O operations from the main API threads.
*   **Scalability**: Workers can be scaled independently based on the workload of specific tasks.
*   **Reliability**: Tasks can be retried in case of failures, ensuring eventual consistency.
*   **Decoupling**: Services communicate via events, reducing direct dependencies.

## Key Technologies

*   **NATS**: Used as the primary messaging system for event-driven communication between services and workers. Workers subscribe to specific subjects to consume events.
*   **Redis**: Utilized for caching, session management, rate limiting, and short-lived data storage. It provides fast access to frequently needed data.
*   **PostgreSQL**: The main relational database for persistent storage of application data.

## Common Worker Patterns

Most workers follow a similar pattern:

1.  **Event Consumption**: Workers subscribe to one or more NATS subjects. When an event is published to a subscribed subject, the worker receives and processes it.
2.  **Database Interaction**: Workers often interact with PostgreSQL to read or write persistent data related to the task.
3.  **Caching/Rate Limiting**: Redis is used to implement various caching strategies or to enforce rate limits on certain operations (e.g., number of emails sent per user, global deletion limits).
4.  **Idempotency**: Workers should be designed to handle duplicate events gracefully to ensure that processing an event multiple times does not lead to incorrect states.
5.  **Error Handling and Retries**: Robust error handling mechanisms are crucial. Failed tasks should ideally be retried, possibly with exponential backoff, and eventually moved to a dead-letter queue if persistent failures occur.

## Worker Examples

### Notification Worker (`cmd/worker/notification_worker.go`)

This worker is responsible for sending various types of notifications (e.g., emails, push notifications) based on events published to notification-related NATS subjects. It interacts with external services (like SMTP for emails) and may use Redis for rate limiting notification sends.

**Key Features:**
- Consumes `email.send` events from NATS
- Supports multiple email templates
- Integrates with SMTP for email delivery
- Logs all email sending activities
- Handles errors gracefully with proper logging

**Event Structure:**
```json
{
  "recipient": "user@example.com",
  "subject": "Email Subject",
  "body": "Email Body Content",
  "template_name": "welcome" // optional
}
```

### User Deletion Worker (`cmd/worker/user_delete_worker.go`)

This worker handles the asynchronous process of user account deletion with a 24-hour recovery period. It consumes `user.delete` events, manages the `user_deletions` table in PostgreSQL, and uses Redis to enforce rate limits on global deletions and per-user recovery email sends.

**Key Features:**
- **24-Hour Recovery Period**: Users have 24 hours to recover their account before permanent deletion
- **Rate Limiting**: 
  - Maximum 10 deletions per day globally
  - Maximum 2 recovery emails per user per day
- **Recovery Email System**: Sends recovery emails 24 hours before scheduled deletion
- **Soft Delete Implementation**: Marks users as deleted without removing data immediately
- **Audit Logging**: Logs all deletion activities for compliance

**Event Structure:**
```json
{
  "user_id": "uuid-of-user-to-delete"
}
```

**Deletion Process:**
1. **Event Received**: Worker receives `user.delete` event
2. **Schedule Deletion**: Inserts record into `user_deletions` table with 24-hour delay
3. **Recovery Email**: Sends recovery email 24 hours before deletion
4. **Execution**: Performs soft delete after 24-hour period
5. **Rate Limiting**: Enforces daily limits using Redis counters

**Database Tables Used:**
- `user_deletions`: Tracks scheduled deletions and their status
- `users`: Main user table (soft delete via `is_deleted` flag)
- `action_logs`: Audit trail for deletion activities

**Redis Keys Used:**
- `global:deletion:count:YYYY-MM-DD`: Global deletion counter
- `user:email:count:USERID:YYYY-MM-DD`: Per-user email counter

## Adding a New Worker

To add a new worker:

1.  **Create a new Go file** in the `cmd/worker/` directory (e.g., `cmd/worker/new_worker.go`).
2.  **Define the `main` function** to initialize connections to NATS, PostgreSQL, and Redis as needed.
3.  **Subscribe to relevant NATS subjects** to consume events.
4.  **Implement event handlers** to process incoming messages, including business logic, database interactions, and any necessary caching or rate limiting.
5.  **Consider idempotency and error handling** for robust processing.
6.  **Update the `Makefile` or CI/CD pipeline** to build and deploy the new worker.

## Running Workers

### Development Mode

To run workers in development mode:

```bash
# Run notification worker
go run cmd/worker/notification_worker.go

# Run user deletion worker
go run cmd/worker/user_delete_worker.go
```

### Production Mode

Build and run workers as standalone executables:

```bash
# Build workers
go build -o notification-worker cmd/worker/notification_worker.go
go build -o user-delete-worker cmd/worker/user_delete_worker.go

# Run workers
./notification-worker
./user-delete-worker
```

### Docker Compose

Workers can be run using the provided `docker-compose.events.yml`:

```bash
# Start all services including workers
docker-compose -f docker-compose.events.yml up

# Run only workers
docker-compose -f docker-compose.events.yml up notification-worker user-delete-worker
```

## Environment Variables

All workers require the following environment variables:

- `DATABASE_URL`: PostgreSQL connection string
- `NATS_URL`: NATS server connection string
- `REDIS_URL`: Redis server connection string
- `SMTP_HOST`: SMTP server hostname
- `SMTP_PORT`: SMTP server port
- `SMTP_USER`: SMTP username
- `SMTP_PASS`: SMTP password
- `SMTP_FROM`: Sender email address

## Monitoring and Logging

Workers provide comprehensive logging for:
- Connection status to external services
- Event processing activities
- Error conditions and recovery
- Rate limiting enforcement
- Database operations

## Deployment

Each worker is typically deployed as a separate process or container, allowing for independent scaling and management. Configuration (e.g., NATS URL, Database URL, Redis URL) should be managed via environment variables.

### Scaling Considerations

- **Notification Worker**: Scale based on email volume and SMTP server capacity
- **User Deletion Worker**: Typically runs as a single instance due to rate limiting constraints
- **Database Connections**: Monitor connection pool usage when scaling workers
- **Redis Usage**: Ensure Redis has sufficient memory for rate limiting counters
