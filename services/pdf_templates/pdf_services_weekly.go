package services

import (
	"fmt"
	"strings"
	"time"

	"project-management-backend/models"

	"github.com/jung-kurt/gofpdf"
)

const (
	lineHeight = 5.0
)

func drawHeader(pdf *gofpdf.Fpdf, project *models.Project) {
	pdf.Image("assets/logo.png", 15, 15, 25, 0, false, "", 0, "")

	pdf.SetXY(45, 15)
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 8, "LAPORAN HASIL KERJA MINGGUAN")

	pdf.Ln(6)
	pdf.SetX(45)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("%s - %s", project.Workspace.Name, project.Name))

	pdf.Ln(5)
	pdf.SetX(45)
	pdf.Cell(0, 6, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")

	pdf.Ln(15)
}

func drawMeta(pdf *gofpdf.Fpdf, project *models.Project, pic models.User, period string) {
	meta := [][]string{
		{"Judul Lapor", fmt.Sprintf("Laporan Hasil Kerja Tim %s", project.Workspace.Name)},
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
}

func drawTableHeader(pdf *gofpdf.Fpdf) {

	headers := []string{
		"No", "Tugas", "Penanggung Jawab",
		"Status", "Kondisi", "Tgl Mulai",
		"Deadline", "Waktu Selesai", "Catatan",
	}

	colWidths := []float64{10, 45, 45, 25, 25, 25, 25, 25, 40}

	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(230, 230, 230)

	x := pdf.GetX()
	y := pdf.GetY()

	for i, header := range headers {
		pdf.Rect(x, y, colWidths[i], 10, "DF")
		pdf.SetXY(x, y+3)
		pdf.MultiCell(colWidths[i], lineHeight, header, "", "C", false)
		x += colWidths[i]
	}

	pdf.SetXY(15, y+10)
}
func GenerateWeeklyReportPDF(project *models.Project, items []models.AgendaItem, pic models.User, period string) (*gofpdf.Fpdf, error) {

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(false, 15)
	pdf.AddPage()

	drawHeader(pdf, project)
	drawMeta(pdf, project, pic, period)
	drawTableHeader(pdf)

	colWidths := []float64{10, 45, 45, 25, 25, 25, 25, 25, 40}

	pdf.SetFont("Arial", "", 9)

	for i, item := range items {
		waktuSelesai := "-"
		if strings.ToLower(item.Status) == "done" {
			waktuSelesai = item.FinishedAt.Format("02-01-2006")
		}

		rowData := []string{
			fmt.Sprintf("%d", i+1),
			item.TaskTitle,
			item.MemberName,
			item.Status,
			item.Kondisi,
			item.StartDate.Format("02-01-2006"),
			item.DueDate.Format("02-01-2006"),
			waktuSelesai,
			item.Notes,
		}

		rowHeight := calculateRowHeight(pdf, rowData, colWidths)

		// Auto page break
		if pdf.GetY()+rowHeight > 190 {
			pdf.AddPage()
			drawHeader(pdf, project)
			drawMeta(pdf, project, pic, period)
			drawTableHeader(pdf)
		}

		drawRow(pdf, rowData, colWidths, rowHeight, i)
	}

	drawFooter(pdf, pic)

	return pdf, nil
}

func drawRow(pdf *gofpdf.Fpdf, row []string, colWidths []float64, rowHeight float64, rowIndex int) {

	startX := pdf.GetX()
	startY := pdf.GetY()

	// Zebra row
	if rowIndex%2 == 0 {
		pdf.SetFillColor(250, 250, 250)
	} else {
		pdf.SetFillColor(255, 255, 255)
	}

	x := startX

	for i, text := range row {

		// Status color
		if i == 4 {
			setStatusColor(pdf, text)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}

		pdf.Rect(x, startY, colWidths[i], rowHeight, "DF")
		pdf.SetXY(x, startY)
		pdf.MultiCell(colWidths[i], lineHeight, text, "", getAlign(i), false)
		x += colWidths[i]
	}

	pdf.SetXY(startX, startY+rowHeight)
}

func calculateRowHeight(pdf *gofpdf.Fpdf, row []string, colWidths []float64) float64 {

	maxLines := 1

	for i, text := range row {
		lines := pdf.SplitLines([]byte(text), colWidths[i])
		if len(lines) > maxLines {
			maxLines = len(lines)
		}
	}

	return float64(maxLines) * lineHeight
}

func getAlign(index int) string {
	centerColumns := []int{0, 3, 4, 5, 6, 7}
	for _, c := range centerColumns {
		if c == index {
			return "C"
		}
	}
	return "L"
}

func setStatusColor(pdf *gofpdf.Fpdf, status string) {
	switch strings.ToLower(status) {
	case "done":
		pdf.SetTextColor(0, 150, 0)
	case "on_progress":
		pdf.SetTextColor(0, 0, 200)
	case "on_board":
		pdf.SetTextColor(200, 120, 0)
	default:
		pdf.SetTextColor(0, 0, 0)
	}
}

func drawFooter(pdf *gofpdf.Fpdf, pic models.User) {

	pdf.Ln(15)

	pdf.SetFont("Arial", "", 10)
	today := time.Now().Format("02-01-2006")

	pdf.Cell(0, 6, fmt.Sprintf("Yogyakarta, %s", today))
	pdf.Ln(6)
	pdf.Cell(0, 6, "Disusun oleh,")
	pdf.Ln(6)
	pdf.Cell(0, 6, "Person In Charge")
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("%s [%s]", pic.Name, pic.Role))
}
