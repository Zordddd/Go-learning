package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Zordddd/learning/taskAPI/internal/storage"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	result := make([]storage.Task, 0, len(storage.Database.Tasks))
	for _, task := range storage.Database.Tasks {
		result = append(result, task)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
