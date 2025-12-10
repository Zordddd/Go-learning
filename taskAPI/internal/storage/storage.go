package storage

import (
	"sync"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    bool      `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Storage struct {
	Mu     sync.RWMutex
	Tasks  map[int]Task
	NextID int
}

var Database = Storage{
	Tasks: make(map[int]Task),
}
