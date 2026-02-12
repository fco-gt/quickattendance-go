package repository

import (
	"autoattendance-go/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendanceRepo struct {
	db *gorm.DB
}

func NewAttendanceRepo(db *gorm.DB) *AttendanceRepo {
	return &AttendanceRepo{db: db}
}

func (r *AttendanceRepo) Create(ctx context.Context, attendance *domain.Attendance) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}
	return db.WithContext(ctx).Create(attendance).Error
}

func (r *AttendanceRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Attendance, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var attendance domain.Attendance
	if err := db.WithContext(ctx).First(&attendance, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAttendanceNotFound
		}
		return nil, err
	}
	return &attendance, nil
}

func (r *AttendanceRepo) GetTodayByUserID(ctx context.Context, agencyID uuid.UUID, userID uuid.UUID) (*domain.Attendance, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var attendance domain.Attendance
	// Considerar zona horaria si fuera necesario, de momento UTC/Local del servidor
	today := time.Now().Format("2006-01-02")

	err := db.WithContext(ctx).
		Where("agency_id = ? AND user_id = ? AND date = ?", agencyID, userID, today).
		First(&attendance).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &attendance, nil
}

func (r *AttendanceRepo) List(ctx context.Context, agencyID uuid.UUID, filter domain.AttendanceFilter) ([]*domain.Attendance, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var attendances []*domain.Attendance
	query := db.WithContext(ctx).Where("agency_id = ?", agencyID)

	if filter.UserID != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.StartDate != nil {
		query = query.Where("date >= ?", filter.StartDate.Format("2006-01-02"))
	}

	if filter.EndDate != nil {
		query = query.Where("date <= ?", filter.EndDate.Format("2006-01-02"))
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Pagination
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	if err := query.Find(&attendances).Error; err != nil {
		return nil, err
	}

	return attendances, nil
}

func (r *AttendanceRepo) Update(ctx context.Context, attendance *domain.Attendance) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}
	return db.WithContext(ctx).Save(attendance).Error
}

func (r *AttendanceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}
	return db.WithContext(ctx).Delete(&domain.Attendance{}, id).Error
}
