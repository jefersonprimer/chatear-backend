# Insomnia Test Guide for Chatear Backend

This guide provides instructions on how to set up and test the Chatear Backend using Insomnia.

## 1. Setup and Start the Services

Before running the application, ensure you have Docker and Docker Compose installed.

### 1.1. Environment Variables

Create a `.env` file in the root of the project based on `env.example`. Make sure to fill in all the necessary details, especially for `DATABASE_URL`, `REDIS_URL`, `NATS_URL`, and `SMTP` settings.

```bash
cp env.example .env
# Open .env and fill in your details
```

**Important:** For `SMTP_USER` and `SMTP_PASS`, if you are using Gmail, you will need to generate an App Password. Refer to Google's documentation on how to do this.

### 1.2. Start Docker Services

Navigate to the root of your project and start the Docker services (Redis, NATS):

```bash
docker-compose -f docker-compose.events.yml up -d
```

This will start the `redis`, and `nats` containers in detached mode.

### 1.3. Run the Go Application

You can run the API and worker services separately.

#### 1.3.1. Run the API Service

Open a new terminal in the project root and run the API service:

```bash
go run cmd/api/main.go
```

#### 1.3.2. Run the Notification Worker Service

Open another new terminal in the project root and run the notification worker service:

```bash
go run cmd/worker/notification_worker.go
```

## 2. Insomnia Requests

You can import the following requests into Insomnia to test the application.

### 2.1. Register User

**Endpoint:** `POST http://localhost:8080/api/v1/register`

**Headers:**
*   `Content-Type: application/json`

**Body (JSON):**

```json
{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123"
}
```

**Expected Response:**
A successful registration will return a user object.

### 2.2. Login User

**Endpoint:** `POST http://localhost:8080/api/v1/login`

**Headers:**
*   `Content-Type: application/json`

**Body (JSON):**

```json
{
    "email": "test@example.com",
    "password": "password123"
}
```

**Expected Response:**
A successful login will return `accessToken` and `refreshToken`.

### 2.3. Password Recovery (Triggers Email Notification)

**Endpoint:** `POST http://localhost:8080/api/v1/password-recovery`

**Headers:**
*   `Content-Type: application/json`

**Body (JSON):**

```json
{
    "email": "test@example.com"
}
```

**Expected Response:**
A successful request will return a success message. Check the logs of your `notification_worker` for email sending attempts. If your SMTP settings are correct, you should receive an email.

### 2.4. Get Current User (Authenticated)

**Endpoint:** `GET http://localhost:8080/api/v1/me`

**Headers:**
*   `Authorization: Bearer <YOUR_ACCESS_TOKEN>` (Replace `<YOUR_ACCESS_TOKEN>` with the token obtained from login)

**Expected Response:**
Returns the details of the currently authenticated user.

### 2.5. Logout User (Authenticated)

**Endpoint:** `POST http://localhost:8080/api/v1/logout`

**Headers:**
*   `Authorization: Bearer <YOUR_ACCESS_TOKEN>` (Replace `<YOUR_ACCESS_TOKEN>` with the token obtained from login)

**Expected Response:**
A successful logout will return a success message. The access token will be blacklisted in Redis.

## 3. Clean Up

To stop and remove the Docker containers:

```bash
docker-compose -f docker-compose.events.yml down
```
