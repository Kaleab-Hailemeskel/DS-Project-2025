package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// envs
	pg := os.Getenv("POSTGRES_URL")
	if pg == "" {
		log.Fatal("POSTGRES_URL is required")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("WARNING: JWT_SECRET is not set — tokens will not work in production")
	}
	jwtExpiry := 60
	if v := os.Getenv("JWT_EXPIRY_MINUTES"); v != "" {
		// ignore error, handlers.may use default
	}

	// open DB
	db, err := sql.Open("postgres", pg)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	store := NewStore(db)

	app := &App{
		store:     store,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
		bcCost:    mustAtoi("BCRYPT_COST", 12),
	}
	// enbable CORS for all routes
	http.Handle("/", enableCORS(http.DefaultServeMux))

	// routes (gateway will proxy /api/users/ -> this service)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User service is running!")
	})

	
	http.HandleFunc("/register", app.RegisterHandler)
	http.HandleFunc("/login", app.LoginHandler)

	// protected /me with middleware
	meHandler := http.HandlerFunc(app.MeHandler)
	http.Handle("/me", AuthMiddleware(jwtSecret, meHandler))
	// validation function for song-service | streaming-service
	http.Handle("/validate", ValidateUser(jwtSecret, nil))

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	fmt.Println("User service running on port", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Allow any origin, or replace "*" with your specific frontend URL
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight OPTIONS request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}