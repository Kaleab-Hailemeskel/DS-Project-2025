package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"song-service/api/config"
	"song-service/api/internal/domain"
	"song-service/api/internal/repository"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	redisSearchRepo repository.IRedisSearchRepo
}

func (m *Middleware) AuthUser(ctx *gin.Context) {
	log.Printf("[AuthUser] Processing request: %s %s", ctx.Request.Method, ctx.Request.URL.Path)

	// 1. Extract Token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("[AuthUser] Error: Missing or malformed Authorization header")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No token provided"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 2. Check Redis Cache first
	cachedUserExist, err := m.redisSearchRepo.IsTokenValid(ctx, token)
	if err == nil && cachedUserExist == true {
		log.Println("[AuthUser] Cache Hit: Token is valid in Redis")
		ctx.Next()
		return
	}
	log.Println("[AuthUser] Cache Miss: Validating token via User Service")

	// 3. Fallback: Validate via User Service
	payload := domain.AuthRequest{
		AccessToken: token,
		Context:     "auth_validation",
	}
	jsonData, _ := json.Marshal(payload)

	url := fmt.Sprintf("http://%s:%s/validate", config.USER_SERVICE_URL, config.USER_SERVICE_PORT)
	log.Printf("[AuthUser] Calling User Service: %s", url)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[AuthUser] Error: Failed to reach User Service: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "auth service unavailable"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[AuthUser] Error: User Service returned status %d", resp.StatusCode)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var result domain.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[AuthUser] Error decoding User Service response: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	// 4. Check Business Logic Result
	if !result.Valid || !result.Data.Linked {
		log.Printf("[AuthUser] Validation failed: Valid=%v, Linked=%v, Message=%s", result.Valid, result.Data.Linked, result.Message)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// 5. Calculate dynamic TTL based on ExpiresAt
	if duration := time.Until(result.Data.ExpiresAt); duration > 0 {
		log.Printf("[AuthUser] Success: Caching token in Redis for %v", duration)
		_ = m.redisSearchRepo.MarkTokenAsValid(ctx, token, duration)
	} else {
		log.Println("[AuthUser] Warning: Token valid but ExpiresAt is in the past, skipping cache")
	}

	// 6. Continue
	log.Println("[AuthUser] Authentication successful")
	ctx.Next()
}

func NewMiddleware(redisSearchRepo_ repository.IRedisSearchRepo) IMiddleware {
	return &Middleware{
		redisSearchRepo: redisSearchRepo_,
	}
}