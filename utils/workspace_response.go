package utils

import "project-management-backend/models"

type WorkspaceResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Mapper single
func ToWorkspaceResponse(user *models.Workspace) WorkspaceResponse {
	return WorkspaceResponse{
		ID:          user.ID,
		Name:        user.Name,
		Description: user.Description,
	}
}

// Mapper list
func ToWorkspaceResponseList(work []models.Workspace) []WorkspaceResponse {
	resp := make([]WorkspaceResponse, len(work))
	for i, u := range work {
		resp[i] = ToWorkspaceResponse(&u)
	}
	return resp
}
