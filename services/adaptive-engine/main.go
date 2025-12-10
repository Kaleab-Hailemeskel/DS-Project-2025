package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Adaptive engine running!")
	})
	fmt.Println("Adaptive engine running on port 8080")
	http.ListenAndServe(":8080", nil)
}
