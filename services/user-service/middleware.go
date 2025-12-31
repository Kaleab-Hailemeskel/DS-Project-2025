package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware returns a middleware that validates Bearer token and injects userID into context.
func AuthMiddleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing authorization", http.StatusUnauthorized)
			return
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}
		token := parts[1]
		uid, err := ParseToken(secret, token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateUser(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req AuthRequest
		// 1. Process the request JSON
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Request processing failed (e.g., malformed JSON)
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

		// 2. Authorize using the Parsing logic
		// We use req.AccessToken from the parsed JSON body
		_, err := ParseToken(secret, req.AccessToken)

		// Create the detail data with default 30 min expiration
		detail := AuthDetailData{
			ExpiresAt: time.Now().Add(30 * time.Minute),
		}

		if err != nil {
			// Request was processed (Valid: true), but authorization failed (Linked: false)
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
		detail.Linked = true

		
	})
}

// GetUserIDFromContext retrieves userID set by middleware.
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	v := ctx.Value(userIDKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(int)
	return id, ok
}
