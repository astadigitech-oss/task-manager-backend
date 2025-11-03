// controllers/project_controller.go
package controllers

import (
	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	Service services.ProjectService
}

func NewProjectController(service services.ProjectService) *ProjectController {
	return &ProjectController{Service: service}
}

func (pc *ProjectController) ListProjects(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)
	projects, err := pc.Service.GetAllProjects(workspaceID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	respProjects := utils.ToProjectResponseList(projects)
	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Project list diambil",
		Data:    respProjects,
	})
}

func (pc *ProjectController) CreateProject(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	project := models.Project{
		Name:        input.Name,
		Description: input.Description,
		WorkspaceID: workspaceID,
	}

	if err := pc.Service.CreateProject(&project, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "CREATE_PROJECT", "project", project.ID, nil, project)

	respProject := utils.ToProjectResponse(&project)
	c.JSON(201, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Project berhasil dibuat",
		Data:    respProject,
	})
}

func (pc *ProjectController) AddMember(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role_in_project"` // string, bukan *string
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := pc.Service.AddMember(projectID, input.UserID, input.Role, currentUser); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "ADD_MEMBER_PROJECT", "project", projectID, nil, input)

	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Member berhasil ditambahkan ke project",
		Data: gin.H{
			"project_id": projectID,
			"user_id":    input.UserID,
			"role":       input.Role,
		},
	})
}

func (pc *ProjectController) GetMembers(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	members, err := pc.Service.GetMembers(projectID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	memberResponses := utils.ToProjectMemberResponseList(members)
	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Members project berhasil diambil",
		Data:    memberResponses,
	})
}

func (pc *ProjectController) DetailProject(c *gin.Context) {
	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	project, err := pc.Service.GetByID(projectID, currentUser)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	respProject := utils.ToProjectResponse(project)
	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Detail project berhasil diambil",
		Data:    respProject,
	})
}
