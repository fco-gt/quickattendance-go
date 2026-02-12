package repository

import (
	"context"
	"quickattendance-go/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}
	return db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}

	var user domain.User
	if err := db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByActivationCode(ctx context.Context, code string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("activation_code = ?", code).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}
	return db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db = r.db
	}
	return db.WithContext(ctx).Delete(&domain.User{}, id).Error
}

func (r *UserRepo) ListByAgencyID(ctx context.Context, agencyID uuid.UUID, filter domain.UserFilter) ([]*domain.User, error) {
	var users []*domain.User
	query := r.db.WithContext(ctx).Where("agency_id = ?", agencyID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Pagination
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
