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
	user, _ := c.Get("currentUser")
	return user.(*models.User)
}

type APIResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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
	workspaceList := make([]gin.H, 0)
	for _, ws := range workspaces {
		workspaceList = append(workspaceList, gin.H{
			"id":    ws.ID,
			"name":  ws.Name,
			"color": ws.Color,
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "List workspace berhasil di ambil",
		Data:    workspaceList,
	})
}

func (wc *WorkspaceController) CreateWorkspace(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	workspace := models.Workspace{
		Name:        input.Name,
		Description: input.Description,
		Color:       input.Color,
	}

	if err := wc.Service.CreateWorkspace(&workspace, currentUser); err != nil {
		utils.Error(currentUser.ID, "CREATE_WORKSPACE", "workspaces", 403, err.Error(), "Failed to create workspaces")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}
	utils.ActivityLog(currentUser.ID, "CREATE_WORKSPACE", "workspace", currentUser.ID, nil, workspace)

	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: "Workspace berhasil di buat",
		Data: gin.H{
			"id":    workspace.ID,
			"name":  workspace.Name,
			"color": workspace.Color,
		},
	})
}

func (wc *WorkspaceController) DetailWorkspace(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	ws, err := wc.Service.GetByID(workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Detail workspace berhasil diambil",
		Data: gin.H{
			"id":          ws.ID,
			"name":        ws.Name,
			"description": ws.Description,
			"createdBy":   ws.CreatedBy,
			"color":       ws.Color,
		},
	})
}

func (wc *WorkspaceController) UpdateWorkspace(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	oldWorkspace, err := wc.Service.GetByID(workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": "Workspace tidak ditemukan"})
		return
	}

	workspace := models.Workspace{
		ID:          workspaceID,
		Name:        input.Name,
		Description: input.Description,
		Color:       input.Color,
	}

	if err := wc.Service.UpdateWorkspace(&workspace, currentUser); err != nil {
		utils.Error(currentUser.ID, "UPDATE_WORKSPACE", "workspaces", 403, err.Error(), "Failed to update workspace")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	// Get updated workspace
	updatedWorkspace, err := wc.Service.GetByID(workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": "Gagal mengambil data workspace setelah update"})
		return
	}

	utils.ActivityLog(currentUser.ID, "UPDATE_WORKSPACE", "workspace", workspaceID, oldWorkspace, workspace)

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Workspace berhasil diupdate",
		Data: gin.H{
			"id":          updatedWorkspace.ID,
			"name":        updatedWorkspace.Name,
			"description": updatedWorkspace.Description,
			"color":       updatedWorkspace.Color,
		},
	})
}

func (wc *WorkspaceController) SoftDeleteWorkspace(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := wc.Service.SoftDeleteWorkspace(workspaceID, currentUser); err != nil {
		utils.Error(currentUser.ID, "SOFT_DELETE_WORKSPACE", "workspaces", 403, err.Error(), "Failed to soft delete workspace")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "SOFT_DELETE_WORKSPACE", "workspace", workspaceID, nil, nil)

	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Workspace berhasil di soft delete",
		Data:    gin.H{"workspace_id": workspaceID},
	})
}

func (wc *WorkspaceController) DeleteWorkspace(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	workspace, err := wc.Service.GetByID(workspaceID, currentUser)
	if err != nil {
		c.JSON(404, gin.H{"error": "Workspace tidak ditemukan"})
		return
	}
	var input struct {
		Confirm bool `json:"confirm"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || !input.Confirm {
		c.JSON(400, gin.H{
			"error":   "Konfirmasi diperlukan untuk hard delete workspace",
			"warning": "Tindakan ini akan menghapus PERMANEN semua data:",
			"data_akan_dihapus": gin.H{
				"workspace":      workspace.Name,
				"projects_count": len(workspace.Projects),
				"members_count":  len(workspace.Members),
				// "total_tasks": wc.countTotalTasks(workspace.Projects),
			},
			"confirmation_required": true,
		})
		return
	}

	if err := wc.Service.DeleteWorkspace(workspaceID, currentUser); err != nil {
		utils.Error(currentUser.ID, "DELETE_WORKSPACE", "workspaces", 403, err.Error(), "Failed to delete workspace")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "DELETE_WORKSPACE", "workspace", workspaceID, nil, nil)

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Workspace berhasil dihapus permanen",
		Data: gin.H{
			"workspace_id":     workspaceID,
			"workspace_name":   workspace.Name,
			"deleted_projects": len(workspace.Projects),
			"deleted_members":  len(workspace.Members),
		},
	})
}

func (wc *WorkspaceController) AddMembers(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Members []struct {
			UserID uint    `json:"user_id" binding:"required"`
			Role   *string `json:"role_in_workspace"`
		} `json:"members" binding:"required,min=1,dive"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	workspaceMembers := make([]services.WorkspaceMember, len(input.Members))
	for i, member := range input.Members {
		workspaceMembers[i] = services.WorkspaceMember{
			UserID: member.UserID,
			Role:   member.Role,
		}
	}

	currentUser := GetCurrentUser(c)

	if err := wc.Service.AddMembers(workspaceID, workspaceMembers, currentUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	addedMembers := make([]gin.H, len(input.Members))
	for i, member := range input.Members {
		roleStr := "member"
		if member.Role != nil {
			roleStr = *member.Role
		}
		addedMembers[i] = gin.H{
			"user_id": member.UserID,
			"role":    roleStr,
		}
	}

	utils.ActivityLog(currentUser.ID, "ADD_MEMBERS", "workspace", workspaceID, nil, gin.H{
		"members_added": len(input.Members),
		"user_ids":      input.Members,
	})

	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: fmt.Sprintf("%d member berhasil ditambahkan ke workspace", len(input.Members)),
		Data: gin.H{
			"workspace_id": workspaceID,
			"members":      addedMembers,
			"total_added":  len(input.Members),
		},
	})
}

func (wc *WorkspaceController) GetMembers(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	members, err := wc.Service.GetMembers(workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	memberList := make([]gin.H, 0)
	for _, member := range members {
		memberList = append(memberList, gin.H{
			"id":          member.User.ID,
			"name":        member.User.Name,
			"role":        member.User.Role,
			"profile_img": member.User.ProfileImage,
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "List member berhasil di ambil",
		Data:    memberList,
	})
}

func (wc *WorkspaceController) RemoveMember(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		UserIDs []uint `json:"user_ids" binding:"required,min=1,dive"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if len(input.UserIDs) == 1 {
		if err := wc.Service.RemoveMember(workspaceID, input.UserIDs[0], currentUser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := wc.Service.RemoveMembers(workspaceID, input.UserIDs, currentUser); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	utils.ActivityLog(currentUser.ID, "REMOVE_MEMBERS", "workspace", workspaceID, nil, gin.H{
		"user_ids_removed": input.UserIDs,
		"total_removed":    len(input.UserIDs),
	})

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: fmt.Sprintf("%d member berhasil dihapus dari workspace", len(input.UserIDs)),
		Data: gin.H{
			"workspace_id":     workspaceID,
			"user_ids_removed": input.UserIDs,
			"total_removed":    len(input.UserIDs),
		},
	})
}

func (wc *WorkspaceController) RemoveSingleMember(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	memberID, err := ParseUintParam(c, "user_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := wc.Service.RemoveMember(workspaceID, memberID, currentUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "REMOVE_MEMBER", "workspace", workspaceID, nil, gin.H{
		"member_id_removed": memberID,
	})

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Member berhasil dihapus dari workspace",
		Data: gin.H{
			"workspace_id":      workspaceID,
			"member_id_removed": memberID,
		},
	})
}
