package models

// TaskWithHistory is a struct that holds a task and its status logs
type TaskWithHistory struct {
	Task       Task            `json:"task"`
	StatusLogs []TaskStatusLog `json:"status_logs"`
}
