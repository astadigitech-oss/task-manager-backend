package models

import "time"

// DailyActivityItem represents a single row in the daily report.
// It captures a snapshot of a task's state at the time of an activity.
type DailyActivityItem struct {
	ActivityTime time.Time
	User         string
	ProjectTitle string
	TaskTitle    string
	TaskPriority string
	StatusAtLog  string // The status of the task as recorded in the activity log.
	Overdue      time.Duration
}
