package dto

import (
	"quickattendance-go/internal/domain"
	"time"

	"github.com/google/uuid"
)

type MarkAttendanceRequest struct {
	UserID        uuid.UUID               `json:"user_id"`
	AgencyID      uuid.UUID               `json:"agency_id"`
	Method        domain.AttendanceMethod `json:"method"`
	Type          domain.AttendanceType   `json:"type"`
	Notes         *string                 `json:"notes"`
	IsRemote      *bool                   `json:"is_remote"`
	Latitude      *float64                `json:"latitude"`
	Longitude     *float64                `json:"longitude"`
	RequesterRole domain.Role             `json:"-"`
}

type AttendanceResponse struct {
	ID                uuid.UUID                `json:"id"`
	UserID            uuid.UUID                `json:"user_id"`
	AgencyID          uuid.UUID                `json:"agency_id"`
	CheckInTime       time.Time                `json:"check_in_time"`
	ScheduleEntryTime time.Time                `json:"schedule_entry_time"`
	Status            domain.AttendanceStatus  `json:"status"`
	CheckOutTime      *time.Time               `json:"check_out_time"`
	ScheduleExitTime  time.Time                `json:"schedule_exit_time"`
	Date              time.Time                `json:"date"`
	MethodIn          domain.AttendanceMethod  `json:"method_in"`
	MethodOut         *domain.AttendanceMethod `json:"method_out"`
	Notes             *string                  `json:"notes"`
	Latitude          *float64                 `json:"latitude"`
	Longitude         *float64                 `json:"longitude"`
}

func ToAttendanceResponse(attendance *domain.Attendance) *AttendanceResponse {
	if attendance == nil {
		return nil
	}

	return &AttendanceResponse{
		ID:                attendance.ID,
		UserID:            attendance.UserID,
		AgencyID:          attendance.AgencyID,
		CheckInTime:       attendance.CheckInTime,
		ScheduleEntryTime: attendance.ScheduleEntryTime,
		Status:            attendance.Status,
		CheckOutTime:      attendance.CheckOutTime,
		ScheduleExitTime:  attendance.ScheduleExitTime,
		Date:              attendance.Date,
		MethodIn:          attendance.MethodIn,
		MethodOut:         attendance.MethodOut,
		Notes:             attendance.Notes,
		Latitude:          attendance.Latitude,
		Longitude:         attendance.Longitude,
	}
}

type AttendanceListParams struct {
	PaginationParams
	UserID    string `form:"user_id" binding:"omitempty"`
	StartDate string `form:"start_date" binding:"omitempty"` // Format: YYYY-MM-DD
	EndDate   string `form:"end_date" binding:"omitempty"`   // Format: YYYY-MM-DD
	Status    string `form:"status" binding:"omitempty"`
}
