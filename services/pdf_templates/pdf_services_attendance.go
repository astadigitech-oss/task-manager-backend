package services

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"project-management-backend/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateAttendanceReport(attendances []models.AttendanceExportResponse, workspaceName string, reportDate string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	parsedDate, err := time.Parse("2006-01-02", reportDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse report date: %w", err)
	}
	formattedDate := parsedDate.Format("02 January 2006")

	for i, attendance := range attendances {
		pdf.AddPage()

		// --- HEADER ---
		pdf.Image("assets/logo.png", 15, 15, 25, 0, false, "", 0, "")
		if pdf.Error() != nil {
			return nil, fmt.Errorf("failed to add logo image to PDF: %w", pdf.Error())
		}
		// Header text
		pdf.SetXY(45, 15)
		pdf.SetFont("Arial", "B", 16)
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(140, 8, "LAPORAN ABSENSI HARIAN")
		pdf.Ln(6)

		pdf.SetX(45)
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(140, 6, fmt.Sprintf("Divisi - %s", workspaceName))
		pdf.Ln(6)

		pdf.SetX(45)
		pdf.SetFont("Arial", "", 8)
		pdf.Cell(140, 6, "PT Asta Digital Agency")
		pdf.Ln(4)

		pdf.SetX(45)
		pdf.SetFont("Arial", "", 8)
		pdf.Cell(140, 6, "Imogiri Timur, Gg. Tobanan V, D.I. Yogyakarta")
		pdf.Ln(4)

		pdf.SetX(45)
		pdf.SetFont("Arial", "", 8)
		pdf.Cell(140, 6, "www.astadigitalagency.com")
		pdf.Ln(15)

		// --- GARIS PEMISAH ---
		pdf.SetDrawColor(200, 200, 200)
		pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
		pdf.Ln(10)

		// --- DETAIL ABSENSI ---
		pdf.SetX(15)
		pdf.SetFont("Arial", "B", 12)
		pdf.SetFillColor(240, 240, 240)
		pdf.CellFormat(50, 8, "DATA KARYAWAN", "0", 1, "L", true, 0, "")
		pdf.Ln(7)

		// --- HELPER MULTIPLY LINES ---
		formatMultiCellContent := func(text string) string {
			pdf.SetX(15)
			lines := strings.Split(strings.ReplaceAll(text, " ", " "), " ")
			var formattedLines []string
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed == "" {
					continue
				}
				if (len(trimmed) > 2 && trimmed[1] == '.' && trimmed[0] >= '0' && trimmed[0] <= '9') || (len(trimmed) > 1 && (trimmed[0] == '*' || trimmed[0] == '-')) {
					if len(trimmed) > 2 {
						formattedLines = append(formattedLines, "• "+strings.TrimSpace(trimmed[2:]))
					} else {
						formattedLines = append(formattedLines, "• ")
					}
				} else {
					formattedLines = append(formattedLines, trimmed)
				}
			}
			return strings.Join(formattedLines, " ")
		}

		pdf.SetX(15)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(35, 7, "Nama Karyawan ")
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(100, 7, fmt.Sprintf("               : %s", attendance.User.Name))
		pdf.Ln(10)

		pdf.SetX(15)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(35, 7, "Waktu Absen")
		pdf.SetFont("Arial", "", 11)
		pdf.Cell(100, 7, fmt.Sprintf("               : %s", attendance.ClockIn.Format("15:04:05 WIB")))
		pdf.Ln(10)

		// Kegiatan
		pdf.SetX(15)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, "Kegiatan yang Dilakukan :")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		formattedActivity := formatMultiCellContent(attendance.Activity)
		pdf.MultiCell(0, 6, formattedActivity, "0", "L", false)
		pdf.Ln(5)

		// Kendala
		pdf.SetX(15)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, "Kendala yang Dihadapi    :")
		pdf.Ln(8)
		pdf.SetFont("Arial", "", 11)
		obstacleText := "-"
		if attendance.Obstacle != nil && *attendance.Obstacle != "" {
			obstacleText = *attendance.Obstacle
		}
		formattedObstacle := formatMultiCellContent(obstacleText)
		pdf.MultiCell(0, 6, formattedObstacle, "0", "L", false)
		pdf.Ln(5)

		// Bukti Foto
		pdf.SetX(15)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, "Bukti Foto:")
		pdf.Ln(8)

		if len(attendance.ImageURLs) > 0 {
			// Hanya tampilkan gambar pertama
			imageURL := attendance.ImageURLs[0]
			filePath := "." + imageURL
			imageType := strings.ToUpper(strings.TrimPrefix(filepath.Ext(filePath), "."))
			if imageType == "JPG" {
				imageType = "JPEG"
			}

			pdf.ImageOptions(filePath, pdf.GetX(), pdf.GetY(), 100, 0, false, gofpdf.ImageOptions{ImageType: imageType, ReadDpi: true}, 0, "")
			if pdf.Error() != nil {
				return nil, fmt.Errorf("failed to add attendance image %s to PDF: %w", filePath, pdf.Error())
			}
		} else {
			pdf.SetFont("Arial", "I", 10)
			pdf.Cell(40, 10, "(Tidak ada gambar)")
		}

		// --- FOOTER HALAMAN ---
		pdf.SetY(-40)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Halaman %d dari %d | Diterbitkan: %s", i+1, len(attendances), formattedDate), "", 0, "C", false, 0, "")
	}

	if pdf.Error() != nil {
		return nil, pdf.Error()
	}

	return pdf, nil
}
