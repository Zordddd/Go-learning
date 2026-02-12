package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Zordddd/learning/taskAPI/internal/config"
	"github.com/Zordddd/learning/taskAPI/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestTaskHandler(t *testing.T) {
	db, err := config.ConnectDB(config.LoadConfig())
	if err != nil {
		panic(err)
	}
	repo := NewTaskRepositoryHandler(store.New(db))
	t.Run("Wrong method test", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPatch, "/task", nil)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(repo.TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
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

		assert.Equal(t, w.Code, http.StatusCreated)
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

		assert.Equal(t, w.Code, http.StatusBadRequest)
	})

	t.Run("Update task handler test", func(t *testing.T) {
		testTask := store.Task{
			ID:        0,
			Title:     "test",
			Text:      "update task",
			Completed: sql.NullBool{false, true},
			CreatedAt: sql.NullTime{time.Now(), true},
			UpdatedAt: sql.NullTime{time.Now(), true},
		}
		newTask := store.UpdateTaskParams{
			ID:        testTask.ID,
			Title:     "test success",
			Text:      "update task",
			Completed: sql.NullBool{true, true},
			UpdatedAt: sql.NullTime{time.Now(), true},
		}
		body, _ := json.Marshal(newTask)
		r := httptest.NewRequest(http.MethodPut, "/task", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(repo.TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusOK)
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

		assert.Equal(t, w.Code, http.StatusBadRequest)
	})

	//t.Run("Delete task handler test", func(t *testing.T) {
	//	testTask := store.Task{
	//		ID:        0,
	//		Title:     "test",
	//		Text:      "delete task",
	//		Completed: sql.NullBool{false, true},
	//		CreatedAt: sql.NullTime{time.Now(), true},
	//		UpdatedAt: sql.NullTime{time.Now(), true},
	//	}
	//	r := httptest.NewRequest("DELETE", "/task?id=0", nil)
	//	w := httptest.NewRecorder()
	//
	//	handler := http.HandlerFunc(repo.TaskHandler)
	//	handler.ServeHTTP(w, r)
	//
	//	_, ok := storage.Database.Tasks[testTask.ID]
	//	assert.False(t, ok)
	//})
}
