package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zordddd/learning/taskAPI/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestTaskHandler(t *testing.T) {
	t.Run("Wrong method test", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPatch, "/task", nil)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusMethodNotAllowed)
	})

	t.Run("Get task handler test", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/task", nil)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create task handler test", func(t *testing.T) {
		body, _ := json.Marshal(storage.Task{
			Name:   "test",
			Status: false,
		})
		r := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(TaskHandler)
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

		handler := http.HandlerFunc(TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusBadRequest)
	})

	t.Run("Update task handler test", func(t *testing.T) {
		testTask := storage.Task{
			ID:     0,
			Name:   "test",
			Status: false,
		}
		newTask := storage.Task{
			ID:     testTask.ID,
			Name:   testTask.Name,
			Status: true,
		}
		body, _ := json.Marshal(newTask)
		storage.Database.Tasks[testTask.ID] = &testTask
		r := httptest.NewRequest(http.MethodPut, "/task", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusOK)
		assert.Equal(t, storage.Database.Tasks[testTask.ID].Status, true)
	})

	t.Run("Update wrong task handler test", func(t *testing.T) {
		body, _ := json.Marshal(struct {
			Name int `json:"name"`
		}{
			Name: 0,
		})
		r := httptest.NewRequest(http.MethodPut, "/task", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(TaskHandler)
		handler.ServeHTTP(w, r)

		assert.Equal(t, w.Code, http.StatusBadRequest)
	})

	t.Run("Delete task handler test", func(t *testing.T) {
		testTask := storage.Task{
			ID:     0,
			Name:   "test",
			Status: false,
		}
		storage.Database.Tasks[testTask.ID] = &testTask
		r := httptest.NewRequest("DELETE", "/task?id=0", nil)
		w := httptest.NewRecorder()

		handler := http.HandlerFunc(TaskHandler)
		handler.ServeHTTP(w, r)

		_, ok := storage.Database.Tasks[testTask.ID]
		assert.False(t, ok)
	})
}
