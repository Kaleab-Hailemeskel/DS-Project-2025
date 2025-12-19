package main

import (
	"context"
	"net/http"
	"strings"
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

// GetUserIDFromContext retrieves userID set by middleware.
func GetUserIDFromContext(ctx context.Context) (int, bool) {
    v := ctx.Value(userIDKey)
    if v == nil {
        return 0, false
    }
    id, ok := v.(int)
    return id, ok
}
