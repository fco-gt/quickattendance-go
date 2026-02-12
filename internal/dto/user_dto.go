package dto

import (
	"quickattendance-go/internal/domain"
	"time"

	"github.com/google/uuid"
)

type InviteUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
}

type ActivateUserRequest struct {
	ActivationToken string          `json:"activation_token" binding:"required"`
	Password        string          `json:"password" binding:"required,min=8"`
	Profile         ActivateProfile `json:"profile" binding:"required"`
}

type ActivateProfile struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	FirstName        *string  `json:"first_name" binding:"required"`
	LastName         *string  `json:"last_name"`
	Email            *string  `json:"email" binding:"required,email"`
	HomeLatitude     *float64 `json:"home_latitude"`
	HomeLongitude    *float64 `json:"home_longitude"`
	HomeRadiusMeters *int     `json:"home_radius_meters"`
}

type UserResponse struct {
	ID               uuid.UUID     `json:"id"`
	FirstName        string        `json:"first_name"`
	LastName         *string       `json:"last_name"`
	Email            string        `json:"email"`
	Status           domain.Status `json:"status"`
	AgencyID         uuid.UUID     `json:"agency_id"`
	HomeLatitude     *float64      `json:"home_latitude"`
	HomeLongitude    *float64      `json:"home_longitude"`
	HomeRadiusMeters *int          `json:"home_radius_meters"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func ToUserResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:               user.ID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Email:            user.Email,
		Status:           user.Status,
		AgencyID:         user.AgencyID,
		HomeLatitude:     user.HomeLatitude,
		HomeLongitude:    user.HomeLongitude,
		HomeRadiusMeters: user.HomeRadiusMeters,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}
}

func ToAuthResponse(user *domain.User, token string) *AuthResponse {
	if user == nil {
		return nil
	}
	if token == "" {
		return nil
	}

	return &AuthResponse{
		Token: token,
		User:  *ToUserResponse(user),
	}
}

type PaginationParams struct {
	Page  int `form:"page,default=1" binding:"omitempty,min=1"`
	Limit int `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

type UserListParams struct {
	PaginationParams
	Status string `form:"status" binding:"omitempty"`
	Search string `form:"search" binding:"omitempty"`
}
