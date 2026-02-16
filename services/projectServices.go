// services/project_service.go
package services

import (
	"bytes"
	"errors"
	"fmt"
	"project-management-backend/config"
	"project-management-backend/models"
	"project-management-backend/repositories"
	"project-management-backend/utils"
	"regexp"
	"sort"
	"strings"

	"time"

	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type ProjectService interface {
	CreateProject(project *models.Project, user *models.User) error
	GetAllProjects(user *models.User) ([]models.Project, error)
	UpdateProject(project *models.Project, user *models.User) error
	SoftDeleteProject(projectID uint, user *models.User) error
	DeleteProject(projectID uint, user *models.User) error
	GetByID(projectID uint, user *models.User) (*models.Project, error)
	AddMembers(projectID uint, members []ProjectMember, currentUser *models.User) error
	GetMembers(projectID uint, user *models.User) ([]models.ProjectUser, error)
	RemoveMember(projectID uint, userID uint, currentUser *models.User) error
	RemoveMembers(projectID uint, userIDs []uint, currentUser *models.User) error
	ExportWeeklyBackward(projectID uint, userID uint) ([]byte, error)
	ExportWeeklyForward(projectID uint, userID uint) ([]byte, error)
	ExportDaily(projectID uint, userID uint) ([]byte, error)
	ExportMonitoring(projectID uint, userID uint) ([]byte, error)
}

type ProjectMember struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role_in_project"`
}

type projectService struct {
	repo              repositories.ProjectRepository
	userRepo          repositories.UserRepository
	workspaceRepo     repositories.WorkspaceRepository
	taskRepo          repositories.TaskRepository
	taskStatusLogRepo repositories.TaskStatusLogRepository
	pdfService        PDFService
	activityLogger    utils.ActivityLogger
}

func NewProjectService(repo repositories.ProjectRepository, userRepo repositories.UserRepository, workspaceRepo repositories.WorkspaceRepository, taskRepo repositories.TaskRepository, taskStatusLogRepo repositories.TaskStatusLogRepository, pdfService PDFService, activityLogger utils.ActivityLogger) ProjectService {
	return &projectService{
		repo:              repo,
		userRepo:          userRepo,
		workspaceRepo:     workspaceRepo,
		taskRepo:          taskRepo,
		taskStatusLogRepo: taskStatusLogRepo,
		pdfService:        pdfService,
		activityLogger:    activityLogger,
	}
}

func (s *projectService) CreateProject(project *models.Project, user *models.User) error {
	project.CreatedBy = user.ID

	if err := s.repo.CreateProject(project); err != nil {
		return err
	}

	creatorMember := &models.ProjectUser{
		ProjectID:     project.ID,
		UserID:        user.ID,
		RoleInProject: "admin",
	}

	if err := s.repo.AddMember(creatorMember); err != nil {
		s.repo.DeleteProject(project.ID)
		return errors.New("gagal menambahkan creator sebagai member project")
	}

	return nil
}

func (s *projectService) GetAllProjects(user *models.User) ([]models.Project, error) {
	if user.Role == "admin" {
		return s.repo.GetAllProjects()
	}

	return s.repo.GetProjectsByUserID(user.ID)
}

func (s *projectService) GetByID(projectID uint, user *models.User) (*models.Project, error) {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("project tidak ditemukan")
	}

	if user.Role != "admin" {
		isProjectMember, err := s.repo.IsUserMember(projectID, user.ID)
		if err != nil {
			return nil, errors.New("gagal memeriksa keanggotaan project")
		}

		if !isProjectMember {
			return nil, errors.New("anda tidak memiliki akses ke project ini")
		}
	}

	return project, nil
}

func (s *projectService) UpdateProject(project *models.Project, user *models.User) error {
	existingProject, err := s.repo.GetByID(project.ID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	existingProject.Name = project.Name
	existingProject.Description = project.Description

	return s.repo.UpdateProject(existingProject)
}

func (s *projectService) SoftDeleteProject(projectID uint, user *models.User) error {
	_, err := s.repo.GetByID(projectID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.taskRepo.SoftDeleteAllTasksInProject(projectID); err != nil {
			return fmt.Errorf("gagal soft delete tasks di dalam project: %w", err)
		}

		if err := s.repo.SoftDeleteProject(projectID); err != nil {
			return fmt.Errorf("gagal soft delete project: %w", err)
		}

		return nil
	})
}

