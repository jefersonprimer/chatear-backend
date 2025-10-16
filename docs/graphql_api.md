# GraphQL API Documentation

This document describes the GraphQL API structure, available queries, mutations, and types.

## Schema Overview

The GraphQL API provides a single endpoint for interacting with the system. It is built using `gqlgen` and follows standard GraphQL conventions.

## Authentication

All protected mutations and queries require a valid JWT access token to be sent in the `Authorization` header as a Bearer token.

## Mutations

### `registerUser(input: RegisterUserInput!): AuthResponse!`

Registers a new user with the provided email and password.

- **Input:** `RegisterUserInput`
    - `email`: User's email address (String!)
    - `password`: User's password (String!)
- **Output:** `AuthResponse`
    - `accessToken`: JWT access token (String!)
    - `refreshToken`: Refresh token (String!)
    - `user`: The newly registered user (User!)

### `login(input: LoginInput!): AuthResponse!`

Authenticates a user and returns access and refresh tokens.

- **Input:** `LoginInput`
    - `email`: User's email address (String!)
    - `password`: User's password (String!)
- **Output:** `AuthResponse`
    - `accessToken`: JWT access token (String!)
    - `refreshToken`: Refresh token (String!)
    - `user`: The logged-in user (User!)

### `logout: MessageResponse!`

Logs out the current user by invalidating their session and tokens.

- **Input:** None
- **Output:** `MessageResponse`
    - `message`: A confirmation message (String!)

### `recoverPassword(input: RecoverPasswordInput!): MessageResponse!`

Initiates the password recovery process for a given email.

- **Input:** `RecoverPasswordInput`
    - `email`: User's email address (String!)
- **Output:** `MessageResponse`
    - `message`: A confirmation message (String!)

### `deleteAccount(input: DeleteAccountInput!): MessageResponse!`

Deletes the authenticated user's account.

- **Input:** `DeleteAccountInput`
    - `password`: Current password for confirmation (String!)
- **Output:** `MessageResponse`
    - `message`: A confirmation message (String!)

### `recoverAccount(input: RecoverAccountInput!): MessageResponse!`

Recovers a user account using a recovery token and sets a new password.

- **Input:** `RecoverAccountInput`
    - `token`: The recovery token (String!)
    - `newPassword`: The new password (String!)
- **Output:** `MessageResponse`
    - `message`: A confirmation message (String!)

## Types

### `AuthResponse`

Represents the response after successful authentication or registration.

- `accessToken`: String!
- `refreshToken`: String!
- `user`: User!

### `MessageResponse`

A generic response type for operations that return a simple message.

- `message`: String!

### `User`

Represents a user in the system.

- `id`: ID!
- `email`: String!

## Input Objects

### `RegisterUserInput`

Input for the `registerUser` mutation.

- `email`: String!
- `password`: String!

### `LoginInput`

Input for the `login` mutation.

- `email`: String!
- `password`: String!

### `RecoverPasswordInput`

Input for the `recoverPassword` mutation.

- `email`: String!

### `DeleteAccountInput`

Input for the `deleteAccount` mutation.

- `password`: String!

### `RecoverAccountInput`

Input for the `recoverAccount` mutation.

- `token`: String!
- `newPassword`: String!
