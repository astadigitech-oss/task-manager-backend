package utils

import "project-management-backend/models"

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Mapper single
func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

// Mapper list
func ToUserResponseList(users []models.User) []UserResponse {
	resp := make([]UserResponse, len(users))
	for i, u := range users {
		resp[i] = ToUserResponse(&u)
	}
	return resp
}
