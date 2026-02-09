package dto

import (
	"autoattendance-go/internal/domain"
	"time"

	"github.com/google/uuid"
)

type RegisterAgencyRequest struct {
	Name       string `json:"name" binding:"required"`
	Domain     string `json:"domain" binding:"required"`
	Address    string `json:"address"`
	Phone      string `json:"phone"`
	AdminEmail string `json:"admin_email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
}

type UpdateAgencyRequest struct {
	Name    *string `json:"name"`
	Address *string `json:"address"`
	Phone   *string `json:"phone"`
}

type AgencyResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToAgencyResponse(agency *domain.Agency) *AgencyResponse {
	if agency == nil {
		return nil
	}

	return &AgencyResponse{
		ID:        agency.ID,
		Name:      agency.Name,
		Domain:    agency.Domain,
		Address:   agency.Address,
		Phone:     agency.Phone,
		IsActive:  agency.IsActive,
		CreatedAt: agency.CreatedAt,
		UpdatedAt: agency.UpdatedAt,
	}
}
