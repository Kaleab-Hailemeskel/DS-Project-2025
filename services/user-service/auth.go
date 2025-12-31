package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string, cost int) ([]byte, error) {
    return bcrypt.GenerateFromPassword([]byte(password), cost)
}

func checkPassword(hash []byte, password string) error {
    return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

func getBcryptCost() int {
    if v := os.Getenv("BCRYPT_COST"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n >= 4 {
            return n
        }
    }
    return 12
}

// GenerateToken creates a signed JWT with user id as `sub` and expiry in minutes.
func GenerateToken(secret string, expiryMinutes int, userID int) (string, time.Time, error) {
    if secret == "" {
        return "", time.Time{}, errors.New("missing JWT secret")
    }
    now := time.Now().UTC()
    exp := now.Add(time.Duration(expiryMinutes) * time.Minute)
    claims := jwt.MapClaims{
        "sub": fmt.Sprintf("%d", userID),
        "iat": now.Unix(),
        "exp": exp.Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(secret))
    return signed, exp, err
}

// ParseToken verifies a token and returns the user id (as int) from `sub` claim.
func ParseToken(secret, tokenStr string) (int, error) {
    if secret == "" {
        return 0, errors.New("missing JWT secret")
    }
    token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(secret), nil
    })
    if err != nil {
        return 0, err
    }
    if !token.Valid {
        return 0, errors.New("invalid token")
    }
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        if sub, found := claims["sub"]; found {
            switch v := sub.(type) {
            case string:
                var id int
                _, err := fmt.Sscanf(v, "%d", &id)
                if err != nil {
                    return 0, err
                }
                return id, nil
            case float64:
                return int(v), nil
            default:
                return 0, errors.New("unexpected sub claim type")
            }
        }
    }
    return 0, errors.New("sub claim not found")
}
