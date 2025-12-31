package main

import (
	"database/sql"
	"encoding/json"
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
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
        return
    }
    if req.Username == "" || req.Password == "" {
        writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "username and password required"})
        return
    }
    hash, err := hashPassword(req.Password, a.bcCost)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to hash password"})
        return
    }
    ctx := r.Context()
    id, created, err := a.store.CreateUser(ctx, req.Username, req.Email, string(hash))
    if err != nil {
        if err == sql.ErrNoRows {
            writeJSON(w, http.StatusConflict, ErrorResponse{Error: "user exists"})
            return
        }
        writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "create user failed"})
        return
    }
    u := &User{ID: id, Username: req.Username, Email: req.Email, CreatedAt: created}
    writeJSON(w, http.StatusCreated, u)
}

func (a *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid json"})
        return
    }
    if req.UsernameOrEmail == "" || req.Password == "" {
        writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing fields"})
        return
    }
    user, err := a.store.GetUserByUsernameOrEmail(r.Context(), req.UsernameOrEmail)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
        return
    }
    if err := checkPassword([]byte(user.PasswordHash), req.Password); err != nil {
        writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
        return
    }
    token, exp, err := GenerateToken(a.jwtSecret, a.jwtExpiry, user.ID)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "token generation failed"})
        return
    }
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
        writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
        return
    }
    user, err := a.store.GetUserByID(r.Context(), uid)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch user"})
        return
    }
    // hide password hash
    user.PasswordHash = ""
    writeJSON(w, http.StatusOK, user)
}

// helper to read int env with default
func mustAtoi(env string, def int) int {
    if v := os.Getenv(env); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            return n
        }
    }
    return def
}
