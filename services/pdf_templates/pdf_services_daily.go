package services

import (
	"fmt"
	"time"

	"project-management-backend/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateDailyReport(
	project *models.Project,
	items []models.DailyActivityItem,
	pic models.User,
	period string,
) (*gofpdf.Fpdf, error) {

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	pdf.Image("assets/logo.png", 15, 15, 25, 0, false, "", 0, "")
	pdf.SetXY(45, 15)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(140, 8, "LAPORAN HASIL KERJA HARIAN")
	pdf.Ln(6)

	pdf.SetX(45)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(140, 6, "Divisi Tim Maintenance dan Development WMS - Liquid8")
	pdf.Ln(4)

	pdf.SetX(45)
	pdf.Cell(277, 6, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")
	pdf.Ln(15)

	pdf.SetFont("Arial", "", 10)
	meta := [][]string{
		{"Judul Laporan", fmt.Sprintf("Laporan Hasil Kerja Tim %s", project.Name)},
		{"Periode Kerja", period},
		{"Hari Kerja", "Senin - Sabtu"},
		{"Divisi", "Maintenance & Development WMS"},
		{"PIC", fmt.Sprintf("%s [%s]", pic.Name, pic.Role)},
	}

	for _, m := range meta {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(30, 6, m[0])
		pdf.Cell(5, 6, ":")
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(100, 6, m[1])
		pdf.Ln(5)
	}

	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)

	headers := []string{
		"No",
		"Terakhir Diperbarui",
		"Penanggung Jawab",
		"Sub-Agenda",
		"Kondisi",
		"Status Terakhir",
		"Wkt Resolusi (Menit)",
	}

	colWidths := []float64{10, 35, 45, 80, 25, 35, 40}

	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 10, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)

	for i, item := range items {
		drawTableRow(pdf, item, i)
	}

	pdf.Ln(15)
	pdf.SetFont("Arial", "", 10)

	today := time.Now().Format("02-01-2006")
	pdf.Cell(190, 6, fmt.Sprintf("Yogyakarta, %s", today))
	pdf.Ln(6)
	pdf.Cell(190, 6, "Disusun oleh,")
	pdf.Ln(6)
	pdf.Cell(190, 6, "Person In Charge")
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(190, 6, fmt.Sprintf("%s [%s]", pic.Name, pic.Role))

	return pdf, nil
}

func drawTableRow(pdf *gofpdf.Fpdf, item models.DailyActivityItem, index int) {
	lineHeight := 8.0

	hSubAgenda := calcHeight(pdf, item.TaskTitle, 80, lineHeight)

	rowHeight := hSubAgenda
	if hSubAgenda > rowHeight {
		rowHeight = hSubAgenda
	}

	startX := pdf.GetX()
	startY := pdf.GetY()

	pdf.CellFormat(10, rowHeight, fmt.Sprintf("%d", index+1), "1", 0, "C", false, 0, "")

	pdf.CellFormat(35, rowHeight, item.ActivityTime.Format("2006/01/02 15:04"), "1", 0, "C", false, 0, "")

	pdf.CellFormat(45, rowHeight, item.User, "1", 0, "L", false, 0, "")

	x := pdf.GetX()
	y := pdf.GetY()
	pdf.MultiCell(80, lineHeight, item.TaskTitle, "1", "L", false)
	pdf.SetXY(x+80, y)

	pdf.CellFormat(25, rowHeight, item.TaskPriority, "1", 0, "C", false, 0, "")

	pdf.CellFormat(35, rowHeight, item.StatusAtLog, "1", 0, "C", false, 0, "")

	pdf.CellFormat(40, rowHeight, formatDuration(item.Overdue), "1", 0, "C", false, 0, "")

	pdf.SetXY(startX, startY+rowHeight)
}

func calcHeight(pdf *gofpdf.Fpdf, text string, width, lineHeight float64) float64 {
	lines := pdf.SplitLines([]byte(text), width)
	return float64(len(lines)) * lineHeight
}

func formatDurationEstimasi(d time.Duration) string {
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, mins)
}
