// utils/project_response.go
package utils

import (
	"project-management-backend/models"
)

type ProjectResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	WorkspaceID uint                    `json:"workspace_id"`
	CreatedBy   uint                    `json:"created_by"`
	Members     []ProjectMemberResponse `json:"members"`
	Tasks       []SimpleTaskResponse    `json:"tasks"`
	Images      []ProjectImageResponse  `json:"images"`
	MemberCount int                     `json:"member_count"`
	TaskCount   int                     `json:"task_count"`
	CreatedAt   string                  `json:"created_at"`
}

type ProjectMemberResponse struct {
	UserID        uint   `json:"user_id"`
	UserName      string `json:"user_name"`
	UserEmail     string `json:"user_email"`
	RoleInProject string `json:"role_in_project"` // string, bukan *string
}

type SimpleTaskResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ProjectImageResponse struct {
	ID         uint   `json:"id"`
	URL        string `json:"url"`         // Match dengan model (URL, bukan ImageURL)
	UploadedBy uint   `json:"uploaded_by"` // Match dengan model
}

func ToProjectResponse(project *models.Project) ProjectResponse {
	// Convert members
	var memberResponses []ProjectMemberResponse
	for _, member := range project.Members {
		memberResponses = append(memberResponses, ProjectMemberResponse{
			UserID:        member.User.ID,
			UserName:      member.User.Name,
			UserEmail:     member.User.Email,
			RoleInProject: member.RoleInProject, // Langsung string
		})
	}

	// Convert tasks
	var taskResponses []SimpleTaskResponse
	for _, task := range project.Tasks {
		taskResponses = append(taskResponses, SimpleTaskResponse{
			ID:   task.ID,
			Name: task.Title,
		})
	}

	// Convert images
	var imageResponses []ProjectImageResponse
	for _, image := range project.Images {
		imageResponses = append(imageResponses, ProjectImageResponse{
			ID:         image.ID,
			URL:        image.URL,        // Match dengan field model
			UploadedBy: image.UploadedBy, // Match dengan field model
		})
	}

	return ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		WorkspaceID: project.WorkspaceID,
		CreatedBy:   project.CreatedBy,
		Members:     memberResponses,
		Tasks:       taskResponses,
		Images:      imageResponses,
		MemberCount: len(project.Members),
		TaskCount:   len(project.Tasks),
		CreatedAt:   project.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToProjectResponseList(projects []models.Project) []ProjectResponse {
	resp := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		resp[i] = ToProjectResponse(&p)
	}
	return resp
}

func ToProjectMemberResponseList(members []models.ProjectUser) []ProjectMemberResponse {
	var memberResponses []ProjectMemberResponse
	for _, member := range members {
		memberResponses = append(memberResponses, ProjectMemberResponse{
			UserID:        member.User.ID,
			UserName:      member.User.Name,
			UserEmail:     member.User.Email,
			RoleInProject: member.RoleInProject, // Langsung string
		})
	}
	return memberResponses
}
