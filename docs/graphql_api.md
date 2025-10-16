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

### `logout: Boolean!`

Logs out the current user by invalidating their session and tokens.

- **Input:** None
- **Output:** `Boolean!`
    - `true` if logout was successful, `false` otherwise.

### `recoverPassword(input: RecoverPasswordInput!): Boolean!`

Initiates the password recovery process for a given email.

- **Input:** `RecoverPasswordInput`
    - `email`: User's email address (String!)
- **Output:** `Boolean!`
    - `true` if recovery email was sent successfully, `false` otherwise.

### `deleteAccount(input: DeleteAccountInput!): Boolean!`

Schedules the authenticated user's account for deletion.

- **Input:** `DeleteAccountInput`
    - `userID`: ID of the user to be deleted (ID!)
- **Output:** `Boolean!`
    - `true` if deletion was scheduled successfully, `false` otherwise.

### `recoverAccount(input: RecoverAccountInput!): AuthResponse!`

Recovers a user account using a recovery token and sets a new password.

- **Input:** `RecoverAccountInput`
    - `token`: The recovery token (String!)
    - `newPassword`: The new password (String!)
- **Output:** `AuthResponse`
    - `accessToken`: JWT access token (String!)
    - `refreshToken`: Refresh token (String!)
    - `user`: The recovered user (User!)

## Types

### `AuthResponse`

Represents the response after successful authentication or registration.

- `accessToken`: String!
- `refreshToken`: String!
- `user`: User!



### `User`

Represents a user in the system.

- `id`: ID!
- `name`: String!
- `email`: String!
- `createdAt`: String!
- `updatedAt`: String!
- `isEmailVerified`: Boolean!
- `deletedAt`: String
- `avatarURL`: String
- `deletionDueAt`: String
- `lastLoginAt`: String
- `isDeleted`: Boolean!

## Input Objects

### `RegisterUserInput`

Input for the `registerUser` mutation.

- `name`: String!
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

- `userID`: ID!

### `RecoverAccountInput`

Input for the `recoverAccount` mutation.

- `token`: String!
- `newPassword`: String!
