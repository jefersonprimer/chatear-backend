# Environment Variables

This document details the environment variables used to configure the application. These variables are typically loaded from a `.env` file during local development and set in the deployment environment for production.

## Application Configuration

*   **`APP_URL`**
    *   **Description**: The base URL of the application. Used for constructing absolute URLs in emails (e.g., magic links) or API responses.
    *   **Example**: `http://localhost:8080`

*   **`PORT`**
    *   **Description**: The port on which the API server will listen for incoming HTTP requests.
    *   **Example**: `8080`

*   **`DATABASE_URL`**
    *   **Description**: The connection string for the PostgreSQL database. This includes credentials, host, port, and database name.
    *   **Example**: `postgresql://user:password@host:port/database_name`

## Supabase Configuration (Optional)

*   **`SUPABASE_URL`**
    *   **Description**: The URL for your Supabase project. Used if integrating with Supabase services.
    *   **Example**: `https://your-project-id.supabase.co`

*   **`SUPABASE_ANON_KEY`**
    *   **Description**: The anonymous key for your Supabase project. Used for client-side interactions with Supabase services.
    *   **Example**: `eyJhbGciOiJIUzI1NiI...`

## Redis Configuration

*   **`REDIS_URL`**
    *   **Description**: The connection string for the Redis server. Used for caching, session management, and rate limiting.
    *   **Example**: `redis://localhost:6379`

## NATS Configuration

*   **`NATS_URL`**
    *   **Description**: The connection string for the NATS messaging server. Used for inter-service communication and event processing (e.g., by workers).
    *   **Example**: `nats://localhost:4222`

## JWT (JSON Web Token) Configuration

*   **`JWT_SECRET`**
    *   **Description**: A strong, randomly generated secret key used for signing and verifying JSON Web Tokens. **Crucial for security; must be kept confidential and strong in production.**
    *   **Example**: `acr3VZCLCbT1nYpnvMW6TCeZDQL6xSJdoG26j56/dcwmC2AdpcrH5V/exDy1UhhHuVy2/Cb8k4FdrY+5oqWIxg==`

*   **`ACCESS_TOKEN_TTL`**
    *   **Description**: Time-To-Live (TTL) for access tokens. Specifies how long an access token remains valid.
    *   **Example**: `15m` (15 minutes)

*   **`REFRESH_TOKEN_TTL`**
    *   **Description**: Time-To-Live (TTL) for refresh tokens. Specifies how long a refresh token remains valid, typically longer than access tokens.
    *   **Example**: `168h` (7 days)

## SMTP (Email Sending) Configuration

*   **`SMTP_HOST`**
    *   **Description**: The hostname of the SMTP server used for sending emails.
    *   **Example**: `smtp.gmail.com`

*   **`SMTP_PORT`**
    *   **Description**: The port of the SMTP server.
    *   **Example**: `587`

*   **`SMTP_USER`**
    *   **Description**: The username for SMTP authentication, typically the sender's email address.
    *   **Example**: `your_email@example.com`

*   **`SMTP_PASS`**
    *   **Description**: The password for SMTP authentication. For services like Gmail, this is often an app-specific password.
    *   **Example**: `your_app_password`

*   **`SMTP_FROM`**
    *   **Description**: The email address that will appear as the sender of outgoing emails.
    *   **Example**: `your_email@example.com`

## Magic Link Configuration

*   **`MAGIC_LINK_EXPIRY`**
    *   **Description**: The duration for which a magic link (used for passwordless login) remains valid.
    *   **Example**: `24h` (24 hours)

## Rate Limiting Configuration

*   **`RATE_LIMIT_ENABLED`**
    *   **Description**: A boolean flag to enable or disable global API rate limiting.
    *   **Example**: `true` or `false`

## Security Configuration

*   **`KEY_ROTATION_INTERVAL`**
    *   **Description**: The interval at which cryptographic keys (e.g., for JWTs) should be rotated. Important for security best practices.
    *   **Example**: `24h` (24 hours)
