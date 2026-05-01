package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Zordddd/learning/taskAPI/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestTaskHandler(t *testing.T) {
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:15",
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"POSTGRES_DB":       "postgres",
			},
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForLog("database system is ready to accept connections"),
			).WithDeadline(2 * time.Minute),
		},
		Started: true,
	})
	assert.NoError(t, err)
	defer container.Terminate(ctx)
	host, err := container.Host(ctx)
	assert.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	connStr := fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres?sslmode=disable", host, port.Port())
	pool, err := pgxpool.New(ctx, connStr)
	defer pool.Close()
	assert.NoError(t, err)

	schema := `
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			text TEXT NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`
	_, err = pool.Exec(ctx, schema)
	assert.NoError(t, err)

	repo := NewTaskRepositoryHandler(store.New(pool))

	t.Run("Wrong method test", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPatch, "/task", nil)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(repo.TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Get task handler test", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/task", nil)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(repo.TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create task handler test", func(t *testing.T) {
		body, _ := json.Marshal(store.CreateTaskParams{
			Title: "test",
			Text:  "create task",
		})
		r := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(repo.TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Create wrong task handler test", func(t *testing.T) {
		body, _ := json.Marshal(struct {
			Name int `json:"name"`
		}{
			Name: 0,
		})
		r := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(repo.TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update task handler test", func(t *testing.T) {
		testTaskCreate := store.CreateTaskParams{
			Title:     "test",
			Text:      "update task",
			Completed: pgtype.Bool{false, true},
		}
		testTaskJSON, _ := json.Marshal(testTaskCreate)
		rCreateTask := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(testTaskJSON))
		wCreateTask := httptest.NewRecorder()

		handler := http.HandlerFunc(repo.TaskHandler)
		handler.ServeHTTP(wCreateTask, rCreateTask)
		assert.Equal(t, http.StatusCreated, wCreateTask.Code)

		taskCreated := store.Task{}
		err = json.Unmarshal(wCreateTask.Body.Bytes(), &taskCreated)
		assert.NoError(t, err)

		newTask := store.UpdateTaskParams{
			ID:        taskCreated.ID,
			Title:     "test success",
			Text:      "update task",
			Completed: pgtype.Bool{true, true},
		}
		body, _ := json.Marshal(newTask)
		r := httptest.NewRequest(http.MethodPut, "/task", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update wrong task handler test", func(t *testing.T) {
		body, _ := json.Marshal(struct {
			Name int `json:"name"`
		}{
			Name: 0,
		})
		r := httptest.NewRequest(http.MethodPut, "/task", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(repo.TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Delete task handler test", func(t *testing.T) {
		testCreateTask := store.CreateTaskParams{
			Title:     "test",
			Text:      "create task",
			Completed: pgtype.Bool{false, false},
		}

		rCreateTaskBody, _ := json.Marshal(testCreateTask)
		rCreateTask := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(rCreateTaskBody))
		rCreateTask.Header.Set("Content-Type", "application/json")
		wCreateTask := httptest.NewRecorder()

		handler := http.HandlerFunc(repo.TaskHandler)

		handler.ServeHTTP(wCreateTask, rCreateTask)
		assert.Equal(t, http.StatusCreated, wCreateTask.Code)
		task := store.Task{}
		err = json.Unmarshal(wCreateTask.Body.Bytes(), &task)
		assert.NoError(t, err)
		taskID := task.ID

		rDeleteTask := httptest.NewRequest("DELETE", fmt.Sprintf("/task?id=%d", taskID), nil)
		rDeleteTask.Header.Set("Content-Type", "application/json")
		wDeleteTask := httptest.NewRecorder()

		handler.ServeHTTP(wDeleteTask, rDeleteTask)
		assert.Equal(t, http.StatusOK, wDeleteTask.Code)

	})
}
