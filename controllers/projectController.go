// controllers/project_controller.go
package controllers

import (
	"fmt"
	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	Service services.ProjectService
}

func NewProjectController(service services.ProjectService) *ProjectController {
	return &ProjectController{Service: service}
}

func (pc *ProjectController) ListProjects(c *gin.Context) {
	currentUser := GetCurrentUser(c)

	projects, err := pc.Service.GetAllProjects(currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "list_projects", "project", 0, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	var projectList []gin.H
	for _, project := range projects {
		members := make([]gin.H, 0)
		for _, member := range project.Members {
			members = append(members, gin.H{
				"id":            member.User.ID,
				"name":          member.User.Name,
				"profile_image": member.User.ProfileImage,
			})
		}

		var completedTasks int
		for _, task := range project.Tasks {
			if task.Status == "done" {
				completedTasks++
			}
		}

		var progress float64
		if len(project.Tasks) > 0 {
			progress = (float64(completedTasks) / float64(len(project.Tasks))) * 100
		} else {
			progress = 0
		}
		projectList = append(projectList, gin.H{
			"id":           project.ID,
			"name":         project.Name,
			"description":  project.Description,
			"workspace_id": project.WorkspaceID,
			"member_count": len(project.Members),
			"task_count":   len(project.Tasks),
			"members":      members,
			"progress":     progress,
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "List project berhasil di ambil",
		Data:    projectList,
	})
}

func (pc *ProjectController) CreateProject(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		WorkspaceID uint   `json:"workspace_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "bind_json", "project", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	project := models.Project{
		Name:        input.Name,
		Description: input.Description,
		WorkspaceID: input.WorkspaceID,
	}

	if err := pc.Service.CreateProject(&project, currentUser); err != nil {
		utils.Error(currentUser.ID, "create_project", "project", 0, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "CREATE_PROJECT", "project", project.ID, nil, project)

	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: "Project berhasil di buat",
		Data: gin.H{
			"id":           project.ID,
			"name":         project.Name,
			"description":  project.Description,
			"workspace_id": project.WorkspaceID,
		},
	})
}

func (pc *ProjectController) DetailProject(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "project", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	project, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_project_by_id", "project", projectID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	members := make([]gin.H, 0)
	for _, member := range project.Members {
		members = append(members, gin.H{
			"id":            member.User.ID,
			"name":          member.User.Name,
			"profile_image": member.User.ProfileImage,
			"position":      member.User.Position,
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Detail project berhasil diambil",
		Data: gin.H{
			"id":           project.ID,
			"name":         project.Name,
			"description":  project.Description,
			"workspace_id": project.WorkspaceID,
			"created_by":   project.CreatedBy,
			"member":       members,
		},
	})
}

func (pc *ProjectController) UpdateProject(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "project", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "bind_json", "project", projectID, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	// Get data lama sebelum update
	oldProject, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_project_before_update", "project", projectID, err.Error(), "")
		c.JSON(403, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	project := models.Project{
		ID:          projectID,
		Name:        input.Name,
		Description: input.Description,
	}

	if err := pc.Service.UpdateProject(&project, currentUser); err != nil {
		utils.Error(currentUser.ID, "update_project", "project", projectID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	// Get data baru setelah update
	updatedProject, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_project_after_update", "project", projectID, err.Error(), "")
		c.JSON(403, gin.H{"error": "Gagal mengambil data project setelah update"})
		return
	}

	utils.ActivityLog(currentUser.ID, "UPDATE_PROJECT", "project", projectID, oldProject, updatedProject)

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Project berhasil di update",
		Data: gin.H{
			"id":           updatedProject.ID,
			"name":         updatedProject.Name,
			"description":  updatedProject.Description,
			"workspace_id": updatedProject.WorkspaceID,
		},
	})
}

func (pc *ProjectController) SoftDeleteProject(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "project", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	oldProject, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_project_before_soft_delete", "project", projectID, err.Error(), "")
		c.JSON(403, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	if err := pc.Service.SoftDeleteProject(projectID, currentUser); err != nil {
		utils.Error(currentUser.ID, "soft_delete_project", "project", projectID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "SOFT_DELETE_PROJECT", "project", projectID, oldProject, "deleted")
	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Project berhasil di soft delete",
		Data:    gin.H{"project_id": projectID},
	})
}

func (pc *ProjectController) DeleteProject(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "project", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	project, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_project_before_hard_delete", "project", projectID, err.Error(), "")
		c.JSON(404, gin.H{"error": "Project tidak ditemukan"})
		return
	}

	var input struct {
		Confirm bool `json:"confirm"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || !input.Confirm {
		utils.Error(currentUser.ID, "confirm_hard_delete", "project", projectID, "Confirmation required for hard delete", "")
		c.JSON(400, gin.H{
			"error":   "Konfirmasi diperlukan untuk hard delete",
			"warning": "Tindakan ini akan menghapus PERMANEN:",
			"data_akan_dihapus": gin.H{
				"project":       project.Name,
				"tasks_count":   len(project.Tasks),
				"members_count": len(project.Members),
				"images_count":  len(project.Images),
			},
			"confirmation_required": true,
		})
		return
	}

	if err := pc.Service.DeleteProject(projectID, currentUser); err != nil {
		utils.Error(currentUser.ID, "hard_delete_project", "project", projectID, err.Error(), "")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "DELETE_PROJECT", "project", projectID, project, "hard deleted")

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Project dan semua data terkait berhasil dihapus permanen",
		Data: gin.H{
			"project_id":      projectID,
			"project_name":    project.Name,
			"deleted_tasks":   len(project.Tasks),
			"deleted_members": len(project.Members),
			"deleted_images":  len(project.Images),
		},
	})
}

