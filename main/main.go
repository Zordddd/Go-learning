package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/health", getHealthHandler)
	http.HandleFunc("/user/", getUserHandler) // Изменил на /user/
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/search", searchHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func getHealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	response := map[string]interface{}{
		"status":    "OK",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из пути: /user/123 → "123"
	path := strings.TrimPrefix(r.URL.Path, "/user/")
	if path == "" {
		http.Error(w, `{"error": "User ID required"}`, http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"user_id": path,
		"name":    "John",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"received_data": data,
		"status":        "processed",
		"timestamp":     time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	response := map[string]interface{}{
		"search_query": params.Get("q"),
		"limit":        params.Get("limit"),
		"all_params":   params,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