func (s *projectService) DeleteProject(projectID uint, user *models.User) error {
	var project models.Project
	err := config.DB.Unscoped().Where("id = ?", projectID).First(&project).Error
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.taskRepo.SoftDeleteAllTasksInProject(projectID); err != nil {
			return fmt.Errorf("gagal soft delete tasks di dalam project: %w", err)
		}

		if err := s.repo.SoftDeleteProject(projectID); err != nil {
			return fmt.Errorf("gagal soft delete project: %w", err)
		}

		return nil
	})
}

func (s *projectService) AddMember(projectID uint, userID uint, role string, currentUser *models.User) error {
	_, err := s.repo.GetByID(projectID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	isMember, err := s.repo.IsUserMember(projectID, userID)
	if err != nil {
		return errors.New("gagal memvalidasi member")
	}
	if isMember {
		return errors.New("user sudah menjadi member di project ini")
	}

	member := &models.ProjectUser{
		ProjectID:     projectID,
		UserID:        userID,
		RoleInProject: role,
	}
	return s.repo.AddMember(member)
}

func (s *projectService) AddMembers(projectID uint, members []ProjectMember, currentUser *models.User) error {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	for _, member := range members {
		isTargetUserInWorkspace, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, member.UserID)
		if err != nil || !isTargetUserInWorkspace {
			return fmt.Errorf("user %d harus menjadi member workspace terlebih dahulu", member.UserID)
		}

		isMember, err := s.repo.IsUserMember(projectID, member.UserID)
		if err != nil {
			return fmt.Errorf("gagal memvalidasi member %d", member.UserID)
		}
		if isMember {
			return fmt.Errorf("user %d sudah menjadi member di project ini", member.UserID)
		}

		projectMember := &models.ProjectUser{
			ProjectID:     projectID,
			UserID:        member.UserID,
			RoleInProject: member.Role,
		}

		if err := s.repo.AddMember(projectMember); err != nil {
			return fmt.Errorf("gagal menambahkan user %d", member.UserID)
		}
	}

	return nil
}

func (s *projectService) GetMembers(projectID uint, user *models.User) ([]models.ProjectUser, error) {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return nil, errors.New("project tidak ditemukan")
	}

	if user.Role == "admin" {
		return s.repo.GetMembers(projectID)
	}

	if user.Role != "admin" {
		isProjectMember, err := s.repo.IsUserMember(projectID, user.ID)
		if err != nil {
			return nil, errors.New("gagal memvalidasi akses project")
		}
		if !isProjectMember {
			return nil, errors.New("akses ditolak untuk melihat members project")
		}

		isWorkspaceMember, err := s.workspaceRepo.IsUserMember(project.WorkspaceID, user.ID)
		if err != nil || !isWorkspaceMember {
			return nil, errors.New("hanya member workspace yang boleh melihat members project")
		}
	}

	return s.repo.GetMembers(projectID)
}

func (s *projectService) RemoveMember(projectID uint, userID uint, currentUser *models.User) error {
	_, err := s.repo.GetByID(projectID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	targetMember, err := s.repo.GetProjectMember(projectID, userID)
	if err != nil {
		return errors.New("member tidak ditemukan di project ini")
	}

	if targetMember.RoleInProject == "admin" {
		return errors.New("tidak bisa menghapus admin project lain")
	}

	return s.repo.RemoveMember(projectID, userID)
}

func (s *projectService) RemoveMembers(projectID uint, userIDs []uint, currentUser *models.User) error {
	_, err := s.repo.GetByID(projectID)
	if err != nil {
		return errors.New("project tidak ditemukan")
	}

	for _, userID := range userIDs {
		targetMember, err := s.repo.GetProjectMember(projectID, userID)
		if err != nil {
			return fmt.Errorf("user %d tidak ditemukan di project ini", userID)
		}

		if userID == currentUser.ID {
			return fmt.Errorf("tidak bisa menghapus diri sendiri (user_id: %d)", userID)
		}

		if targetMember.RoleInProject == "admin" {
			return fmt.Errorf("tidak bisa menghapus admin project (user_id: %d)", userID)
		}
	}

	return s.repo.RemoveMembers(projectID, userIDs)
}

// Private helper function to fetch project and PIC details
func (s *projectService) getProjectAndPIC(projectID uint, userID uint) (*models.Project, *models.User, error) {
	project, err := s.repo.GetByID(projectID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get project: %w", err)
	}

	pic, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get PIC details: %w", err)
	}

	return project, pic, nil
}

