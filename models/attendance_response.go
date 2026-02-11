package models

// AttendanceExportResponse is used for exporting attendance data.
type AttendanceExportResponse struct {
	Attendance
	User      UserResponse `json:"user"`
	ImageURLs []string     `json:"image_urls"`
}