func (pc *ProjectController) AddMember(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "project", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Members []struct {
			UserID uint   `json:"user_id" binding:"required"`
			Role   string `json:"role_in_project"`
		} `json:"members" binding:"required,min=1,dive"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "bind_json_add_member", "project", projectID, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	projectMembers := make([]services.ProjectMember, len(input.Members))
	for i, member := range input.Members {
		projectMembers[i] = services.ProjectMember{
			UserID: member.UserID,
			Role:   member.Role,
		}
	}

	currentUser := GetCurrentUser(c)

	if err := pc.Service.AddMembers(projectID, projectMembers, currentUser); err != nil {
		utils.Error(currentUser.ID, "add_members", "project", projectID, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//Log Activity
	addedMembers := make([]gin.H, len(input.Members))
	for i, member := range input.Members {
		roleStr := "member"
		if member.Role != "" {
			roleStr = member.Role
		}
		addedMembers[i] = gin.H{
			"user_id": member.UserID,
			"role":    roleStr,
		}
	}
	utils.ActivityLog(currentUser.ID, "ADD_MEMBERS_PROJECT", "project", projectID, nil, gin.H{
		"members_added": len(input.Members),
		"user_ids":      addedMembers,
	})

	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: fmt.Sprintf("%d member berhasil ditambahkan ke project", len(input.Members)),
		Data: gin.H{
			"project_id":  projectID,
			"members":     addedMembers,
			"total_added": len(input.Members),
		},
	})
}

func (pc *ProjectController) GetMembers(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "project", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	members, err := pc.Service.GetMembers(projectID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_members", "project", projectID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	memberProject := make([]gin.H, 0)
	for _, member := range members {
		memberProject = append(memberProject, gin.H{
			"id":          member.User.ID,
			"name":        member.User.Name,
			"user_email":  member.User.Email,
			"role":        member.User.Role,
			"profile_img": member.User.ProfileImage,
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Members project berhasil diambil",
		Data:    memberProject,
	})
}

func (pc *ProjectController) RemoveMember(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "project", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		UserIDs []uint `json:"user_ids" binding:"required,min=1,dive"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "bind_json_remove_member", "project", projectID, err.Error(), "")
		c.JSON(400, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if len(input.UserIDs) == 1 {
		if err := pc.Service.RemoveMember(projectID, input.UserIDs[0], currentUser); err != nil {
			utils.Error(currentUser.ID, "remove_single_member", "project", projectID, err.Error(), "")
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := pc.Service.RemoveMembers(projectID, input.UserIDs, currentUser); err != nil {
			utils.Error(currentUser.ID, "remove_multiple_members", "project", projectID, err.Error(), "")
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	utils.ActivityLog(currentUser.ID, "REMOVE_MEMBERS_PROJECT", "project", projectID, nil, gin.H{
		"user_ids_removed": input.UserIDs,
		"total_removed":    len(input.UserIDs),
	})

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: fmt.Sprintf("%d member berhasil dihapus dari project", len(input.UserIDs)),
		Data: gin.H{
			"project_id":       projectID,
			"user_ids_removed": input.UserIDs,
			"total_removed":    len(input.UserIDs),
		},
	})
}

func (pc *ProjectController) RemoveSingleMember(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "project", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	memberIDStr := c.Param("user_id")
	if memberIDStr == "" {
		utils.Error(0, "get_user_id_param", "project", projectID, "user_id is required", "")
		c.JSON(400, gin.H{"error": "user_id is required"})
		return
	}

	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		utils.Error(0, "parse_user_id", "project", projectID, "invalid user_id", "")
		c.JSON(400, gin.H{"error": "invalid user_id"})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := pc.Service.RemoveMember(projectID, uint(memberID), currentUser); err != nil {
		utils.Error(currentUser.ID, "remove_single_member_param", "project", projectID, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "REMOVE_MEMBER_PROJECT", "project", projectID, nil, gin.H{
		"member_id_removed": memberID,
	})

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Member berhasil dihapus dari project",
		Data: gin.H{
			"project_id":        projectID,
			"member_id_removed": memberID,
		},
	})
}
