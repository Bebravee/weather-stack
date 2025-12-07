package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/chup1x/weather-stack/internal/domain"
	"gorm.io/gorm"
)

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *newsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) CreateNewsRequest(ctx context.Context, new *domain.NewsEntity) error {
	return r.db.WithContext(ctx).Table("news").Create(new).Error
}

func (r *newsRepository) GetNewsByCityID(ctx context.Context, cityID string) (*domain.NewsEntity, error) {
	news := &domain.NewsEntity{}

	startOfDay := time.Now().Truncate(24 * time.Hour)

	if err := r.db.WithContext(ctx).
		Table("news").
		Where("city_id = ? AND created_at >= ?", cityID, startOfDay).
		Order("created_at DESC").
		First(news).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewsNotFound
		}
		return nil, err
	}

	return news, nil
}
