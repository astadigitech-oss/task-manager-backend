// utils/task_response.go
package utils

import (
	"project-management-backend/models"
	"time"
)

type TaskResponse struct {
	ID              uint                 `json:"id"`
	Title           string               `json:"title"`
	Description     string               `json:"description"`
	Status          string               `json:"status"`
	Priority        string               `json:"priority"`
	StartDate       time.Time            `json:"start_date"`
	DueDate         time.Time            `json:"due_date"`
	Notes           *string              `json:"notes"`
	ProjectID       uint                 `json:"project_id"`
	Members         []TaskMemberResponse `json:"members"`
	Images          []TaskImageResponse  `json:"images"`
	MemberCount     int                  `json:"member_count"`
	OverDueDuration int                  `json:"overdue_duration"`
	CreatedAt       string               `json:"created_at"`
}

type SimpleTaskResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type TaskMemberResponse struct {
	UserID           uint   `json:"user_id"`
	UserName         string `json:"user_name"`
	UserEmail        string `json:"user_email"`
	UserProfileImage string `json:"user_profile_image"`
	RoleInTask       string `json:"role_in_task"`
	AssignedAt       string `json:"assigned_at"`
}

type TaskImageResponse struct {
	ID  uint   `json:"id"`
	URL string `json:"url"`
}

func ToTaskResponse(task *models.Task) TaskResponse {
	// Convert members
	var memberResponses []TaskMemberResponse
	for _, member := range task.Members {
		profileImage := ""
		if member.User.ProfileImage != nil {
			profileImage = *member.User.ProfileImage
		}
		memberResponses = append(memberResponses, TaskMemberResponse{
			UserID:           member.User.ID,
			UserName:         member.User.Name,
			UserEmail:        member.User.Email,
			UserProfileImage: profileImage,
			RoleInTask:       member.RoleInTask,
			AssignedAt:       member.AssignedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Convert images
	var imageResponses []TaskImageResponse
	for _, image := range task.Images {
		imageResponses = append(imageResponses, TaskImageResponse{
			ID:  image.ID,
			URL: image.URL,
		})
	}

	return TaskResponse{
		ID:              task.ID,
		Title:           task.Title,
		Description:     task.Description,
		Status:          task.Status,
		Priority:        task.Priority,
		StartDate:       task.StartDate,
		DueDate:         task.DueDate,
		Notes:           task.Notes,
		ProjectID:       task.ProjectID,
		Members:         memberResponses,
		Images:          imageResponses,
		MemberCount:     len(task.Members),
		OverDueDuration: int(task.OverdueDuration),
		CreatedAt:       task.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToTaskResponseList(tasks []models.Task) []TaskResponse {
	resp := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = ToTaskResponse(&t)
	}
	return resp
}
