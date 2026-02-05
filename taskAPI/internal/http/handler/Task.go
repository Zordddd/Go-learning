package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Zordddd/learning/taskAPI/internal/models"
	"github.com/Zordddd/learning/taskAPI/internal/postgres"
)

type TaskRepositoryHandler struct {
	Repo *postgres.TaskRepository
}

func NewTaskRepositoryHandler(repo *postgres.TaskRepository) *TaskRepositoryHandler {
	return &TaskRepositoryHandler{Repo: repo}
}

func (t *TaskRepositoryHandler) PingContext(ctx context.Context) error {
	err := t.Repo.DB.PingContext(ctx)
	return err
}

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
func (t *TaskRepositoryHandler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	Method := r.Method
	switch Method {
	case http.MethodGet:
		t.GetTasksHandler(w, r)
	case http.MethodPost:
		t.CreateTaskHandler(w, r)
	case http.MethodPut:
		t.UpdateTaskHandler(w, r)
	case http.MethodDelete:
		t.DeleteTaskHandler(w, r)
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
func (t *TaskRepositoryHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.Repo.GetTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
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
func (t *TaskRepositoryHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.TaskCreated
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := t.Repo.CreateTask(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

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
func (t *TaskRepositoryHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var currentTask models.TaskUpdated
	if err := json.NewDecoder(r.Body).Decode(&currentTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := t.Repo.UpdateTask(r.Context(), currentTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

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
func (t *TaskRepositoryHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	err = t.Repo.DeleteTask(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": "success",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
