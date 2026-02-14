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
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%d Hari", days)
	}
	if hours > 0 {
		return fmt.Sprintf("%d Jam %dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	return "< 1 menit"
}

func splitTextToLines(pdf *gofpdf.Fpdf, text string, maxWidth float64, fontSize float64) []string {
	pdf.SetFontSize(fontSize)
	text = strings.TrimSpace(text)

	if pdf.GetStringWidth(text) <= maxWidth {
		return []string{text}
	}

	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if pdf.GetStringWidth(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	if len(lines) == 0 || (len(lines) == 1 && pdf.GetStringWidth(lines[0]) > maxWidth) {
		for len(text) > 0 {
			if pdf.GetStringWidth(text) <= maxWidth {
				return []string{text}
			}
			text = text[:len(text)-1]
		}
	}

	return lines
}

func GenerateAgendaReport(project *models.Project, items []models.AgendaItem, pic models.User) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	pdf.Image("assets/logo.png", 15, 15, 25, 0, false, "", 0, "")
	pdf.SetXY(45, 15)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(140, 8, "LAPORAN AGENDA KERJA MINGGUAN")
	pdf.Ln(5)

	pdf.SetX(45)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(140, 6, fmt.Sprintf("Divisi - %s - Liquid8", project.Workspace.Name))
	pdf.Ln(4)

	pdf.SetX(45)
	pdf.Cell(140, 6, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 11)
	meta := [][]string{
		{"Judul Laporan", fmt.Sprintf("Laporan Agenda Kerja Mingguan Tim %s", project.Name)},
		{"Periode Kerja", "Senin - Sabtu"},
		{"Hari Kerja", "Senin - Sabtu"},
		{"Divisi", project.Workspace.Name},
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

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)

	colWidths := []float64{10, 35, 45, 30, 20, 22, 22, 20, 53}
	headers := []string{"No", "Judul Proyek", "Judul Tugas", "PIC Tugas", "Status", "Tgl Mulai", "Tgl Selesai", "Durasi", "Catatan"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)
	for i, item := range items {
		pdf.CellFormat(colWidths[0], 10, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[1], 10, item.ProjectTitle, "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidths[2], 10, item.TaskTitle, "1", 0, "L", false, 0, "")
		pdf.CellFormat(colWidths[3], 10, item.MemberName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[4], 10, item.Status, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[5], 10, item.StartDate.Format("02-01-2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[6], 10, item.DueDate.Format("02-01-2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[7], 10, item.WorkDuration, "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[8], 10, item.Notes, "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(15)
	pdf.SetFont("Arial", "", 10)
	today := time.Now().Format("02/01/2006")
	pdf.Cell(140, 6, fmt.Sprintf("Yogyakarta, %s", today))
	pdf.Ln(5)
	pdf.Cell(140, 6, "Disusun oleh,")
	pdf.Ln(5)
	pdf.Cell(140, 6, "Person In Charge")
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(140, 6, fmt.Sprintf("%s [admin]", pic.Name))

	return pdf, nil
}
