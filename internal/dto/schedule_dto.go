package dto

import (
	"quickattendance-go/internal/domain"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateScheduleRequest struct {
	Name               string      `json:"name"`
	DaysOfWeek         []int       `json:"days_of_week"`
	EntryTimeMinutes   int         `json:"entry_time_minutes"`
	ExitTimeMinutes    int         `json:"exit_time_minutes"`
	GracePeriodMinutes int         `json:"grace_period_minutes"`
	IsDefault          bool        `json:"is_default"`
	AssignedUsersIDs   []uuid.UUID `json:"assigned_users_ids"`
}

type UpdateScheduleRequest struct {
	Name               *string      `json:"name"`
	DaysOfWeek         *[]int       `json:"days_of_week"`
	EntryTimeMinutes   *int         `json:"entry_time_minutes"`
	ExitTimeMinutes    *int         `json:"exit_time_minutes"`
	GracePeriodMinutes *int         `json:"grace_period_minutes"`
	IsDefault          *bool        `json:"is_default"`
	AssignedUsersIDs   *[]uuid.UUID `json:"assigned_users_ids"`
}

type ScheduleResponse struct {
	ID       uuid.UUID `json:"id"`
	AgencyID uuid.UUID `json:"agency_id"`
	Name     string    `json:"name"`
	// El cliente deberia recibir un arreglo de enteros para los dias de la semana
	DaysOfWeek         []int          `json:"days_of_week"`
	EntryTimeMinutes   int            `json:"entry_time_minutes"`
	ExitTimeMinutes    int            `json:"exit_time_minutes"`
	GracePeriodMinutes int            `json:"grace_period_minutes"`
	IsDefault          bool           `json:"is_default"`
	AssignedUsers      []UserResponse `json:"assigned_users"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

func ToScheduleResponse(schedule *domain.Schedule) *ScheduleResponse {
	if schedule == nil {
		return nil
	}

	days := []int{}
	for str := range strings.SplitSeq(schedule.DaysOfWeek, ",") {
		if str == "" {
			continue
		}

		if val, err := strconv.Atoi(str); err == nil {
			days = append(days, val)
		}
	}

	users := []UserResponse{}
	for _, user := range schedule.AssignedUsers {
		users = append(users, *ToUserResponse(&user))
	}

	return &ScheduleResponse{
		ID:                 schedule.ID,
		AgencyID:           schedule.AgencyID,
		Name:               schedule.Name,
		DaysOfWeek:         days,
		EntryTimeMinutes:   schedule.EntryTimeMinutes,
		ExitTimeMinutes:    schedule.ExitTimeMinutes,
		GracePeriodMinutes: schedule.GracePeriodMinutes,
		IsDefault:          schedule.IsDefault,
		AssignedUsers:      users,
		CreatedAt:          schedule.CreatedAt,
		UpdatedAt:          schedule.UpdatedAt,
	}
}

type ScheduleListParams struct {
	PaginationParams
	Name      string `form:"name" binding:"omitempty"`
	IsDefault *bool  `form:"is_default" binding:"omitempty"`
}

type GetApplicableScheduleParams struct {
	UserID uuid.UUID `form:"user_id" binding:"omitempty"`
	Date   string    `form:"date" binding:"required"` // Format: YYYY-MM-DD
}
