package repository

import (
	"context"
	"quickattendance-go/internal/domain"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
	return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) Create(ctx context.Context, schedule *domain.Schedule) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}
	return db.WithContext(ctx).Create(schedule).Error
}

func (r *ScheduleRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Schedule, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var schedule domain.Schedule
	// Pre carga de usuarios asignados
	if err := db.WithContext(ctx).Preload("AssignedUsers").First(&schedule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrScheduleNotFound
		}
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepo) GetByAgencyID(ctx context.Context, agencyID uuid.UUID, filter domain.ScheduleFilter) ([]*domain.Schedule, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var schedules []*domain.Schedule
	query := db.WithContext(ctx).Where("agency_id = ?", agencyID)

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.IsDefault != nil {
		query = query.Where("is_default = ?", *filter.IsDefault)
	}

	// Pagination
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	if err := query.Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *ScheduleRepo) GetByName(ctx context.Context, agencyID uuid.UUID, name string) ([]*domain.Schedule, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var schedules []*domain.Schedule
	if err := db.WithContext(ctx).Where("agency_id = ? AND name = ?", agencyID, name).Find(&schedules).Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

// GetByDate busca horarios que apliquen a un día específico (0=Dom, 1=Lun...)
func (r *ScheduleRepo) GetByDate(ctx context.Context, agencyID uuid.UUID, date time.Time) ([]*domain.Schedule, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	weekday := strconv.Itoa(int(date.Weekday()))

	var schedules []*domain.Schedule
	err := db.WithContext(ctx).
		Where("agency_id = ? AND days_of_week LIKE ?", agencyID, "%"+weekday+"%").
		Find(&schedules).Error

	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *ScheduleRepo) Update(ctx context.Context, schedule *domain.Schedule) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	return db.WithContext(ctx).Save(schedule).Error
}

func (r *ScheduleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	return db.WithContext(ctx).Delete(&domain.Schedule{}, id).Error
}

func (r *ScheduleRepo) GetDefault(ctx context.Context, agencyID uuid.UUID) (*domain.Schedule, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var schedules []domain.Schedule
	err := db.WithContext(ctx).
		Where("agency_id = ? AND is_default = ?", agencyID, true).
		Limit(1).
		Find(&schedules).Error

	if err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return nil, nil
	}

	return &schedules[0], nil
}

func (r *ScheduleRepo) GetUserScheduleByDay(ctx context.Context, agencyID uuid.UUID, userID uuid.UUID, weekday string) (*domain.Schedule, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var schedule domain.Schedule
	err := db.WithContext(ctx).
		Joins("JOIN schedule_users ON schedule_users.schedule_id = schedules.id").
		Where("schedules.agency_id = ? AND schedule_users.user_id = ? AND schedules.days_of_week LIKE ?",
			agencyID, userID, "%"+weekday+"%").
		First(&schedule).Error

	if err != nil {
		return nil, err
	}
	return &schedule, nil
}
