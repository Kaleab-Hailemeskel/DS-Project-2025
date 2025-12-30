package domain

import "time"

type AuthRequest struct {
	AccessToken string `json:"access_token"`
	Context     string `json:"context"`
}

type AuthResponse struct {
	Valid   bool           `json:"valid"`
	Message string         `json:"message"`
	Data    AuthDetailData `json:"data"`
}

type AuthDetailData struct {
	Linked    bool      `json:"linked"`
	ExpiresAt time.Time `json:"expires_at"`
}
