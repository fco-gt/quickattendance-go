package repository

import (
	"autoattendance-go/internal/domain"
	"context"

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
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}

func (r *UserRepo) ListByAgencyID(ctx context.Context, agencyID uuid.UUID) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.WithContext(ctx).Where("agency_id = ?", agencyID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
