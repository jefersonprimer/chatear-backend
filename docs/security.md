# Security Design Document

## Authentication and Token Logic

This document outlines the security design for authentication and token management within the system.

### 1. Access Tokens (JWT)
- **Type:** JSON Web Tokens (JWT)
- **Purpose:** Used for authenticating API requests.
- **Lifespan:** Short-lived to minimize the impact of token compromise.
- **Storage:** Stored client-side (e.g., in memory, HTTP-only cookies).
- **Validation:** Validated on each protected API request.

### 2. Refresh Tokens
- **Purpose:** Used to obtain new access tokens without requiring the user to re-authenticate.
- **Storage:** Stored securely in the PostgreSQL database in a dedicated `refresh_tokens` table.
- **Lifespan:** Longer-lived than access tokens.
- **Usage:** Exchanged for a new access token and refresh token pair.
- **Invalidation:** Can be revoked by the user (e.g., logout from all devices) or by the system.

### 3. Blacklisted/Revoked Tokens
- **Purpose:** To immediately invalidate compromised or explicitly revoked access/refresh tokens.
- **Storage:** Stored in Redis with a Time-To-Live (TTL) corresponding to the token's original expiry.
- **Mechanism:** When a token is blacklisted, its signature is added to a Redis set or hash. During token validation, this list is checked.

### 4. Token Creation, Parsing, and Validation Helpers
- **Location:** `shared/auth/` directory.
- **Functionality:**
    - Generate new JWT access tokens (`auth.GenerateAccessToken`).
    - Generate new refresh tokens (`auth.GenerateRefreshToken`).
    - Parse and validate JWT access tokens (`auth.ValidateAccessToken`).
    - Blacklist access tokens (`blacklistRepo.Add`).

### 5. Key Management
- **JWT Secret:** A strong, randomly generated secret key used for signing JWTs. Defined as `constants.JwtSecret`.
- **Access Token Expiration:** Configured via `constants.AccessTokenExpiration`.
- **Storage:** Environment variable or secure secret management system.

### 6. Security Considerations
- **HTTPS:** All communication must occur over HTTPS.
- **CSRF Protection:** Implement CSRF protection for state-changing requests.
- **XSS Protection:** Sanitize all user-generated content.
- **Rate Limiting:** Apply rate limiting to authentication endpoints to prevent brute-force attacks.
- **Secure Cookie Flags:** Use `HttpOnly`, `Secure`, and `SameSite` flags for cookies storing tokens.
