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
	ExportAgenda(projectID uint, userID uint) ([]byte, error)
}

type ProjectMember struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role_in_project"`
}

type projectService struct {
	repo           repositories.ProjectRepository
	userRepo       repositories.UserRepository
	workspaceRepo  repositories.WorkspaceRepository
	taskRepo       repositories.TaskRepository
	pdfService     PDFService
	activityLogger utils.ActivityLogger
}

func NewProjectService(repo repositories.ProjectRepository, userRepo repositories.UserRepository, workspaceRepo repositories.WorkspaceRepository, taskRepo repositories.TaskRepository, pdfService PDFService, activityLogger utils.ActivityLogger) ProjectService {
	return &projectService{
		repo:           repo,
		userRepo:       userRepo,
		workspaceRepo:  workspaceRepo,
		taskRepo:       taskRepo,
		pdfService:     pdfService,
		activityLogger: activityLogger,
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

	inProgressTasks, err := s.taskRepo.GetTasksInProgressSince(projectID, oneWeekAgo)
	if err != nil {
		return nil, err
	}

	onBoardTasks, err := s.taskRepo.GetTasksOnBoardSince(projectID, oneWeekAgo)
	if err != nil {
		return nil, err
	}

	doneTasks, err := s.taskRepo.GetTasksDoneSince(projectID, oneWeekAgo)
	if err != nil {
		return nil, err
	}

	tasks := append(inProgressTasks, doneTasks...)
	tasks = append(tasks, onBoardTasks...)
	period := fmt.Sprintf("%s - %s", oneWeekAgo.Format("02 Jan 2006"), now.Format("02 Jan 2006"))

	project, pic, err := s.getProjectAndPIC(projectID, userID)
	if err != nil {
		return nil, err
	}

	pdf, err := s.pdfService.GenerateWeeklyReportPDF(project, tasks, *pic, period)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return s.generatePDFAndLog(pdf, project, userID, "Weekly Backward Report")
}

// Export 2: Weekly Forward Report
func (s *projectService) ExportWeeklyForward(projectID uint, userID uint) ([]byte, error) {
	now := time.Now()
	oneWeekLater := now.AddDate(0, 0, 7)

	startingTasks, err := s.taskRepo.GetTasksStartingBetween(projectID, now, oneWeekLater)
	if err != nil {
		return nil, err
	}

	dueTasks, err := s.taskRepo.GetOnProgressTasksDueBetween(projectID, now, oneWeekLater)
	if err != nil {
		return nil, err
	}

	tasks := append(startingTasks, dueTasks...)
	period := fmt.Sprintf("%s - %s", now.Format("02 Jan 2006"), oneWeekLater.Format("02 Jan 2006"))

	project, pic, err := s.getProjectAndPIC(projectID, userID)
	if err != nil {
		return nil, err
	}

	pdf, err := s.pdfService.GenerateWeeklyReportPDF(project, tasks, *pic, period)
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

	// var tasks []models.Task
	// for _, activity := range activities {
	// 	task, err := s.taskRepo.GetByID(activity.ItemID)
	// 	if err == nil {
	// 		tasks = append(tasks, *task)
	// 	}
	// }
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

			// Log a single, combined state change.
			newStatus := matches[2]
			dailyItems = append(dailyItems, models.DailyActivityItem{
				ActivityTime: activity.CreatedAt,
				User:         user.Name,
				TaskTitle:    task.Title,
				TaskPriority: task.Priority,
				StatusAtLog:  newStatus,
				Overdue:      task.OverdueDuration,
			})
		}
	}

	period := fmt.Sprintf("Daily Report - %s", now.Format("02 Jan 2006"))

	project, pic, err := s.getProjectAndPIC(projectID, userID)
	if err != nil {
		return nil, err
	}

	pdf, err := s.pdfService.GenerateDailyReportPDF(project, dailyItems, *pic, period)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return s.generatePDFAndLog(pdf, project, userID, "Daily Report")
}

// Export 4: Agenda Report
func (s *projectService) ExportAgenda(projectID uint, userID uint) ([]byte, error) {
	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, -7)
	oneWeekLater := now.AddDate(0, 0, 7)

	inProgressTasks, err := s.taskRepo.GetTasksInProgressSince(projectID, oneWeekAgo)
	if err != nil {
		return nil, err
	}
	doneTasks, err := s.taskRepo.GetTasksDoneSince(projectID, oneWeekAgo)
	if err != nil {
		return nil, err
	}

	startingTasks, err := s.taskRepo.GetTasksStartingBetween(projectID, now, oneWeekLater)
	if err != nil {
		return nil, err
	}

	onBoardTasks, err := s.taskRepo.GetTasksOnBoardSince(projectID, oneWeekAgo)
	if err != nil {
		return nil, err
	}

	dueTasks, err := s.taskRepo.GetOnProgressTasksDueBetween(projectID, now, oneWeekLater)
	if err != nil {
		return nil, err
	}

	// --- Combine all tasks ---
	var tasks []models.Task
	tasks = append(tasks, inProgressTasks...)
	tasks = append(tasks, doneTasks...)
	tasks = append(tasks, onBoardTasks...)
	tasks = append(tasks, startingTasks...)
	tasks = append(tasks, dueTasks...)

	// Define the period string
	period := fmt.Sprintf("%s - %s", oneWeekAgo.Format("02 Jan 2006"), oneWeekLater.Format("02 Jan 2006"))

	// --- Data for Daily Report ---
	oneDayAgo := now.AddDate(0, 0, -1)
	dailyInProgress, _ := s.taskRepo.GetTasksInProgressSince(projectID, oneDayAgo)
	dailyDone, _ := s.taskRepo.GetTasksDoneSince(projectID, oneDayAgo)
	dailyOnBoard, _ := s.taskRepo.GetTasksOnBoardSince(projectID, oneDayAgo)

	var dailyTasks []models.Task
	dailyTasks = append(dailyTasks, dailyInProgress...)
	dailyTasks = append(dailyTasks, dailyDone...)
	dailyTasks = append(dailyTasks, dailyOnBoard...)
	dailyDate := now.Format("Monday, 02-01-2006")

	project, pic, err := s.getProjectAndPIC(projectID, userID)
	if err != nil {
		return nil, err
	}

	pdf, err := s.pdfService.GenerateAgendaReportPDF(project, tasks, dailyTasks, *pic, period, dailyDate)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Use the generic export function
	return s.generatePDFAndLog(pdf, project, userID, "Agenda Report")
}
