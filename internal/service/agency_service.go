package service

import (
	"autoattendance-go/internal/domain"
	"autoattendance-go/internal/dto"
	"autoattendance-go/pkg/security"
	"context"

	"github.com/google/uuid"
)

type AgencyService struct {
	agencyRepo domain.AgencyRepo
	userRepo   domain.UserRepo
	hasher     *security.PasswordHasher
	txManager  domain.Transactor
}

func NewAgencyService(agencyRepo domain.AgencyRepo, userRepo domain.UserRepo, hasher *security.PasswordHasher, txManager domain.Transactor) *AgencyService {
	return &AgencyService{
		agencyRepo: agencyRepo,
		userRepo:   userRepo,
		hasher:     hasher,
		txManager:  txManager,
	}
}

func (s *AgencyService) Register(ctx context.Context, req *dto.RegisterAgencyRequest) (*dto.AgencyResponse, error) {
	// 1. Validaciones previas de existencia
	if _, err := s.agencyRepo.GetByName(ctx, req.Name); err == nil {
		return nil, domain.ErrAgencyExists
	}
	if _, err := s.agencyRepo.GetByDomain(ctx, req.Domain); err == nil {
		return nil, domain.ErrAgencyExists
	}
	if _, err := s.userRepo.GetByEmail(ctx, req.AdminEmail); err == nil {
		return nil, domain.ErrUserExists
	}

	var agency *domain.Agency
	err := s.txManager.WithinTransaction(ctx, func(tCtx context.Context) error {
		agency = &domain.Agency{
			Name:    req.Name,
			Domain:  req.Domain,
			Address: req.Address,
			Phone:   req.Phone,
		}
		if err := s.agencyRepo.Create(tCtx, agency); err != nil {
			return err
		}

		hashedPassword, _ := s.hasher.Hash(tCtx, req.Password)

		adminUser := &domain.User{
			FirstName:    "Admin",
			Email:        req.AdminEmail,
			Role:         domain.RoleAdmin,
			PasswordHash: hashedPassword,
			AgencyID:     agency.ID,
			Status:       domain.StatusActive,
		}

		// En caso de que falle la creación del usuario, GORM hace ROLLBACK de la agencia automáticamente
		if err := s.userRepo.Create(tCtx, adminUser); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.ToAgencyResponse(agency), nil
}

func (s *AgencyService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateAgencyRequest) (*dto.AgencyResponse, error) {
	agency, err := s.agencyRepo.GetByID(ctx, id)

	if err != nil {
		return nil, domain.ErrAgencyNotFound
	}

	if req.Name != nil {
		agency.Name = *req.Name
	}
	if req.Address != nil {
		agency.Address = *req.Address
	}
	if req.Phone != nil {
		agency.Phone = *req.Phone
	}

	if err := s.agencyRepo.Update(ctx, agency); err != nil {
		return nil, err
	}

	return dto.ToAgencyResponse(agency), nil
}
