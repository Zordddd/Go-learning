package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Zordddd/learning/taskAPI/internal/storage"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task storage.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	storage.Database.Mu.Lock()
	task.ID = storage.Database.NextID
	task.Timestamp = time.Now()
	storage.Database.Tasks[storage.Database.NextID] = task
	storage.Database.NextID++
	storage.Database.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
