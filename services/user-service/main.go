package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User service is running!")
	})
	fmt.Println("User service running on port 8080")
	http.ListenAndServe(":8080", nil)
}

