package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidUser           = errors.New("invalid user")
	ErrUserExists            = errors.New("user already exists")
	ErrUserAlreadyActivated  = errors.New("user already activated")
	ErrActivationCodeExpired = errors.New("activation code expired")
	ErrInvalidActivationCode = errors.New("invalid activation code")
	ErrUserNotActive         = errors.New("user not active")
)

// Utilizar punteros en Home Latitude, Home Longitude y Home Radius Meters para aceptar valores nulos al igual que para CodeExpiry
type User struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName        string    `gorm:"not null"`
	LastName         *string
	Email            string    `gorm:"uniqueIndex;not null"`
	PasswordHash     string    `gorm:"not null"`
	Status           Status    `gorm:"not null"`
	AgencyID         uuid.UUID `gorm:"type:uuid;not null"`
	Agency           Agency    `gorm:"foreignKey:AgencyID"`
	Role             Role      `gorm:"not null;default:'employee'"`
	HomeLatitude     *float64
	HomeLongitude    *float64
	HomeRadiusMeters *int
	ActivationCode   *string `gorm:"uniqueIndex"`
	CodeExpiry       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Status string

const (
	StatusPending  Status = "pending"
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByActivationCode(ctx context.Context, code string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByAgencyID(ctx context.Context, agencyID uuid.UUID) ([]*User, error)
}
