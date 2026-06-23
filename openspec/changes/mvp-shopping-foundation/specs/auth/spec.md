# Auth Specification

## Purpose

Authenticate the two household users and protect API routes with JWT.

## Requirements

### Requirement: Fixed Users

The system SHALL seed exactly two users on first startup: one primary user and one partner user. Credentials SHALL come from environment variables.

#### Scenario: Login with valid credentials

- GIVEN the seeded users exist
- WHEN a user sends `POST /api/v1/auth/login` with valid email and password
- THEN the response SHALL contain a JWT and a 200 status
- AND the token SHALL expire after `JWT_EXPIRATION_HOURS`

#### Scenario: Login with invalid credentials

- GIVEN the seeded users exist
- WHEN a user sends `POST /api/v1/auth/login` with an invalid password
- THEN the response SHALL return 401 with code `UNAUTHORIZED`
- AND the response SHALL NOT contain a token

### Requirement: Token Protection

The system SHALL reject requests to protected routes when the JWT is missing, expired, or invalid.

#### Scenario: Access protected route without token

- GIVEN a protected endpoint `/api/v1/products`
- WHEN a request is sent without an `Authorization` header
- THEN the response SHALL return 401 with code `UNAUTHORIZED`

#### Scenario: Access protected route with expired token

- GIVEN a valid token that has expired
- WHEN a request is sent with that token
- THEN the response SHALL return 401 with code `TOKEN_EXPIRED`

## Notes

- No registration endpoint in MVP.
- Passwords SHALL be hashed with bcrypt.
