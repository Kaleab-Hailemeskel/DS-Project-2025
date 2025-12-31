package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"streaming-service/config"
	"streaming-service/internal/domain"

	"streaming-service/internal/repository"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	redisSearchRepo repository.IRedisSearchRepo
}

func (m *Middleware) AuthUser(ctx *gin.Context) {
	// 1. Extract Token
	authHeader := ctx.GetHeader("Authorization")
	log.Printf("✅ AuthHeader: %s", authHeader)
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No token provided"})
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 2. Check Redis Cache first
	cachedUserExist, err := m.redisSearchRepo.IsTokenValid(ctx, token)
	if err == nil && cachedUserExist == true {
		ctx.Next()
		return
	}

	// 3. Fallback: Validate via User Service
	payload := domain.AuthRequest{
		AccessToken: token,
		Context:     "auth_validation",
	}
	jsonData, _ := json.Marshal(payload)

	url := fmt.Sprintf("http://%s:%s/validate", config.USER_SERVICE_URL, config.USER_SERVICE_PORT)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil || resp.StatusCode != http.StatusOK {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "auth service unavailable"})
		return
	}
	defer resp.Body.Close()

	var result domain.AuthResponse
	json.NewDecoder(resp.Body).Decode(&result)
	log.Printf("✅ AuthResponse: %+v", result)
	// 4. Check Business Logic Result
	if !result.Valid || !result.Data.Linked {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// 5. Calculate dynamic TTL based on ExpiresAt
	// If the token is already expired or expires very soon, we set a small default or don't cache
	if duration := time.Until(result.Data.ExpiresAt); duration > 0 {
		_ = m.redisSearchRepo.MarkTokenAsValid(ctx, token, duration)
	}

	// 6. Continue
	ctx.Next()
}

func NewMiddleware(redisSearchRepo_ repository.IRedisSearchRepo) IMiddleware {
	return &Middleware{
		redisSearchRepo: redisSearchRepo_,
	}
}
