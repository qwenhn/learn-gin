package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/demo", demoHandler)

	log.Println("Server is starting...")
	err := http.ListenAndServe(":8080", nil) // localhost:8080
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

func demoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)

	if r.Method != http.MethodGet {
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	resp := map[string]string{
		"message": "Welcome you",
		"info":    "Qwen",
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-API-Key", "123456")

	json.NewEncoder(w).Encode(resp)
}
