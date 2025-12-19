package main

import "time"

// User represents a user record stored in Postgres.
type User struct {
    ID           int       `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

// RegisterRequest is the JSON payload for user registration.
type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

// LoginRequest is the JSON payload for login.
type LoginRequest struct {
    UsernameOrEmail string `json:"username_or_email"`
    Password        string `json:"password"`
}

// LoginResponse is returned after successful login.
type LoginResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
    User        *User  `json:"user"`
}

// ErrorResponse is a simple JSON error.
type ErrorResponse struct {
    Error string `json:"error"`
}
