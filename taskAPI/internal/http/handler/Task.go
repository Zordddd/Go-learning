package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Zordddd/learning/taskAPI/internal/storage"
	"github.com/Zordddd/learning/taskAPI/pkg/http/responseWriter"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	Method := r.Method
	switch Method {
	case http.MethodGet:
		GetTasksHandler(w, r)
	case http.MethodPost:
		CreateTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	}
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	rw := responseWriter.NewResponseWriter(w)
	result := make([]storage.Task, 0, len(storage.Database.Tasks))
	storage.Database.Mu.RLock()
	for _, task := range storage.Database.Tasks {
		result = append(result, *task)
	}
	storage.Database.Mu.RUnlock()

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(result); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	rw := responseWriter.NewResponseWriter(w)
	var task storage.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	storage.Database.Mu.Lock()
	task.ID = storage.Database.NextID
	task.Timestamp = time.Now()
	storage.Database.Tasks[storage.Database.NextID] = &task
	storage.Database.NextID++
	storage.Database.Mu.Unlock()

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	rw := responseWriter.NewResponseWriter(w)
	var currentTask storage.Task
	if err := json.NewDecoder(r.Body).Decode(&currentTask); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	storage.Database.Mu.Lock()
	storage.Database.Tasks[currentTask.ID] = &currentTask
	storage.Database.Mu.Unlock()

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": "success",
	}
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	rw := responseWriter.NewResponseWriter(w)
	data := r.URL.Query().Get("id")
	if data == "" {
		http.Error(rw, "id is required", http.StatusBadRequest)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	storage.Database.Mu.Lock()
	delete(storage.Database.Tasks, id)
	storage.Database.Mu.Unlock()

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": "success",
	}
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
