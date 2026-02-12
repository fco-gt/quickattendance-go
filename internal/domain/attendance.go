package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrAttendanceNotFound  = errors.New("attendance not found")
	ErrNoAttendancesExist  = errors.New("this agency has no attendances")
	ErrAttendanceExists    = errors.New("attendance already exists")
	ErrAttendanceInvalid   = errors.New("invalid attendance")
	ErrInvalidUserOrAgency = errors.New("invalid user or agency")
	ErrGeofenceViolation   = errors.New("location out of range")
	ErrManualNotAllowed    = errors.New("only admins can mark attendance manually")
	ErrInvalidAttendance   = errors.New("invalid attendance data")
)

type Attendance struct {
	ID                uuid.UUID         `gorm:"type:uuid;primaryKey"`
	UserID            uuid.UUID         `gorm:"type:uuid;not null;index"`
	User              User              `gorm:"foreignKey:UserID"`
	AgencyID          uuid.UUID         `gorm:"type:uuid;not null;index"`
	Agency            Agency            `gorm:"foreignKey:AgencyID"`
	CheckInTime       time.Time         `gorm:"not null"` // Si se crea al entrar, es obligatorio
	ScheduleEntryTime time.Time         `gorm:"not null"`
	Status            AttendanceStatus  `gorm:"not null"`
	CheckOutTime      *time.Time        // Puede ser NULL hasta que salgan
	ScheduleExitTime  time.Time         `gorm:"not null"`
	Date              time.Time         `gorm:"type:date;not null"`
	MethodIn          AttendanceMethod  `gorm:"not null"`
	MethodOut         *AttendanceMethod // Opcional hasta el checkout
	Notes             *string
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Geolocation
	Latitude  *float64
	Longitude *float64
}

type AttendanceStatus string

var (
	StatusPresent AttendanceStatus = "present"
	StatusAbsent  AttendanceStatus = "absent"
	StatusLate    AttendanceStatus = "late"
	StatusEarly   AttendanceStatus = "early"
)

type AttendanceMethod string

var (
	MethodManual   AttendanceMethod = "manual"
	MethodQR       AttendanceMethod = "qr"
	MethodTelework AttendanceMethod = "telework"
	MethodNFC      AttendanceMethod = "nfc"
)

type AttendanceType string

var (
	TypeIn  AttendanceType = "in"
	TypeOut AttendanceType = "out"
)

func (a *Attendance) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type AttendanceFilter struct {
	UserID    uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
	Status    AttendanceStatus
	Page      int
	Limit     int
}

type AttendanceRepo interface {
	Create(ctx context.Context, attendance *Attendance) error
	GetByID(ctx context.Context, id uuid.UUID) (*Attendance, error)
	GetTodayByUserID(ctx context.Context, agencyID uuid.UUID, userID uuid.UUID) (*Attendance, error)
	List(ctx context.Context, agencyID uuid.UUID, filter AttendanceFilter) ([]*Attendance, error)
	Update(ctx context.Context, attendance *Attendance) error
	Delete(ctx context.Context, id uuid.UUID) error
}
