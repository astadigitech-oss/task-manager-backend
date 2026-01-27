package services

import (
	"fmt"
	"math"
	"strings"
	"time"

	"project-management-backend/models"

	"github.com/jung-kurt/gofpdf"
)

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "-"
	}
	days := int(math.Floor(d.Hours() / 24))
	hours := int(math.Mod(d.Hours(), 24))

	if days > 0 {
		return fmt.Sprintf("%d Hari", days)
	}
	if hours > 0 {
		return fmt.Sprintf("%d Jam", hours)
	}
	return "Kurang dari 1 jam"
}

func GenerateAgendaReport(project *models.Project, agendaTasks []models.Task, dailyTasks []models.Task, pic models.User, period string, date string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.AddPage()
	generateAgendaPage(pdf, project, agendaTasks, pic, period)

	pdf.AddPage()
	generateDailyPage(pdf, project, dailyTasks, pic, date)

	return pdf, nil
}

func generateAgendaPage(pdf *gofpdf.Fpdf, project *models.Project, tasks []models.Task, pic models.User, period string) {
	// --- HEADER ---
	// pdf.Image("assets/logo.png", 10, 10, 30, 0, false, "", 0, "")
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(277, 6, "LAPORAN HASIL KERJA HARIAN")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(277, 6, "Divisi Tim Maintenance dan Development WMS - Liquid8")
	pdf.Ln(4)
	pdf.Cell(277, 6, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")
	pdf.Ln(15)

	// --- METADATA ---
	pdf.SetFont("Arial", "", 11)
	meta := [][]string{
		{"Judul Lapor", fmt.Sprintf("Laporan Hasil Kerja Harian Tim %s", project.Name)},
		{"Periode Kerja", period},
		{"Hari Kerja", "Senin - Sabtu"},
		{"Divisi", "Maintenance & Development WMS"},
		{"PIC", fmt.Sprintf("%s [admin]", pic.Name)},
	}
	for _, item := range meta {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(40, 6, item[0])
		pdf.Cell(5, 6, ":")
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(100, 6, item[1])
		pdf.Ln(5)
	}
	pdf.Ln(10)

	// --- TABLE HEADER ---
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	headers := []string{"No", "Waktu", "Penanggung Jawab", "Agenda", "Sub-Agenda", "Kondisi", "Status", "Estimasi"}
	colWidths := []float64{10, 27, 35, 30, 30, 20, 20, 20}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// --- TABLE ROWS ---
	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)
	for i, task := range tasks {
		if i >= 25 {
			break
		}
		var members []string
		for _, member := range task.Members {
			members = append(members, member.User.Name)
		}
		penanggungJawab := strings.Join(members, ", ")
		if len(penanggungJawab) == 0 {
			penanggungJawab = "N/A"
		}

		estimasi := formatDuration(task.DueDate.Sub(task.StartDate))
		if !task.DueDate.IsZero() && task.StartDate.IsZero() {
			estimasi = "N/A"
		}

		rowData := []string{
			fmt.Sprintf("%d", i+1),
			task.StartDate.Format("12/27/2025"),
			penanggungJawab,
			project.Name,
			task.Title,
			task.Priority,
			task.Status,
			estimasi,
		}

		for j, data := range rowData {
			pdf.CellFormat(colWidths[j], 10, data, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// --- FOOTER ---
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 10)
	today := time.Now().Format("02-01-2006")
	pdf.Cell(277, 6, fmt.Sprintf("Yogyakarta, %s", today))
	pdf.Ln(5)
	pdf.Cell(277, 6, "Disusun oleh,")
	pdf.Ln(5)
	pdf.Cell(277, 6, "Person In Charge")
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(277, 6, fmt.Sprintf("%s [admin]", pic.Name))
}

func generateDailyPage(pdf *gofpdf.Fpdf, project *models.Project, tasks []models.Task, pic models.User, date string) {
	// --- HEADER ---
	// pdf.Image("assets/logo.png", 10, 10, 30, 0, false, "", 0, "")
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 6, "LAPORAN HASIL KERJA HARIAN")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 6, "Divisi Tim Maintenance dan Development WMS - Liquid8")
	pdf.Ln(4)
	pdf.Cell(190, 6, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")
	pdf.Ln(15)

	// --- METADATA ---
	pdf.SetFont("Arial", "", 11)
	meta := [][]string{
		{"Judul Lapor", fmt.Sprintf("Laporan Hasil Kerja Harian Tim %s", project.Name)},
		{"Hari/Tanggal", date},
		{"Divisi", "Maintenance & Development WMS"},
		{"PIC", fmt.Sprintf("%s [admin]", pic.Name)},
	}
	for _, item := range meta {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(40, 6, item[0])
		pdf.Cell(5, 6, ":")
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(100, 6, item[1])
		pdf.Ln(5)
	}
	pdf.Ln(10)

	// --- TABLE HEADER ---
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	headers := []string{"No", "Jam", "Penanggung Jawab", "Agenda", "Sub-Agenda", "Kondisi", "Status", "Wkt Resolusi (Menit)"}
	colWidths := []float64{10, 15, 35, 30, 35, 20, 25, 20}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// --- TABLE ROWS ---
	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)
	for i, task := range tasks {
		if i >= 25 {
			break
		}
		var members []string
		for _, member := range task.Members {
			members = append(members, member.User.Name)
		}
		penanggungJawab := strings.Join(members, ", ")
		if len(penanggungJawab) == 0 {
			penanggungJawab = "N/A"
		}

		var resolutionTime string
		if task.FinishedAt != nil && !task.FinishedAt.IsZero() && !task.StartDate.IsZero() {
			duration := task.FinishedAt.Sub(task.StartDate)
			resolutionTime = fmt.Sprintf("%.0f", duration.Minutes())
		} else {
			resolutionTime = "0"
		}

		rowData := []string{
			fmt.Sprintf("%d", i+1),
			task.StartDate.Format("02/01/2006"),
			penanggungJawab,
			project.Name,
			task.Title,
			task.Priority,
			task.Status,
			resolutionTime,
		}

		for j, data := range rowData {
			pdf.CellFormat(colWidths[j], 10, data, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Fill remaining rows to make a total of 25
	for i := len(tasks); i < 10; i++ {
		for _, w := range colWidths {
			pdf.CellFormat(w, 10, "", "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// --- FOOTER ---
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 10)
	today := time.Now().Format("02-01-2006") // Perbaiki format date
	pdf.Cell(277, 6, fmt.Sprintf("Yogyakarta, %s", today))
	pdf.Ln(5)
	pdf.Cell(277, 6, "Disusun oleh,")
	pdf.Ln(5)
	pdf.Cell(277, 6, "Person In Charge")
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(277, 6, fmt.Sprintf("%s [admin]", pic.Name))
}
