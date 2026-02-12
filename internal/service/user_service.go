package service

import (
	"autoattendance-go/internal/domain"
	"autoattendance-go/internal/dto"
	"autoattendance-go/pkg/security"
	"context"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo   domain.UserRepo
	agencyRepo domain.AgencyRepo
	jwt        *security.JWTService
	hasher     *security.PasswordHasher
	tokenTTL   time.Duration
}

func NewUserService(userRepo domain.UserRepo, agencyRepo domain.AgencyRepo, jwt *security.JWTService, hasher *security.PasswordHasher, tokenTTL time.Duration) *UserService {
	return &UserService{
		userRepo:   userRepo,
		agencyRepo: agencyRepo,
		jwt:        jwt,
		hasher:     hasher,
		tokenTTL:   tokenTTL,
	}
}

func (s *UserService) Invite(ctx context.Context, agencyID uuid.UUID, req *dto.InviteUserRequest) error {
	// Verificar si el usuario ya existe
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return domain.ErrUserExists
	}

	activationToken, err := security.GenerateRandomToken(32)
	if err != nil {
		return err
	}

	expiry := time.Now().Add(time.Hour * 24)

	user := &domain.User{
		FirstName:      req.FirstName,
		LastName:       &req.LastName,
		Email:          req.Email,
		AgencyID:       agencyID,
		Role:           domain.RoleEmployee,
		Status:         domain.StatusPending,
		ActivationCode: &activationToken,
		CodeExpiry:     &expiry,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) ActivateByCode(ctx context.Context, req *dto.ActivateUserRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByActivationCode(ctx, req.ActivationToken)
	if err != nil {
		return nil, domain.ErrInvalidActivationCode
	}

	if user.Status != domain.StatusPending {
		return nil, domain.ErrUserAlreadyActivated
	}

	if user.CodeExpiry == nil || time.Now().After(*user.CodeExpiry) {
		return nil, domain.ErrActivationCodeExpired
	}

	hashedPassword, err := s.hasher.Hash(ctx, req.Password)
	if err != nil {
		return nil, err
	}

	user.Status = domain.StatusActive
	user.FirstName = req.Profile.FirstName
	user.LastName = &req.Profile.LastName
	user.PasswordHash = hashedPassword

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.jwt.Sign(user.ID, user.AgencyID, domain.Role(user.Role), s.tokenTTL)
	if err != nil {
		return nil, err
	}

	return dto.ToAuthResponse(user, token), nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginUserRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if user.Status != domain.StatusActive {
		return nil, domain.ErrUserNotActive
	}

	ok, err := s.hasher.Compare(ctx, user.PasswordHash, req.Password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !ok {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := s.jwt.Sign(user.ID, user.AgencyID, domain.Role(user.Role), s.tokenTTL)
	if err != nil {
		return nil, err
	}

	return dto.ToAuthResponse(user, token), nil
}

func (s *UserService) GetUserByID(ctx context.Context, agencyID uuid.UUID, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if user.AgencyID != agencyID {
		return nil, domain.ErrUserNotFound
	}

	return dto.ToUserResponse(user), nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, agencyID uuid.UUID, userID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if user.AgencyID != agencyID {
		return nil, domain.ErrUserNotFound
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = req.LastName
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.HomeLatitude != nil {
		user.HomeLatitude = req.HomeLatitude
	}

	if req.HomeLongitude != nil {
		user.HomeLongitude = req.HomeLongitude
	}

	if req.HomeRadiusMeters != nil {
		user.HomeRadiusMeters = req.HomeRadiusMeters
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return dto.ToUserResponse(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, agencyID uuid.UUID, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	if user.AgencyID != agencyID {
		return domain.ErrUserNotFound
	}

	if err := s.userRepo.Delete(ctx, user.ID); err != nil {
		return err
	}

	return nil
}

func (s *UserService) ListByAgencyID(ctx context.Context, agencyID uuid.UUID, params *dto.UserListParams) ([]*dto.UserResponse, error) {
	filter := domain.UserFilter{
		Status: params.Status,
		Search: params.Search,
		Page:   params.Page,
		Limit:  params.Limit,
	}

	users, err := s.userRepo.ListByAgencyID(ctx, agencyID, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = dto.ToUserResponse(user)
	}

	return responses, nil
}
