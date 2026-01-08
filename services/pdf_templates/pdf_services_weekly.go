package services

import (
	"fmt"
	"strings"
	"time"

	"project-management-backend/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateWeeklyReport(project *models.Project, tasks []models.Task, pic models.User, period string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// --- HEADER ---
	// Logo placeholder
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(40, 10, "ASTA DIGITAL")
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "AGENCY")
	pdf.Ln(10)

	// Main Title
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "LAPORAN HASIL KERJA MINGGUAN")
	pdf.Ln(5)

	// Sub-title
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 10, "Divisi Tim Maintenance dan Development WMS - Liquid8")
	pdf.Ln(4)
	pdf.Cell(190, 10, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")
	pdf.Ln(15)

	// --- METADATA ---
	pdf.SetFont("Arial", "", 11)
	meta := [][]string{
		{"Judul Lapor", fmt.Sprintf("Laporan Hasil Kerja Tim %s", project.Name)},
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

	// --- TABLE HEADER ---
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	headers := []string{"No", "Hari/\nTanggal", "Penanggung\nJawab", "Agenda", "Sub-Agenda", "Kondisi", "Estimasi", "Realisasi"}
	colWidths := []float64{10, 25, 25, 35, 35, 15, 25, 25} // Lebar kolom yang disesuaikan

	headerRowHeight := 12.0
	lineHeight := 5.5

	x := pdf.GetX()
	y := pdf.GetY()

	for i, header := range headers {
		pdf.Rect(x, y, colWidths[i], headerRowHeight, "DF")
		pdf.SetXY(x, y+2)
		pdf.MultiCell(colWidths[i], lineHeight, header, "", "C", false)
		x += colWidths[i]
	}
	pdf.SetY(y + headerRowHeight)

	// --- TABLE ROWS ---
	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255)
	for i, task := range tasks {
		// Handle multiple members
		var members []string
		for _, member := range task.Members {
			members = append(members, member.User.Name)
		}
		penanggungJawab := strings.Join(members, ", ")
		if len(penanggungJawab) == 0 {
			penanggungJawab = "N/A"
		}

		// Calculate estimation
		estimasi := formatDuration(task.DueDate.Sub(task.StartDate))

		var realisasi string
		if task.FinishedAt != nil && !task.FinishedAt.IsZero() {
			// Tugas sudah selesai, hitung durasi aktual
			realisasiDuration := task.FinishedAt.Sub(task.StartDate)
			realisasi = formatDuration(realisasiDuration)
		} else {
			// Tugas belum selesai
			realisasi = "On Progress"
		}

		rowData := []string{
			fmt.Sprintf("%d", i+1),
			task.StartDate.Format("02/01/2006"),
			penanggungJawab,
			task.Project.Name,
			task.Title,
			task.Priority,
			estimasi,
			realisasi,
		}

		// Add cell for each data point
		for j, data := range rowData {
			pdf.CellFormat(colWidths[j], 10, data, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// --- FOOTER ---
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
