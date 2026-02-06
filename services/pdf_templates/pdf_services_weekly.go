package services

import (
	"fmt"
	"strings"
	"time"

	"project-management-backend/models"

	"github.com/jung-kurt/gofpdf"
)

func splitTextToLinesWeekly(pdf *gofpdf.Fpdf, text string, maxWidth float64, fontSize float64) []string {
	pdf.SetFontSize(fontSize)
	text = strings.TrimSpace(text)

	if text == "" {
		return []string{"N/A"}
	}

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

func maxInt(values ...int) int {
	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func GenerateWeeklyReport(project *models.Project, tasks []models.Task, pic models.User, period string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	pdf.Image("assets/logo.png", 15, 15, 25, 0, false, "", 0, "")

	pdf.SetXY(45, 15)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(140, 8, "LAPORAN HASIL KERJA MINGGUAN")
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

	colWidths := []float64{10, 30, 40, 80, 25, 25, 25, 30}
	headers := []string{"No", "Hari/Tanggal\nMulai", "Penanggung\nJawab", "Sub-Agenda", "Prioritas", "Estimasi", "Realisasi", "Status"}

	headerY := pdf.GetY()
	headerHeight := 12.0

	for i, header := range headers {
		pdf.SetXY(15+sumWidths(colWidths[:i]), headerY)
		pdf.CellFormat(colWidths[i], headerHeight, "", "1", 0, "", true, 0, "")
		pdf.SetXY(15+sumWidths(colWidths[:i]), headerY+2)
		pdf.MultiCell(colWidths[i], 4, header, "", "C", false)
	}

	pdf.SetY(headerY + headerHeight)

	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)

	for i, task := range tasks {
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

		var realisasi string
		if task.FinishedAt != nil && !task.FinishedAt.IsZero() && !task.StartDate.IsZero() {
			realisasiDuration := task.FinishedAt.Sub(task.StartDate)
			realisasi = formatDuration(realisasiDuration)
		} else {
			realisasi = "On Progress"
		}

		tanggal := "N/A"
		if !task.StartDate.IsZero() {
			tanggal = task.StartDate.Format("02/01/2006")
		}

		penanggungJawabLines := splitTextToLinesWeekly(pdf, penanggungJawab, colWidths[2]-2, 9)
		subAgendaLines := splitTextToLinesWeekly(pdf, task.Title, colWidths[3]-2, 9)

		maxLines := maxInt(len(penanggungJawabLines), len(subAgendaLines))
		rowHeight := float64(maxLines) * 4.5
		if rowHeight < 8 {
			rowHeight = 8
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
			pdf.MultiCell(colWidths[2], 4.5, strings.Join(penanggungJawabLines, "\n"), "1", "L", false)
		}
		currentX += colWidths[2]

		pdf.SetXY(currentX, startY)
		if len(subAgendaLines) == 1 {
			pdf.CellFormat(colWidths[3], rowHeight, subAgendaLines[0], "1", 0, "L", false, 0, "")
		} else {
			pdf.MultiCell(colWidths[3], 4.5, strings.Join(subAgendaLines, "\n"), "1", "L", false)
		}
		currentX += colWidths[3]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[4], rowHeight, task.Priority, "1", 0, "C", false, 0, "")
		currentX += colWidths[4]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[5], rowHeight, estimasi, "1", 0, "C", false, 0, "")
		currentX += colWidths[5]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[6], rowHeight, realisasi, "1", 0, "C", false, 0, "")
		currentX += colWidths[6]

		pdf.SetXY(currentX, startY)
		pdf.CellFormat(colWidths[7], rowHeight, task.Status, "1", 0, "C", false, 0, "")

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
	pdf.Cell(140, 6, fmt.Sprintf("%s [%s]", pic.Name, pic.Role))

	return pdf, nil
}

func sumWidths(widths []float64) float64 {
	total := 0.0
	for _, w := range widths {
		total += w
	}
	return total
}
