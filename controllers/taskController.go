// controllers/task_controller.go
package controllers

import (
	"project-management-backend/models"
	"project-management-backend/services"
	"project-management-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	Service services.TaskService
}

func NewTaskController(service services.TaskService) *TaskController {
	return &TaskController{Service: service}
}

func (tc *TaskController) ListTasks(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		utils.Error(0, "parse_workspace_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)
	tasks, err := tc.Service.GetAllTasks(projectID, workspaceID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "list_tasks", "task", 0, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	respTasks := utils.ToTaskResponseList(tasks)
	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "List task berhasil di ambil",
		Data:    respTasks,
	})
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		utils.Error(0, "parse_workspace_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		Priority    string    `json:"priority"`
		StartDate   time.Time `json:"start_date"`
		DueDate     time.Time `json:"due_date"`
		Notes       *string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "bind_json", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	task := models.Task{
		ProjectID:   projectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		StartDate:   input.StartDate,
		DueDate:     input.DueDate,
		Notes:       input.Notes,
	}

	if err := tc.Service.CreateTask(&task, workspaceID, currentUser); err != nil {
		utils.Error(currentUser.ID, "create_task", "task", 0, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "CREATE_TASK", "task", task.ID, nil, task)

	respTask := utils.ToTaskResponse(&task)
	c.JSON(201, APIResponse{
		Success: true,
		Code:    201,
		Message: "Task berhasil dibuat",
		Data:    respTask,
	})
}

func (tc *TaskController) DetailTask(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		utils.Error(0, "parse_workspace_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		utils.Error(0, "parse_task_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	task, err := tc.Service.GetByID(taskID, workspaceID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_task_by_id", "task", taskID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	respTask := utils.ToTaskResponse(task)
	c.JSON(200, utils.APIResponse{
		Success: true,
		Code:    200,
		Message: "Detail task berhasil diambil",
		Data:    respTask,
	})
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		utils.Error(0, "parse_workspace_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		utils.Error(0, "parse_task_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.Error(0, "bind_json", "task", taskID, err.Error(), "")
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	currentUser := GetCurrentUser(c)

	oldTask, err := tc.Service.GetByID(taskID, workspaceID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_task_before_update", "task", taskID, err.Error(), "")
		c.JSON(403, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	if err := tc.Service.UpdateTask(taskID, updates, workspaceID, currentUser); err != nil {
		utils.Error(currentUser.ID, "update_task", "task", taskID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	// Get data baru setelah update
	updatedTask, err := tc.Service.GetByID(taskID, workspaceID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_task_after_update", "task", taskID, err.Error(), "")
		c.JSON(403, gin.H{"error": "Gagal mengambil data task setelah update"})
		return
	}

	utils.ActivityLog(currentUser.ID, "UPDATE_TASK", "task", taskID, oldTask, updatedTask)

	respTask := utils.ToTaskResponse(updatedTask)
	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Task berhasil diupdate",
		Data:    respTask,
	})
}

// SoftDeleteTask - Soft delete task
func (tc *TaskController) SoftDeleteTask(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		utils.Error(0, "parse_workspace_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		utils.Error(0, "parse_task_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := tc.Service.SoftDeleteTask(taskID, workspaceID, currentUser); err != nil {
		utils.Error(currentUser.ID, "soft_delete_task", "task", taskID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Task berhasil di soft delete",
		Data:    gin.H{"task_id": taskID},
	})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		utils.Error(0, "parse_workspace_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		utils.Error(0, "parse_task_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	// Get task info untuk show warning
	task, err := tc.Service.GetByID(taskID, workspaceID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_task_before_hard_delete", "task", taskID, err.Error(), "")
		c.JSON(404, gin.H{"error": "Task tidak ditemukan"})
		return
	}

	var input struct {
		Confirm bool `json:"confirm"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || !input.Confirm {
		utils.Error(currentUser.ID, "confirm_hard_delete", "task", taskID, "Confirmation required for hard delete", "")
		c.JSON(400, gin.H{
			"error":   "Konfirmasi diperlukan untuk hard delete",
			"warning": "Tindakan ini akan menghapus PERMANEN:",
			"data_akan_dihapus": gin.H{
				"task":          task.Title,
				"members_count": len(task.Members),
				"images_count":  len(task.Images),
			},
			"confirmation_required": true,
		})
		return
	}

	if err := tc.Service.DeleteTask(taskID, workspaceID, currentUser); err != nil {
		utils.Error(currentUser.ID, "hard_delete_task", "task", taskID, err.Error(), "")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Task dan semua data terkait berhasil dihapus permanen",
		Data: gin.H{
			"task_id":         taskID,
			"task_title":      task.Title,
			"deleted_members": len(task.Members),
			"deleted_images":  len(task.Images),
		},
	})
}

func (tc *TaskController) AddMember(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		utils.Error(0, "parse_workspace_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		utils.Error(0, "parse_task_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ProjectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role_in_task"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(0, "bind_json_add_member", "task", taskID, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	if err := tc.Service.AddMember(taskID, ProjectID, workspaceID, input.UserID, input.Role, currentUser); err != nil {
		utils.Error(currentUser.ID, "add_member_task", "task", taskID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "ADD_MEMBER_TASK", "task", taskID, nil, input)

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Member berhasil ditambahkan ke task",
		Data: gin.H{
			"task_id": taskID,
			"user_id": input.UserID,
			"role":    input.Role,
		},
	})
}

func (tc *TaskController) GetMembers(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		utils.Error(0, "parse_workspace_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		utils.Error(0, "parse_task_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	members, err := tc.Service.GetMembers(taskID, projectID, workspaceID, currentUser)
	if err != nil {
		utils.Error(currentUser.ID, "get_members_task", "task", taskID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	memberTask := make([]gin.H, 0)
	for _, member := range members {
		memberTask = append(memberTask, gin.H{
			"id":           member.User.ID,
			"name":         member.User.Name,
			"user_email":   member.User.Email,
			"role_in_task": member.RoleInTask,
			"profile_img":  member.User.ProfileImage,
			"workspace_id": workspaceID,
		})
	}

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Members task berhasil diambil",
		Data:    memberTask,
	})
}

func (tc *TaskController) DeleteMember(c *gin.Context) {
	workspaceID, err := ParseUintParam(c, "workspace_id")
	if err != nil {
		utils.Error(0, "parse_workspace_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	projectID, err := ParseUintParam(c, "project_id")
	if err != nil {
		utils.Error(0, "parse_project_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskID, err := ParseUintParam(c, "task_id")
	if err != nil {
		utils.Error(0, "parse_task_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	currentUser := GetCurrentUser(c)

	userID, err := ParseUintParam(c, "user_id")
	if err != nil {
		utils.Error(0, "parse_user_id", "task", 0, err.Error(), "")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := tc.Service.DeleteMember(taskID, projectID, workspaceID, userID, currentUser); err != nil {
		utils.Error(currentUser.ID, "delete_member_task", "task", taskID, err.Error(), "")
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	utils.ActivityLog(currentUser.ID, "DELETE_MEMBER_TASK", "task", taskID, gin.H{"deleted_user_id": userID}, nil)

	c.JSON(200, APIResponse{
		Success: true,
		Code:    200,
		Message: "Member berhasil dihapus dari task",
		Data: gin.H{
			"task_id": taskID,
			"user_id": userID,
		},
	})
}
