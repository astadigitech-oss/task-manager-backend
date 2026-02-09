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
	return "Kurang dari 1 menit"
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

func GenerateAgendaReport(project *models.Project, agendaTasks []models.Task, dailyTasks []models.Task, pic models.User, period string, date string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)

	pdf.AddPage()
	generateAgendaPage(pdf, project, agendaTasks, pic, period)

	pdf.AddPage()
	generateDailyPage(pdf, project, dailyTasks, pic, date)

	return pdf, nil
}

func generateAgendaPage(pdf *gofpdf.Fpdf, project *models.Project, tasks []models.Task, pic models.User, period string) {
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
		{"Periode Kerja", period},
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

	colWidths := []float64{10, 25, 45, 90, 30, 30, 40}
	headers := []string{"No", "Tanggal", "Penanggung Jawab", "Sub-Agenda", "Prioritas", "Status", "Estimasi"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

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

		var estimasi string
		if !task.DueDate.IsZero() && !task.StartDate.IsZero() {
			duration := task.DueDate.Sub(task.StartDate)
			estimasi = formatDuration(duration)
		} else {
			estimasi = "N/A"
		}

		tanggal := "N/A"
		if !task.StartDate.IsZero() {
			tanggal = task.StartDate.Format("02/01/2006")
		}

		penanggungJawabLines := splitTextToLines(pdf, penanggungJawab, colWidths[2]-2, 9)
		subAgendaLines := splitTextToLines(pdf, task.Title, colWidths[3]-2, 9)

		maxLines := max(len(penanggungJawabLines), len(subAgendaLines))
		rowHeight := float64(maxLines) * 5
		if rowHeight < 10 {
			rowHeight = 10
		}

		startY := pdf.GetY()
		currentX := 15.0

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[0], rowHeight, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		currentX += colWidths[0]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[1], rowHeight, tanggal, "1", 0, "C", false, 0, "")
		currentX += colWidths[1]

		pdf.SetXY(currentX, startY)
		if len(penanggungJawabLines) == 1 {
			pdf.CellFormat(colWidths[2], rowHeight, penanggungJawabLines[0], "1", 0, "L", false, 0, "")
		} else {
			pdf.MultiCell(colWidths[2], 5, strings.Join(penanggungJawabLines, "\n"), "1", "L", false)
		}
		currentX += colWidths[2]

		pdf.SetXY(currentX, startY)
		if len(subAgendaLines) == 1 {
			pdf.CellFormat(colWidths[3], rowHeight, subAgendaLines[0], "1", 0, "L", false, 0, "")
		} else {
			pdf.MultiCell(colWidths[3], 5, strings.Join(subAgendaLines, "\n"), "1", "L", false)
		}
		currentX += colWidths[3]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[4], rowHeight, task.Priority, "1", 0, "C", false, 0, "")
		currentX += colWidths[4]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[5], rowHeight, task.Status, "1", 0, "C", false, 0, "")
		currentX += colWidths[5]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[6], rowHeight, estimasi, "1", 0, "C", false, 0, "")

		pdf.SetY(startY + rowHeight)
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
}

func generateDailyPage(pdf *gofpdf.Fpdf, project *models.Project, tasks []models.Task, pic models.User, date string) {
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
	pdf.Cell(140, 6, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 11)
	meta := [][]string{
		{"Judul Laporan", fmt.Sprintf("Laporan Hasil Kerja Harian Tim %s", project.Name)},
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

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)

	colWidths := []float64{10, 25, 45, 90, 30, 30, 40}
	headers := []string{"No", "Jam", "Penanggung Jawab", "Sub-Agenda", "Prioritas", "Status", "Waktu Resolusi"}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

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
			resolutionTime = formatDurationEstimasi(duration)
		} else {
			resolutionTime = "0m"
		}

		jam := "N/A"
		if !task.StartDate.IsZero() {
			jam = task.StartDate.Format("15:04")
		}

		penanggungJawabLines := splitTextToLines(pdf, penanggungJawab, colWidths[2]-2, 9)
		subAgendaLines := splitTextToLines(pdf, task.Title, colWidths[3]-2, 9)

		maxLines := max(len(penanggungJawabLines), len(subAgendaLines))
		rowHeight := float64(maxLines) * 5
		if rowHeight < 10 {
			rowHeight = 10
		}

		startY := pdf.GetY()
		currentX := 15.0

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[0], rowHeight, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		currentX += colWidths[0]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[1], rowHeight, jam, "1", 0, "C", false, 0, "")
		currentX += colWidths[1]

		pdf.SetXY(currentX, startY)
		if len(penanggungJawabLines) == 1 {
			pdf.CellFormat(colWidths[2], rowHeight, penanggungJawabLines[0], "1", 0, "L", false, 0, "")
		} else {
			pdf.MultiCell(colWidths[2], 5, strings.Join(penanggungJawabLines, "\n"), "1", "L", false)
		}
		currentX += colWidths[2]

		pdf.SetXY(currentX, startY)
		if len(subAgendaLines) == 1 {
			pdf.CellFormat(colWidths[3], rowHeight, subAgendaLines[0], "1", 0, "L", false, 0, "")
		} else {
			pdf.MultiCell(colWidths[3], 5, strings.Join(subAgendaLines, "\n"), "1", "L", false)
		}
		currentX += colWidths[3]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[4], rowHeight, task.Priority, "1", 0, "C", false, 0, "")
		currentX += colWidths[4]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[5], rowHeight, task.Status, "1", 0, "C", false, 0, "")
		currentX += colWidths[5]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[6], rowHeight, resolutionTime, "1", 0, "C", false, 0, "")

		pdf.SetY(startY + rowHeight)
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
}

func max(values ...int) int {
	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}
