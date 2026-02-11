package service

import (
	"autoattendance-go/internal/domain"
	"autoattendance-go/internal/dto"
	"context"
	"time"
)

type AttendanceService struct {
	attendanceRepo domain.AttendanceRepo
	userRepo       domain.UserRepo
	agencyRepo     domain.AgencyRepo
}

func NewAttendanceService(attendanceRepo domain.AttendanceRepo, userRepo domain.UserRepo, agencyRepo domain.AgencyRepo) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		userRepo:       userRepo,
		agencyRepo:     agencyRepo,
	}
}

func (s *AttendanceService) MarkAttendance(ctx context.Context, req *dto.MarkAttendanceRequest) (*dto.AttendanceResponse, error) {
	// Validate user and agency
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	agency, err := s.agencyRepo.GetByID(ctx, req.AgencyID)
	if err != nil {
		return nil, domain.ErrAgencyNotFound
	}

	if user == nil || agency == nil {
		return nil, domain.ErrInvalidUserOrAgency
	}

	today := time.Now()

	// Check for an schedule for today
	//  TODO

}
