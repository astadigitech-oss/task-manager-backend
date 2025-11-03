package utils

import "project-management-backend/models"

type SimpleMemberResponse struct {
	UserID          uint    `json:"user_id"`
	UserName        string  `json:"user_name"`
	UserEmail       string  `json:"user_email"`
	RoleInWorkspace *string `json:"role_in_workspace"`
}

func ToMemberResponseList(members []models.WorkspaceUser) []SimpleMemberResponse {
	var simpleResponses []SimpleMemberResponse
	for _, member := range members {
		simpleResponses = append(simpleResponses, SimpleMemberResponse{
			UserID:          member.User.ID,
			UserName:        member.User.Name,
			UserEmail:       member.User.Email,
			RoleInWorkspace: member.RoleInWorkspace,
		})
	}
	return simpleResponses
}
