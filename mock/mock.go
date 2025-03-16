package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var mockResponse = map[string]interface{}{
	"releaseDate": "16.07.2006",
	"text":        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?",
	"link":        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
}

func main() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Mock API server is running on http://localhost:8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Failed to start mock API server: %v", err)
	}
}