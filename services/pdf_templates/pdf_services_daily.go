package services

import (
	"fmt"
	"time"

	"project-management-backend/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateDailyReport(project *models.Project, items []models.DailyActivityItem, pic models.User, period string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Image("assets/logo.png", 15, 15, 25, 0, false, "", 0, "")
	pdf.SetXY(45, 15)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(140, 8, "LAPORAN HASIL KERJA HARIAN")
	pdf.Ln(5)
	pdf.SetX(45)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(140, 6, "Divisi Tim Maintenance dan Development WMS - Liquid8")
	pdf.Ln(4)
	pdf.SetX(45)
	pdf.Cell(277, 6, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 11)
	meta := [][]string{
		{"Judul Laporan", fmt.Sprintf("Laporan Hasil Kerja Tim %s", project.Name)},
		{"Periode Kerja", period},
		{"Hari Kerja", "Senin - Sabtu"},
		{"Divisi", "Maintenance & Development WMS"},
		{"PIC", fmt.Sprintf("%s [%s]", pic.Name, pic.Role)},
	}
	for _, item := range meta {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(30, 6, item[0])
		pdf.Cell(5, 6, ":")
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(100, 6, item[1])
		pdf.Ln(5)
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	headers := []string{"No", "Terakhir Diperbarui", "Penanggung Jawab", "Agenda", "Sub-Agenda", "Kondisi", "Status Terakhir", "Wkt Resolusi (Menit)"}
	colWidths := []float64{10, 35, 35, 45, 45, 25, 35, 40}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)

	for i, item := range items {
		tableRow(pdf, item, i)
	}

	// for i, task := range tasks {
	// 	var members []string
	// 	for _, member := range task.Members {
	// 		members = append(members, member.User.Name)
	// 	}
	// 	penanggungJawab := strings.Join(members, ", ")
	// 	if len(penanggungJawab) == 0 {
	// 		penanggungJawab = "N/A"
	// 	}

	// 	estimasi := formatDurationEstimasi(task.DueDate.Sub(task.StartDate))

	// 	rowData := []string{
	// 		fmt.Sprintf("%d", i+1),
	// 		item.ActivityTime.Format("15:04"),
	// 		item.User,
	// 		task.Project.Name,
	// 		item.TaskTitle,
	// 		task.Priority,
	// 		task.Status,
	// 		estimasi,
	// 	}

	// 	for j, data := range rowData {
	// 		pdf.CellFormat(colWidths[j], 10, data, "1", 0, "C", false, 0, "")
	// 	}
	// 	pdf.Ln(-1)
	// }
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 10)
	today := time.Now().Format("02-01-2006")
	pdf.Cell(190, 6, fmt.Sprintf("Yogyakarta, %s", today))
	pdf.Ln(5)
	pdf.Cell(190, 6, "Disusun oleh,")
	pdf.Ln(5)
	pdf.Cell(190, 6, "Person In Charge")
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(190, 6, fmt.Sprintf("%s [%s]", pic.Name, pic.Role))

	return pdf, nil
}

func tableRow(pdf *gofpdf.Fpdf, item models.DailyActivityItem, index int) {
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(10, 10, fmt.Sprintf("%d", index+1), "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 10, item.ActivityTime.Format("2006/01/02 15:04"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 10, item.User, "1", 0, "L", false, 0, "")
	pdf.CellFormat(45, 10, item.ProjectTitle, "1", 0, "L", false, 0, "")
	pdf.CellFormat(45, 10, item.TaskTitle, "1", 0, "L", false, 0, "")
	pdf.CellFormat(25, 10, item.TaskPriority, "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 10, item.StatusAtLog, "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, formatDurationEstimasi(item.Overdue), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
}

func formatDurationEstimasi(d time.Duration) string {
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, mins)
}
