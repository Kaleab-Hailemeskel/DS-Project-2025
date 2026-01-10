package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware returns a middleware that validates Bearer token and injects userID into context.
func AuthMiddleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[AuthMiddleware] Intercepting request: %s %s", r.Method, r.URL.Path)

		auth := r.Header.Get("Authorization")
		if auth == "" {
			log.Println("[AuthMiddleware] Error: Missing Authorization header")
			http.Error(w, "missing authorization", http.StatusUnauthorized)
			return
		}

		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Printf("[AuthMiddleware] Error: Invalid authorization format. Parts count: %d", len(parts))
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		uid, err := ParseToken(secret, token)
		if err != nil {
			log.Printf("[AuthMiddleware] Error: Token validation failed: %v", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("[AuthMiddleware] Success: Valid token for UserID %d", uid)
		ctx := context.WithValue(r.Context(), userIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateUser(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[ValidateUser] Received validation request")
		w.Header().Set("Content-Type", "application/json")

		var req AuthRequest
		// 1. Process the request JSON
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[ValidateUser] Error: Failed to decode request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(AuthResponse{
				Valid:   false,
				Message: "Invalid request format",
				Data: AuthDetailData{
					Linked:    false,
					ExpiresAt: time.Now().Add(30 * time.Minute),
				},
			})
			return
		}

		log.Printf("[ValidateUser] Attempting to validate AccessToken for external service")

		// 2. Authorize using the Parsing logic
		uid, err := ParseToken(secret, req.AccessToken)

		detail := AuthDetailData{
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		if err != nil {
			log.Printf("[ValidateUser] Validation failed: %v", err)
			detail.Linked = false
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(AuthResponse{
				Valid:   true,
				Message: "Authorization failed: " + err.Error(),
				Data:    detail,
			})
			return
		}

		// 3. Success Case
		log.Printf("[ValidateUser] Success: Token is valid for UserID %d", uid)
		detail.Linked = true

		// If a next handler is provided, we pass the user ID down via context and continue.
		// This allows you to chain further logic after validation.
		if next != nil {
			log.Println("[ValidateUser] Forwarding to next handler")
			ctx := context.WithValue(r.Context(), userIDKey, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		// If no next handler, we simply return the success JSON to the caller
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(AuthResponse{
			Valid:   true,
			Message: "Token validated successfully",
			Data:    detail,
		})
	})
}

// GetUserIDFromContext retrieves userID set by middleware.
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	v := ctx.Value(userIDKey)
	if v == nil {
		log.Println("[Context] Warning: No UserID found in context")
		return 0, false
	}
	id, ok := v.(int)
	if !ok {
		log.Printf("[Context] Error: Value in context is not an int: %T", v)
	}
	return id, ok
}
