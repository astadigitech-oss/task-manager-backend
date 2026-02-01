package services

import (
	"project-management-backend/models"
	pdf_templates "project-management-backend/services/pdf_templates"

	"github.com/jung-kurt/gofpdf"
)

type PDFService interface {
	GenerateAgendaReportPDF(project *models.Project, agendaTasks []models.Task, dailyTasks []models.Task, pic models.User, period string, date string) (*gofpdf.Fpdf, error)
	GenerateDailyReportPDF(project *models.Project, items []models.DailyActivityItem, pic models.User, date string) (*gofpdf.Fpdf, error)
	GenerateWeeklyReportPDF(project *models.Project, tasks []models.Task, pic models.User, date string) (*gofpdf.Fpdf, error)
}

type pdfService struct{}

func NewPDFService() PDFService {
	return &pdfService{}
}

func (s *pdfService) GenerateAgendaReportPDF(project *models.Project, agendaTasks []models.Task, dailyTasks []models.Task, pic models.User, period string, date string) (*gofpdf.Fpdf, error) {
	return pdf_templates.GenerateAgendaReport(project, agendaTasks, dailyTasks, pic, period, date)
}

func (s *pdfService) GenerateDailyReportPDF(project *models.Project, items []models.DailyActivityItem, pic models.User, date string) (*gofpdf.Fpdf, error) {
	return pdf_templates.GenerateDailyReport(project, items, pic, date)
}

func (s *pdfService) GenerateWeeklyReportPDF(project *models.Project, tasks []models.Task, pic models.User, date string) (*gofpdf.Fpdf, error) {
	return pdf_templates.GenerateWeeklyReport(project, tasks, pic, date)
}
