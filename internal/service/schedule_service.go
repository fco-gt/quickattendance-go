package service

import (
	"context"
	"quickattendance-go/internal/domain"
	"quickattendance-go/internal/dto"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ScheduleService struct {
	scheduleRepo domain.ScheduleRepo
	userRepo     domain.UserRepo
	transactor   domain.Transactor
}

func NewScheduleService(scheduleRepo domain.ScheduleRepo, userRepo domain.UserRepo, transactor domain.Transactor) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
		transactor:   transactor,
	}
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, agencyID uuid.UUID, req *dto.CreateScheduleRequest) (*dto.ScheduleResponse, error) {
	defaultSchedule, err := s.scheduleRepo.GetDefault(ctx, agencyID)
	if err != nil {
		return nil, err
	}

	if req.IsDefault && defaultSchedule != nil {
		return nil, domain.ErrDefaultScheduleAlreadyExists
	}

	repeated, err := s.scheduleRepo.GetByName(ctx, agencyID, req.Name)
	if err != nil {
		return nil, err
	}

	if len(repeated) > 0 {
		return nil, domain.ErrScheduleNameAlreadyExists
	}

	var daysStr []string
	for _, day := range req.DaysOfWeek {
		daysStr = append(daysStr, strconv.Itoa(day))
	}

	schedule := &domain.Schedule{
		Name:               req.Name,
		DaysOfWeek:         strings.Join(daysStr, ","),
		EntryTimeMinutes:   req.EntryTimeMinutes,
		ExitTimeMinutes:    req.ExitTimeMinutes,
		GracePeriodMinutes: req.GracePeriodMinutes,
		IsDefault:          req.IsDefault,
		AgencyID:           agencyID,
	}

	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	return dto.ToScheduleResponse(schedule), nil
}

func (s *ScheduleService) GetAgencySchedules(ctx context.Context, agencyID uuid.UUID, params *dto.ScheduleListParams) ([]*dto.ScheduleResponse, error) {
	filter := domain.ScheduleFilter{
		Name:      params.Name,
		IsDefault: params.IsDefault,
		Page:      params.Page,
		Limit:     params.Limit,
	}

	schedules, err := s.scheduleRepo.GetByAgencyID(ctx, agencyID, filter)
	if err != nil {
		return nil, err
	}

	var scheduleResponses []*dto.ScheduleResponse
	for _, schedule := range schedules {
		scheduleResponses = append(scheduleResponses, dto.ToScheduleResponse(schedule))
	}

	return scheduleResponses, nil
}

func (s *ScheduleService) GetSchedule(ctx context.Context, scheduleID uuid.UUID, agencyID uuid.UUID) (*dto.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return nil, err
	}

	if schedule.AgencyID != agencyID {
		return nil, domain.ErrScheduleNotFound
	}

	return dto.ToScheduleResponse(schedule), nil
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, req *dto.UpdateScheduleRequest, scheduleID uuid.UUID, agencyID uuid.UUID) (*dto.ScheduleResponse, error) {
	var response *dto.ScheduleResponse
	err := s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		schedule, err := s.scheduleRepo.GetByID(txCtx, scheduleID)
		if err != nil {
			return err
		}
		if schedule.AgencyID != agencyID {
			return domain.ErrScheduleNotFound
		}

		if req.Name != nil {
			repeated, _ := s.scheduleRepo.GetByName(txCtx, agencyID, *req.Name)
			for _, r := range repeated {
				if r.ID != scheduleID {
					return domain.ErrScheduleNameAlreadyExists
				}
			}
			schedule.Name = *req.Name
		}

		if req.IsDefault != nil && *req.IsDefault && !schedule.IsDefault {
			currentDefault, _ := s.scheduleRepo.GetDefault(txCtx, agencyID)
			if currentDefault != nil {
				currentDefault.IsDefault = false
				if err := s.scheduleRepo.Update(txCtx, currentDefault); err != nil {
					return err
				}
			}
			schedule.IsDefault = true
		}

		if req.DaysOfWeek != nil {
			var days []string
			for _, d := range *req.DaysOfWeek {
				days = append(days, strconv.Itoa(d))
			}
			schedule.DaysOfWeek = strings.Join(days, ",")
		}

		if req.EntryTimeMinutes != nil {
			schedule.EntryTimeMinutes = *req.EntryTimeMinutes
		}
		if req.ExitTimeMinutes != nil {
			schedule.ExitTimeMinutes = *req.ExitTimeMinutes
		}
		if req.GracePeriodMinutes != nil {
			schedule.GracePeriodMinutes = *req.GracePeriodMinutes
		}

		if req.AssignedUsersIDs != nil {
			var users []domain.User
			for _, id := range *req.AssignedUsersIDs {
				u, err := s.userRepo.GetByID(txCtx, id)
				if err != nil {
					return err
				}

				if u.AgencyID != agencyID {
					return domain.ErrUserNotFound
				}
				users = append(users, *u)
			}
			schedule.AssignedUsers = users
		}

		if err := s.scheduleRepo.Update(txCtx, schedule); err != nil {
			return err
		}

		response = dto.ToScheduleResponse(schedule)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, scheduleID uuid.UUID, agencyID uuid.UUID) error {
	schedule, err := s.scheduleRepo.GetByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	if schedule.AgencyID != agencyID {
		return domain.ErrScheduleNotFound
	}

	if schedule.IsDefault {
		return domain.ErrDeleteDefaultSchedule
	}

	if err := s.scheduleRepo.Delete(ctx, scheduleID); err != nil {
		return err
	}

	return nil
}

func (s *ScheduleService) GetApplicableSchedule(ctx context.Context, agencyID uuid.UUID, userID uuid.UUID, date time.Time) (*dto.ScheduleResponse, error) {
	weekday := strconv.Itoa(int(date.Weekday()))

	userSchedule, _ := s.scheduleRepo.GetUserScheduleByDay(ctx, agencyID, userID, weekday)
	if userSchedule != nil {
		return dto.ToScheduleResponse(userSchedule), nil
	}

	defaultSchedule, _ := s.scheduleRepo.GetDefault(ctx, agencyID)
	if defaultSchedule != nil {
		if strings.Contains(defaultSchedule.DaysOfWeek, weekday) {
			return dto.ToScheduleResponse(defaultSchedule), nil
		}
	}

	return nil, domain.ErrNoScheduleFound
}
