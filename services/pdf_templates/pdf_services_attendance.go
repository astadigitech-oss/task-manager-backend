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
		pdf.SetXY(45, 15)
		pdf.SetFont("Arial", "B", 16)
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(140, 8, "Laporan Absensi Harian")
		pdf.Ln(6)

		pdf.SetX(45)
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(140, 6, fmt.Sprintf("Divisi - %s - Liquid8", workspaceName))
		pdf.Ln(6)

		pdf.SetX(45)
		pdf.Cell(140, 6, "www.astadigitalagency | Imogiri Timur, Gg. Tobanan V | D.I.Yogyakarta")
		pdf.Ln(6)

		pdf.SetX(45)
		pdf.Cell(140, 6, formattedDate)
		pdf.Ln(15)

		// --- DETAIL ABSENSI ---
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, attendance.User.Name)
		pdf.Ln(8)
		pdf.SetFont("Arial", "I", 9)
		pdf.Cell(40, 10, fmt.Sprintf("Waktu Absen: %s", attendance.ClockIn.Format("15:04:05 WIB")))
		pdf.Ln(15)

		// Kegiatan
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, "Kegiatan yang Dilakukan:")
		pdf.Ln(10)
		pdf.SetFont("Arial", "", 11)
		pdf.MultiCell(0, 6, attendance.Activity, "0", "L", false)
		pdf.Ln(8)

		// Kendala
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, "Kendala yang Dihadapi:")
		pdf.Ln(10)
		pdf.SetFont("Arial", "", 11)
		obstacleText := "tidak ada kendala"
		if attendance.Obstacle != nil && *attendance.Obstacle != "" {
			obstacleText = *attendance.Obstacle
		}
		pdf.MultiCell(0, 6, obstacleText, "0", "L", false)
		pdf.Ln(8)

		// Bukti Foto
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, "Bukti Foto:")
		pdf.Ln(10)

		if len(attendance.ImageURLs) > 0 {
			for _, imageURL := range attendance.ImageURLs {
				filePath := "." + imageURL

				imageType := strings.ToUpper(strings.TrimPrefix(filepath.Ext(filePath), "."))
				if imageType == "JPG" {
					imageType = "JPEG"
				}

				pdf.ImageOptions(filePath, pdf.GetX(), pdf.GetY(), 100, 0, false, gofpdf.ImageOptions{ImageType: imageType, ReadDpi: true}, 0, "")
				if pdf.Error() != nil {
					return nil, fmt.Errorf("failed to add attendance image %s to PDF: %w", filePath, pdf.Error())
				}
				pdf.Ln(10)
			}
		} else {
			pdf.SetFont("Arial", "I", 10)
			pdf.Cell(40, 10, "(Tidak ada gambar)")
		}

		// --- FOOTER HALAMAN ---
		pdf.SetY(-30)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Halaman %d dari %d", i+1, len(attendances)), "", 0, "C", false, 0, "")
	}

	if pdf.Error() != nil {
		return nil, pdf.Error()
	}

	return pdf, nil
}
