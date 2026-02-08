package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrAgencyNotFound     = errors.New("agency not found")
	ErrInvalidAgency      = errors.New("invalid agency")
	ErrAgencyExists       = errors.New("agency already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Agency struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"uniqueIndex;not null"`
	Domain    string    `gorm:"uniqueIndex;not null"`
	Address   string
	Phone     string
	IsActive  bool   `gorm:"default:true"`
	Users     []User `gorm:"foreignKey:AgencyID;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Agency) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type AgencyRepo interface {
	Create(ctx context.Context, agency *Agency) error
	GetByID(ctx context.Context, id uuid.UUID) (*Agency, error)
	GetByName(ctx context.Context, name string) (*Agency, error)
	GetByDomain(ctx context.Context, domain string) (*Agency, error)
	Update(ctx context.Context, agency *Agency) error
}
