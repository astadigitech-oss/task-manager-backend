package services

import (
	"errors"
	"fmt"
	"time"

	"project-management-backend/models"
	"project-management-backend/repositories"

	"github.com/go-sql-driver/mysql"
)

var ErrAttendanceAlreadyExists = errors.New("attendance for this day already submitted")

type AttendanceService struct {
	repo          repositories.AttendanceRepository
	imageRepo     repositories.AttendanceImageRepository
	userRepo      repositories.UserRepository
	workspaceRepo repositories.WorkspaceRepository
}

func NewAttendanceService(repo repositories.AttendanceRepository, imageRepo repositories.AttendanceImageRepository, userRepo repositories.UserRepository, workspaceRepo repositories.WorkspaceRepository) *AttendanceService {
	return &AttendanceService{
		repo:          repo,
		imageRepo:     imageRepo,
		userRepo:      userRepo,
		workspaceRepo: workspaceRepo,
	}
}

func (s *AttendanceService) GetAttendanceByID(attendanceID uint) (*models.Attendance, error) {
	return s.repo.GetByID(attendanceID)
}

func (s *AttendanceService) SubmitAttendance(attendance *models.Attendance) error {
	hasAccess, err := s.workspaceRepo.IsUserMember(attendance.WorkspaceID, attendance.UserID)
	if err != nil {
		return fmt.Errorf("tidak dapat memeriksa apakah pengguna memiliki akses ke workspace: %w", err)
	}
	if !hasAccess {
		return errors.New("user bukan anggota workspace")
	}

	err = s.repo.Create(attendance)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrAttendanceAlreadyExists
		}
		return err
	}
	return nil
}

func (s *AttendanceService) GetAttendancesForExport(workspaceID uint, date string) ([]models.AttendanceExportResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	attendances, err := s.repo.GetAttendancesByWorkspaceIDAndDateRange(workspaceID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	var exportData []models.AttendanceExportResponse

	for _, att := range attendances {
		user := att.User

		images, err := s.imageRepo.GetByAttendanceID(att.ID)
		var imageURLs []string
		if err == nil {
			for _, img := range images {
				imageURLs = append(imageURLs, img.URL)
			}
		}

		exportData = append(exportData, models.AttendanceExportResponse{
			Attendance: att,
			User: models.UserResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
			},
			ImageURLs: imageURLs,
		})
	}

	return exportData, nil
}
