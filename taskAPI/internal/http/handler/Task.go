package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Zordddd/learning/taskAPI/internal/storage"
)

// TaskHandler godoc
// @Summary Main task handler
// @Description Routes task requests to appropriate handlers based on HTTP method
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} storage.Task "GET method response"
// @Success 201 {object} storage.Task "POST method response"
// @Success 200 {object} map[string]string "PUT/DELETE method response"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 405 {object} map[string]string "Method not allowed"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tasks [get]
// @Router /tasks [post]
// @Router /tasks [put]
// @Router /tasks [delete]
// @Security ApiKeyAuth
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
	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetTasksHandler godoc
// @Summary Get all tasks
// @Description Retrieve all tasks from the database
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} storage.Task
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
// @Security ApiKeyAuth
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	result := make([]storage.Task, 0, len(storage.Database.Tasks))
	storage.Database.Mu.RLock()
	for _, task := range storage.Database.Tasks {
		result = append(result, *task)
	}
	storage.Database.Mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateTaskHandler godoc
// @Summary Create a new task
// @Description Add a new task to the database
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body storage.Task true "Task object"
// @Success 201 {object} storage.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
// @Security ApiKeyAuth
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task storage.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	storage.Database.Mu.Lock()
	task.ID = storage.Database.NextID
	task.Timestamp = time.Now()
	storage.Database.Tasks[storage.Database.NextID] = &task
	storage.Database.NextID++
	storage.Database.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UpdateTaskHandler godoc
// @Summary Update an existing task
// @Description Update task information
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body storage.Task true "Updated task object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [put]
// @Security ApiKeyAuth
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var currentTask storage.Task
	if err := json.NewDecoder(r.Body).Decode(&currentTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := storage.Database.Tasks[currentTask.ID]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	storage.Database.Mu.Lock()
	storage.Database.Tasks[currentTask.ID] = &currentTask
	storage.Database.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": "success",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteTaskHandler godoc
// @Summary Delete a task
// @Description Remove a task from the database by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id query int true "Task ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [delete]
// @Security ApiKeyAuth
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("id")
	if data == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, exists := storage.Database.Tasks[id]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	storage.Database.Mu.Lock()
	delete(storage.Database.Tasks, id)
	storage.Database.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": "success",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
