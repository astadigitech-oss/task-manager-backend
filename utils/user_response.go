package utils

import (
	"time"
)

// SimpleWorkspace adalah representasi ramping dari sebuah workspace.
type SimpleWorkspace struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// SimpleProject adalah representasi ramping dari sebuah project.
type SimpleProject struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// SimpleTask adalah representasi ramping dari sebuah task.
type SimpleTask struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

// OnlineUserResponse mendefinisikan struktur untuk endpoint /api/online-users.
type OnlineUserResponse struct {
	ID         uint              `json:"id"`
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	Role       string            `json:"role"`
	IsOnline   bool              `json:"is_online"`
	LastSeen   *time.Time        `json:"last_seen,omitempty"`
	Workspaces []SimpleWorkspace `json:"workspaces"`
	Projects   []SimpleProject   `json:"projects"`
	Tasks      []SimpleTask      `json:"tasks"`
}
