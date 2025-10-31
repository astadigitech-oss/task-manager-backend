package controllers

import (
	"fmt"
	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkspaceController struct {
	Service services.WorkspaceService
}

func NewWorkspaceController(service services.WorkspaceService) *WorkspaceController {
	return &WorkspaceController{Service: service}
}

func GetCurrentUser(c *gin.Context) *models.User {
	return &models.User{ID: 1, Name: "Admin", Role: "admin"} //dummy aja ini
}

func ParseUintParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, fmt.Errorf("%s is required", paramName)
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %s", paramName, idStr)
	}

	return uint(id), nil
}

func (wc *WorkspaceController) ListWorkspaces(c *gin.Context) {
	currentUser := GetCurrentUser(c)
	workspaces, err := wc.Service.GetAllWorkspaces(currentUser)
	if err != nil {
		utils.Error(0, "GET_WORKSPACE", "workspaces", 403, err.Error(), "Failed to get workspaces")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}
	respWorkspaces := utils.ToWorkspaceResponseList(workspaces)
	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Workspace list diambil",
		Data:    respWorkspaces,
	})
}

func (wc *WorkspaceController) CreateWorkspace(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	workspace := models.Workspace{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := wc.Service.CreateWorkspace(&workspace, currentUser); err != nil {
		utils.Error(currentUser.ID, "CREATE_WORKSPACE", "workspaces", 403, err.Error(), "Failed to create workspaces")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}
	utils.ActivityLog(currentUser.ID, "CREATE_WORKSPACE", "workspace", currentUser.ID, nil, workspace)

	respWorkspaces := utils.ToWorkspaceResponse(&workspace)
	c.JSON(201, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Workspace berhasil di buat",
		Data:    respWorkspaces,
	})
}

func (wc *WorkspaceController) AddMember(c *gin.Context) {
	var input struct {
		UserID uint    `json:"user_id"`
		Role   *string `json:"role_in_workspace"`
	}
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := wc.Service.AddMember(workspaceID, input.UserID, input.Role); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Berhasil tambah member"})
}

func (wc *WorkspaceController) GetMembers(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	members, err := wc.Service.GetMembers(workspaceID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": members})
}

func (wc *WorkspaceController) DetailWorkspace(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ws, err := wc.Service.GetByID(workspaceID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": ws})
}

// // POST /api/workspaces
// func (wc *WorkspaceController) CreateWorkspace(c *gin.Context) {
// 	userIDVal, _ := c.Get("user_id")
// 	roleVal, _ := c.Get("role")
// 	userID := userIDVal.(uint)
// 	role := roleVal.(string)

// 	var input struct {
// 		Name string `json:"name"`
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
// 		return
// 	}

// 	workspace, err := wc.Service.CreateWorkspace(userID, role, input.Name)
// 	if err != nil {
// 		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, workspace)
// }
