package postgres

import (
	"context"
	"errors"

	"github.com/chup1x/weather-stack/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, new *domain.UserEntity) error {
	if err := r.db.WithContext(ctx).Table("users").Create(new).Error; err != nil {
		r.db.Migrator().AddColumn("users", "city_n")
	}
	return r.db.WithContext(ctx).Table("users").Create(new).Error
}

func (r *userRepository) SelectByID(ctx context.Context, id domain.UserID) (*domain.UserEntity, error) {
	user := &domain.UserEntity{}
	if err := r.db.WithContext(ctx).Table("users").Where("id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) SelectByTelegramID(ctx context.Context, id int64) (*domain.UserEntity, error) {
	user := &domain.UserEntity{}
	if err := r.db.WithContext(ctx).Table("users").Where("telegram_id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
