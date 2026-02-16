package services

import (
	"bytes"
	"project-management-backend/models"
	pdf_templates "project-management-backend/services/pdf_templates"

	"github.com/jung-kurt/gofpdf"
)

type PDFService interface {
	GenerateMonitoringReportPDF(project *models.Project, tasks []models.TaskWithHistory, pic models.User, period string) (*gofpdf.Fpdf, error)
	GenerateDailyReportPDF(project *models.Project, items []models.DailyActivityItem, pic models.User, date string) (*gofpdf.Fpdf, error)
	GenerateWeeklyReportPDF(project *models.Project, agendaItems []models.AgendaItem, pic models.User, period string) (*gofpdf.Fpdf, error)
	CreateAttendanceReportPDF(attendances []models.AttendanceExportResponse, workspaceName string, date string) ([]byte, error)
}

type pdfService struct{}

func NewPDFService() PDFService {
	return &pdfService{}
}

func (s *pdfService) GenerateMonitoringReportPDF(project *models.Project, tasks []models.TaskWithHistory, pic models.User, period string) (*gofpdf.Fpdf, error) {
	return pdf_templates.GenerateMonitoringReportPDF(project, tasks, pic, period)
}

func (s *pdfService) GenerateDailyReportPDF(project *models.Project, items []models.DailyActivityItem, pic models.User, date string) (*gofpdf.Fpdf, error) {
	return pdf_templates.GenerateDailyReport(project, items, pic, date)
}

func (s *pdfService) GenerateWeeklyReportPDF(project *models.Project, agendaItems []models.AgendaItem, pic models.User, period string) (*gofpdf.Fpdf, error) {
	return pdf_templates.GenerateWeeklyReportPDF(project, agendaItems, pic, period)
}

func (s *pdfService) CreateAttendanceReportPDF(attendances []models.AttendanceExportResponse, workspaceName string, date string) ([]byte, error) {
	pdf, err := pdf_templates.GenerateAttendanceReport(attendances, workspaceName, date)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
