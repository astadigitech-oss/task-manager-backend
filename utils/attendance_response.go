package utils

import (
	"time"

	"project-management-backend/models"
)

type AttendanceImageResponse struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}

type SimpleUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SimpleWorkspaceResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AttendanceResponse struct {
	ID        uint                      `json:"id"`
	Activity  string                    `json:"activity"`
	Obstacle  *string                   `json:"obstacle"`
	ClockIn   time.Time                 `json:"clock_in"`
	ClockOut  *time.Time                `json:"clock_out,omitempty"`
	CreatedAt time.Time                 `json:"created_at"`
	User      SimpleUserResponse        `json:"user"`
	Workspace SimpleWorkspaceResponse   `json:"workspace"`
	Images    []AttendanceImageResponse `json:"images,omitempty"`
}

func ToAttendanceResponse(attendance models.Attendance) AttendanceResponse {
	var images []AttendanceImageResponse
	for _, img := range attendance.Images {
		images = append(images, AttendanceImageResponse{
			ID:  img.ID,
			URL: img.URL,
		})
	}

	return AttendanceResponse{
		ID:        attendance.ID,
		Activity:  attendance.Activity,
		Obstacle:  attendance.Obstacle,
		ClockIn:   attendance.ClockIn,
		ClockOut:  attendance.ClockOut,
		CreatedAt: attendance.CreatedAt,
		User: SimpleUserResponse{
			ID:    attendance.User.ID,
			Name:  attendance.User.Name,
			Email: attendance.User.Email,
		},
		Workspace: SimpleWorkspaceResponse{
			ID:   attendance.Workspace.ID,
			Name: attendance.Workspace.Name,
		},
		Images: images,
	}
}