func (s *projectService) generatePDFAndLog(
	pdf *gofpdf.Fpdf,
	project *models.Project,
	userID uint,
	reportName string,
) ([]byte, error) {
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to write PDF to buffer: %w", err)
	}

	activity := models.ActivityLog{
		UserID:    userID,
		Action:    fmt.Sprintf("User exported project '%s' - %s", project.Name, reportName),
		TableName: "projects",
		ItemID:    project.ID,
	}
	s.activityLogger.Log(activity)

	return buf.Bytes(), nil
}

// Export 1: Weekly Backward Report
func (s *projectService) ExportWeeklyBackward(projectID uint, userID uint) ([]byte, error) {
	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, -7)

	tasks, err := s.taskRepo.GetTasksStartingBetween(projectID, oneWeekAgo, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var agendaItems []models.AgendaItem
	for _, task := range tasks {
		var memberName string

		if len(task.Members) > 0 && task.Members[0].User.ID != 0 {
			memberName = task.Members[0].User.Name
		} else {
			memberName = "N/A"
		}

		logs, err := s.taskStatusLogRepo.GetLogsByTaskID(task.ID)
		var totalDuration time.Duration
		if err == nil {
			for _, log := range logs {
				if log.ClockOut != nil {
					totalDuration += log.ClockOut.Sub(log.ClockIn)
				}
			}
		}

		agendaItems = append(agendaItems, models.AgendaItem{
			ProjectTitle: task.Project.Name,
			TaskTitle:    task.Title,
			MemberName:   memberName,
			Status:       task.Status,
			Kondisi:      task.Priority,
			StartDate:    task.StartDate,
			DueDate:      task.DueDate,
			Notes:        *task.Notes,
			WorkDuration: formatDuration(totalDuration),
			FinishedAt:   task.FinishedAt,
		})
	}

	period := fmt.Sprintf("%s - %s", oneWeekAgo.Format("02 Jan 2006"), now.Format("02 Jan 2006"))

	sort.SliceStable(agendaItems, func(i, j int) bool {
		if agendaItems[i].MemberName != agendaItems[j].MemberName {
			return agendaItems[i].MemberName < agendaItems[j].MemberName
		}
		return agendaItems[i].StartDate.Before(agendaItems[j].StartDate)
	})

	project, pic, err := s.getProjectAndPIC(projectID, userID)
	if err != nil {
		return nil, err
	}

	pdf, err := s.pdfService.GenerateWeeklyReportPDF(project, agendaItems, *pic, period)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return s.generatePDFAndLog(pdf, project, userID, "Weekly Backward Report")
}

