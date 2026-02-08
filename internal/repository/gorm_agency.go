package repository

import (
	"autoattendance-go/internal/domain"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AgencyRepo struct {
	db *gorm.DB
}

func NewAgencyRepo(db *gorm.DB) domain.AgencyRepo {
	return &AgencyRepo{db: db}
}

func (r *AgencyRepo) Create(ctx context.Context, agency *domain.Agency) error {
	// Verificamos si hay alguna transaccion en el contexto
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		return tx.Create(agency).Error
	}

	return r.db.WithContext(ctx).Create(agency).Error
}

func (r *AgencyRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Agency, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var agency domain.Agency
	if err := db.WithContext(ctx).First(&agency, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAgencyNotFound
		}
		return nil, err
	}
	return &agency, nil
}

func (r *AgencyRepo) GetByName(ctx context.Context, name string) (*domain.Agency, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var agency domain.Agency
	if err := db.WithContext(ctx).Where("name = ?", name).First(&agency).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAgencyNotFound
		}
		return nil, err
	}
	return &agency, nil
}

func (r *AgencyRepo) GetByDomain(ctx context.Context, domainName string) (*domain.Agency, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var agency domain.Agency
	if err := db.WithContext(ctx).Where("domain = ?", domainName).First(&agency).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAgencyNotFound
		}
		return nil, err
	}
	return &agency, nil
}

func (r *AgencyRepo) Update(ctx context.Context, agency *domain.Agency) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}
	// Usamos Session para garantizar que todas las relaciones sean actualizadas
	return db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(agency).Error
}

func (r *AgencyRepo) GetUsers(ctx context.Context, agencyID uuid.UUID) ([]*domain.User, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var users []*domain.User
	if err := db.WithContext(ctx).Where("agency_id = ?", agencyID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
