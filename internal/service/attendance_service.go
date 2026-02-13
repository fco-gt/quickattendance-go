package service

import (
	"context"
	"quickattendance-go/internal/domain"
	"quickattendance-go/internal/dto"
	"quickattendance-go/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type AttendanceService struct {
	attendanceRepo domain.AttendanceRepo
	userRepo       domain.UserRepo
	scheduleSvc    *ScheduleService
	transactor     domain.Transactor
}

func NewAttendanceService(
	attendanceRepo domain.AttendanceRepo,
	userRepo domain.UserRepo,
	scheduleSvc *ScheduleService,
	transactor domain.Transactor,
) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		userRepo:       userRepo,
		scheduleSvc:    scheduleSvc,
		transactor:     transactor,
	}
}

func (s *AttendanceService) MarkAttendance(ctx context.Context, req *dto.MarkAttendanceRequest) (*dto.AttendanceResponse, error) {
	now := time.Now()

	sched, err := s.scheduleSvc.GetApplicableSchedule(ctx, req.AgencyID, req.UserID, now)
	if err != nil {
		return nil, err
	}

	var response *dto.AttendanceResponse

	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		existing, err := s.attendanceRepo.GetTodayByUserID(txCtx, req.AgencyID, req.UserID)
		if err != nil {
			return err
		}

		if req.Method == domain.MethodManual {
			if req.RequesterRole != domain.RoleAdmin {
				return domain.ErrManualNotAllowed
			}
		}

		if req.IsRemote != nil && *req.IsRemote {
			if user.HomeLatitude == nil || user.HomeLongitude == nil || user.HomeRadiusMeters == nil {
				return domain.ErrHomeLocationNotSet
			}
			if req.Latitude == nil || req.Longitude == nil {
				return domain.ErrInvalidAttendance
			}

			dist := utils.Haversine(*user.HomeLatitude, *user.HomeLongitude, *req.Latitude, *req.Longitude)
			if dist > float64(*user.HomeRadiusMeters) {
				return domain.ErrGeofenceViolation
			}
		}

		if req.Type == domain.TypeIn {
			if existing != nil {
				return domain.ErrAttendanceExists
			}

			entryTime := parseTimeMinutes(now, sched.EntryTimeMinutes)
			exitTime := parseTimeMinutes(now, sched.ExitTimeMinutes)

			status := domain.StatusPresent
			lateLimit := entryTime.Add(time.Duration(sched.GracePeriodMinutes) * time.Minute)
			if now.After(lateLimit) {
				status = domain.StatusLate
			}

			attendance := &domain.Attendance{
				UserID:            req.UserID,
				AgencyID:          req.AgencyID,
				CheckInTime:       now,
				ScheduleEntryTime: entryTime,
				ScheduleExitTime:  exitTime,
				Date:              now,
				Status:            status,
				MethodIn:          req.Method,
				Notes:             req.Notes,
				Latitude:          req.Latitude,
				Longitude:         req.Longitude,
			}

			if err := s.attendanceRepo.Create(txCtx, attendance); err != nil {
				return err
			}
			response = dto.ToAttendanceResponse(attendance)
			return nil
		}

		if existing == nil {
			return domain.ErrAttendanceNotFound
		}

		if existing.CheckOutTime != nil {
			return domain.ErrAttendanceExists
		}

		existing.CheckOutTime = &now
		existing.MethodOut = &req.Method

		if err := s.attendanceRepo.Update(txCtx, existing); err != nil {
			return err
		}
		response = dto.ToAttendanceResponse(existing)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *AttendanceService) GetAgencyAttendances(ctx context.Context, agencyID uuid.UUID, params *dto.AttendanceListParams) ([]*dto.AttendanceResponse, error) {
	filter := domain.AttendanceFilter{
		Page:   params.Page,
		Limit:  params.Limit,
		Status: domain.AttendanceStatus(params.Status),
	}

	if params.UserID != "" {
		if id, err := uuid.Parse(params.UserID); err == nil {
			filter.UserID = id
		}
	}

	if params.StartDate != "" {
		if t, err := time.Parse("2006-01-02", params.StartDate); err == nil {
			filter.StartDate = &t
		}
	}
	if params.EndDate != "" {
		if t, err := time.Parse("2006-01-02", params.EndDate); err == nil {
			filter.EndDate = &t
		}
	}

	attendances, err := s.attendanceRepo.List(ctx, agencyID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.AttendanceResponse, len(attendances))
	for i, a := range attendances {
		responses[i] = dto.ToAttendanceResponse(a)
	}
	return responses, nil
}

// Helper to set hours/minutes on a base date
func parseTimeMinutes(base time.Time, minutes int) time.Time {
	hours := minutes / 60
	mins := minutes % 60
	return time.Date(base.Year(), base.Month(), base.Day(), hours, mins, 0, 0, base.Location())
}
