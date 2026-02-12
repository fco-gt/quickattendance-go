package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrScheduleNotFound             = errors.New("schedule not found")
	ErrNoScheduleFound              = errors.New("no schedule found")
	ErrDefaultScheduleNotFound      = errors.New("default schedule not found")
	ErrDefaultScheduleAlreadyExists = errors.New("default schedule already exists")
	ErrScheduleNameAlreadyExists    = errors.New("schedule name already exists")
	ErrDeleteDefaultSchedule        = errors.New("cannot delete the default schedule of an agency")
)

type Schedule struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	AgencyID uuid.UUID `gorm:"type:uuid;not null;index"`
	Agency   Agency    `gorm:"foreignKey:AgencyID"`
	Name     string    `gorm:"not null"`
	// Recomendación: String para evitar líos de drivers con arrays
	DaysOfWeek         string `gorm:"not null"` // Ej: "1,2,3,4,5,"
	EntryTimeMinutes   int    `gorm:"not null"`
	ExitTimeMinutes    int    `gorm:"not null"`
	GracePeriodMinutes int    `gorm:"not null"`
	IsDefault          bool   `gorm:"not null"`
	AssignedUsers      []User `gorm:"many2many:schedule_users;"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type ScheduleFilter struct {
	Name      string
	IsDefault *bool
	Page      int
	Limit     int
}

type ScheduleRepo interface {
	Create(ctx context.Context, schedule *Schedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*Schedule, error)
	GetByAgencyID(ctx context.Context, agencyID uuid.UUID, filter ScheduleFilter) ([]*Schedule, error)
	GetByDate(ctx context.Context, agencyID uuid.UUID, date time.Time) ([]*Schedule, error)
	GetDefault(ctx context.Context, agencyID uuid.UUID) (*Schedule, error)
	GetByName(ctx context.Context, agencyID uuid.UUID, name string) ([]*Schedule, error)
	GetUserScheduleByDay(ctx context.Context, agencyID uuid.UUID, userID uuid.UUID, weekday string) (*Schedule, error)
	Update(ctx context.Context, schedule *Schedule) error
	Delete(ctx context.Context, id uuid.UUID) error
}