// Export 2: Weekly Forward Report
func (s *projectService) ExportWeeklyForward(projectID uint, userID uint) ([]byte, error) {
	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, 7)

	tasks, err := s.taskRepo.GetTasksStartingBetween(projectID, oneWeekAgo, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	var agendaItems []models.AgendaItem
	for _, task := range tasks {
		var memberName string

		if len(task.Members) > 0 && task.Members[0].User.ID != 0 {
			memberName = task.Members[0].User.Name
		} else {
			memberName = "N/A"
		}

		logs, err := s.taskStatusLogRepo.GetLogsByTaskID(task.ID)
		var totalDuration time.Duration
		if err == nil {
			for _, log := range logs {
				if log.ClockOut != nil {
					totalDuration += log.ClockOut.Sub(log.ClockIn)
				}
			}
		}

		agendaItems = append(agendaItems, models.AgendaItem{
			ProjectTitle: task.Project.Name,
			TaskTitle:    task.Title,
			MemberName:   memberName,
			Status:       task.Status,
			Kondisi:      task.Priority,
			StartDate:    task.StartDate,
			DueDate:      task.DueDate,
			Notes:        *task.Notes,
			WorkDuration: formatDuration(totalDuration),
			FinishedAt:   task.FinishedAt,
		})
	}

	period := fmt.Sprintf("%s - %s", oneWeekAgo.Format("02 Jan 2006"), now.Format("02 Jan 2006"))

	sort.SliceStable(agendaItems, func(i, j int) bool {
		if agendaItems[i].MemberName != agendaItems[j].MemberName {
			return agendaItems[i].MemberName < agendaItems[j].MemberName
		}
		return agendaItems[i].StartDate.Before(agendaItems[j].StartDate)
	})

	project, pic, err := s.getProjectAndPIC(projectID, userID)
	if err != nil {
		return nil, err
	}

	pdf, err := s.pdfService.GenerateWeeklyReportPDF(project, agendaItems, *pic, period)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return s.generatePDFAndLog(pdf, project, userID, "Weekly Forward Report")
}

// Export 3: Daily Report
func (s *projectService) ExportDaily(projectID uint, userID uint) ([]byte, error) {
	now := time.Now()
	year, month, day := now.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	activities, err := s.repo.GetActivityLogsBetween(projectID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	project, pic, err := s.getProjectAndPIC(projectID, userID)
	if err != nil {
		return nil, err
	}

	var dailyItems []models.DailyActivityItem
	statusChangeRegex := regexp.MustCompile(`changed status of task '.*' from '(.*)' to '(.*)'`)

	for _, activity := range activities {
		matches := statusChangeRegex.FindStringSubmatch(activity.Action)
		if len(matches) == 3 {
			task, err := s.taskRepo.GetByID(activity.ItemID)
			if err != nil {
				continue
			}

			user, err := s.userRepo.GetUserByID(activity.UserID)
			if err != nil {
				continue
			}

			newStatus := matches[2]
			dailyItems = append(dailyItems, models.DailyActivityItem{
				ActivityTime: activity.CreatedAt,
				User:         user.Name,
				ProjectTitle: project.Name,
				TaskTitle:    task.Title,
				TaskPriority: task.Priority,
				StatusAtLog:  newStatus,
				Overdue:      task.OverdueDuration,
			})
		}
	}

	period := fmt.Sprintf("Daily Report - %s", now.Format("02 Jan 2006"))

	pdf, err := s.pdfService.GenerateDailyReportPDF(project, dailyItems, *pic, period)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return s.generatePDFAndLog(pdf, project, userID, "Daily Report")
}

// Export 4: Monitoring Report
func (s *projectService) ExportMonitoring(projectID uint, userID uint) ([]byte, error) {
	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, -7)

	inProgressTasks, err := s.taskRepo.GetTasksInProgressSince(projectID, oneWeekAgo)
	if err != nil {
		return nil, err
	}

	doneTasks, err := s.taskRepo.GetTasksDoneSince(projectID, oneWeekAgo)
	if err != nil {
		return nil, err
	}

	tasks := append(inProgressTasks, doneTasks...)

	var tasksWithHistory []models.TaskWithHistory
	for _, task := range tasks {
		logs, err := s.taskStatusLogRepo.GetLogsByTaskID(task.ID)
		if err != nil {
			// Handle error, maybe log it and continue
			continue
		}
		tasksWithHistory = append(tasksWithHistory, models.TaskWithHistory{
			Task:       task,
			StatusLogs: logs,
		})
	}

	period := fmt.Sprintf("%s - %s", oneWeekAgo.Format("02 Jan 2006"), now.Format("02 Jan 2006"))

	project, pic, err := s.getProjectAndPIC(projectID, userID)
	if err != nil {
		return nil, err
	}

	pdf, err := s.pdfService.GenerateMonitoringReportPDF(project, tasksWithHistory, *pic, period)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return s.generatePDFAndLog(pdf, project, userID, "Monitoring Report")
}

// HELPER FORMAT DURATION
func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "-"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dj", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%dd", seconds))
	}
	return strings.Join(parts, " ")
}
