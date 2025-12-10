package storage

import (
	"sync"
	"time"
)

// Task represents a task entity
// @Description Task information
type Task struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    bool      `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// Storage holds the in-memory database
type Storage struct {
	Mu     sync.RWMutex
	Tasks  map[int]*Task
	NextID int
}

var Database = Storage{
	Tasks: make(map[int]*Task),
}
