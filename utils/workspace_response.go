// utils/workspace_response.go
package utils

import "project-management-backend/models"

type WorkspaceResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	CreatedBy   uint                    `json:"created_by,omitempty"`
	Projects    []SimpleProjectResponse `json:"projects"`
	MemberCount int                     `json:"member_count,omitempty"`
}

type SimpleProjectResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToWorkspaceResponse(workspace *models.Workspace) WorkspaceResponse {
	// Convert projects to simple list
	var projectResponses []SimpleProjectResponse
	for _, project := range workspace.Projects {
		projectResponses = append(projectResponses, SimpleProjectResponse{
			ID:   project.ID,
			Name: project.Name,
		})
	}

	return WorkspaceResponse{
		ID:          workspace.ID,
		Name:        workspace.Name,
		Description: workspace.Description,
		CreatedBy:   workspace.CreatedBy,
		Projects:    projectResponses,
		MemberCount: len(workspace.Members),
	}
}

func ToWorkspaceResponseList(workspaces []models.Workspace) []WorkspaceResponse {
	resp := make([]WorkspaceResponse, len(workspaces))
	for i, w := range workspaces {
		resp[i] = ToWorkspaceResponse(&w)
	}
	return resp
}
