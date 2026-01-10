package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type App struct {
	store     *Store
	jwtSecret string
	jwtExpiry int
	bcCost    int
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

	func (a *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Register] Received request: %s %s", r.Method, r.URL.Path)

		if r.Method != http.MethodPost {
			log.Printf("[Register] Error: Method %s not allowed", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[Register] Error decoding JSON: %v", err)
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
			return
		}

		log.Printf("[Register] Attempting to register user: %s (Email: %s)", req.Username, req.Email)

		if req.Username == "" || req.Password == "" {
			log.Println("[Register] Error: Missing username or password")
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "username and password required"})
			return
		}

		hash, err := hashPassword(req.Password, a.bcCost)
		if err != nil {
			log.Printf("[Register] Error hashing password: %v", err)
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to hash password"})
			return
		}

		ctx := r.Context()
		id, created, err := a.store.CreateUser(ctx, req.Username, req.Email, string(hash))
		if err != nil {
			log.Printf("[Register] Store.CreateUser failed for %s: %v", req.Username, err)
			if err == sql.ErrNoRows {
				writeJSON(w, http.StatusConflict, ErrorResponse{Error: "user exists"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "create user failed"})
			return
		}

		log.Printf("[Register] Success: Created user %s with ID %d", req.Username, id)
		u := &User{ID: id, Username: req.Username, Email: req.Email, CreatedAt: created}
		writeJSON(w, http.StatusCreated, u)
	}

func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Login] Received request: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("[Login] Error: Method %s not allowed", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[Login] Error decoding JSON: %v", err)
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
		return
	}

	log.Printf("[Login] Attempting login for: %s", req.UsernameOrEmail)

	if req.UsernameOrEmail == "" || req.Password == "" {
		log.Println("[Login] Error: Missing fields")
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing fields"})
		return
	}

	user, err := a.store.GetUserByUsernameOrEmail(r.Context(), req.UsernameOrEmail)
	if err != nil {
		log.Printf("[Login] Store.GetUserByUsernameOrEmail failed for %s: %v", req.UsernameOrEmail, err)
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	if err := checkPassword([]byte(user.PasswordHash), req.Password); err != nil {
		log.Printf("[Login] Password verification failed for user %d: %v", user.ID, err)
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	token, exp, err := GenerateToken(a.jwtSecret, a.jwtExpiry, user.ID)
	if err != nil {
		log.Printf("[Login] Token generation failed for user %d: %v", user.ID, err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "token generation failed"})
		return
	}

	log.Printf("[Login] Success: User %d logged in, token expires at %v", user.ID, exp)
	resp := LoginResponse{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   int(time.Until(exp).Seconds()),
		User:        user,
	}
	writeJSON(w, http.StatusOK, resp)
}

func (a *App) MeHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := GetUserIDFromContext(r.Context())
	if !ok {
		log.Println("[Me] Error: No UserID found in context")
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	log.Printf("[Me] Fetching profile for UserID: %d", uid)

	user, err := a.store.GetUserByID(r.Context(), uid)
	if err != nil {
		log.Printf("[Me] Store.GetUserByID failed for ID %d: %v", uid, err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch user"})
		return
	}

	user.PasswordHash = ""
	log.Printf("[Me] Success: Returning profile for %s", user.Username)
	writeJSON(w, http.StatusOK, user)
}

func mustAtoi(env string, def int) int {
	if v := os.Getenv(env); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}