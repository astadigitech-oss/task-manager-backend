package services

import (
	"fmt"
	"math"
	"strings"
	"time"

	"project-management-backend/models"

	"github.com/jung-kurt/gofpdf"
)

func maxInt(values ...int) int {
	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func drawWeeklyReportHeader(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	headers := []string{"No", "Hari/Tanggal\nMulai", "Penanggung Jawab", "Tugas", "Kondisi", "Estimasi", "Status", "Durasi", "Catatan"}
	colWidths := []float64{10, 25, 30, 45, 25, 25, 30, 30, 45}
	headerRowHeight := 12.0
	lineHeight := 5.5
	x := 15.0
	y := pdf.GetY()
	for i, header := range headers {
		pdf.Rect(x, y, colWidths[i], headerRowHeight, "DF")
		pdf.SetXY(x, y+2)
		pdf.MultiCell(colWidths[i], lineHeight, header, "", "C", false)
		x += colWidths[i]
	}
	pdf.SetXY(15.0, y+headerRowHeight)
}

func GenerateWeeklyReportPDF(project *models.Project, tasksWithHistory []models.TaskWithHistory, pic models.User, period string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.Image("assets/logo.png", 15, 15, 25, 0, false, "", 0, "")
	pdf.SetXY(45, 15)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(140, 8, "LAPORAN HASIL KERJA MINGGUAN")
	pdf.Ln(5)
	pdf.SetX(45)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(140, 6, fmt.Sprintf("Divisi %s - %s", project.Workspace.Name, project.Name))
	pdf.Ln(4)
	pdf.SetX(45)
	pdf.Cell(277, 6, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")
	pdf.Ln(15)

	// Meta Info
	pdf.SetFont("Arial", "", 11)
	meta := [][]string{
		{"Judul Lapor", fmt.Sprintf("Laporan Hasil Kerja Tim %s", project.Name)},
		{"Periode Kerja", period},
		{"Hari Kerja", "Senin - Sabtu"},
		{"Agenda", project.Name},
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

	drawWeeklyReportHeader(pdf)

	pdf.SetFont("Arial", "", 9)
	cellHeight := 10.0
	no := 1
	colWidths := []float64{10, 25, 30, 45, 25, 25, 30, 30, 45}
	for _, taskWithHistory := range tasksWithHistory {
		pdf.SetFillColor(255, 255, 255)
		task := taskWithHistory.Task
		var members []string
		for _, member := range task.Members {
			members = append(members, member.User.Name)
		}
		penanggungJawab := strings.Join(members, ", ")
		if len(penanggungJawab) == 0 {
			penanggungJawab = "N/A"
		}

		estimasi := formatDuration(task.DueDate.Sub(task.StartDate))

		spanningData := []string{
			fmt.Sprintf("%d", no),
			task.StartDate.Format("01-01-2006"),
			penanggungJawab,
			task.Title,
			task.Priority,
			estimasi,
		}
		spanningDataWidths := colWidths[0:7]

		maxLines := 1
		for i, data := range spanningData {
			lines := pdf.SplitLines([]byte(data), spanningDataWidths[i])
			if len(lines) > maxLines {
				maxLines = len(lines)
			}
		}

		lines := pdf.SplitLines([]byte(task.Description), colWidths[8])
		if len(lines) > maxLines {
			maxLines = len(lines)
		}

		mainDataHeight := float64(maxLines)*6.0 + 4.0

		statusLogsHeight := float64(len(taskWithHistory.StatusLogs)) * cellHeight
		if statusLogsHeight == 0 {
			statusLogsHeight = cellHeight
		}

		overallRowHeight := math.Max(mainDataHeight, statusLogsHeight)

		if pdf.GetY()+overallRowHeight > 185 {
			pdf.AddPage()
			drawWeeklyReportHeader(pdf)
			pdf.SetFillColor(255, 255, 255)
			pdf.SetFont("Arial", "", 9)
		}

		startY := pdf.GetY()
		currentX := 15.0

		for i, data := range spanningData {
			pdf.Rect(currentX, startY, spanningDataWidths[i], overallRowHeight, "DF")
			pdf.SetXY(currentX+1, startY+2)
			pdf.MultiCell(spanningDataWidths[i]-2, 6, data, "", "C", false)
			currentX += spanningDataWidths[i]
		}

		if len(taskWithHistory.StatusLogs) > 0 {
			for i, log := range taskWithHistory.StatusLogs {
				pdf.SetXY(currentX, startY+float64(i)*cellHeight)
				var duration string
				if log.ClockOut != nil && !log.ClockOut.IsZero() {
					duration = formatDuration(log.ClockOut.Sub(log.ClockIn))
				} else {
					duration = "-"
				}
				pdf.CellFormat(colWidths[6], cellHeight, log.Status, "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidths[7], cellHeight, duration, "1", 0, "C", false, 0, "")
			}
		} else {

			pdf.SetXY(currentX, startY)
			pdf.CellFormat(colWidths[6], overallRowHeight, task.Status, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[7], overallRowHeight, "N/A", "1", 0, "C", false, 0, "")
		}

		notesX := currentX + colWidths[6] + colWidths[7]
		pdf.Rect(notesX, startY, colWidths[8], overallRowHeight, "DF")
		pdf.SetXY(notesX+1, startY+2)
		pdf.MultiCell(colWidths[8]-2, 6, task.Description, "", "L", false)

		pdf.SetY(startY + overallRowHeight)
		no++
	}

	if pdf.GetY() > 160 {
		pdf.AddPage()
	}
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

func sumWidths(widths []float64) float64 {
	total := 0.0
	for _, w := range widths {
		total += w
	}
	return total
}
