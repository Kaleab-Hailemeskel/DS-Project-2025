package main
 // Simple streaming service main file
import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "streaming service is running!")
	})
	fmt.Println("Song service running on port 8080")
	http.ListenAndServe(":8080", nil)
}
