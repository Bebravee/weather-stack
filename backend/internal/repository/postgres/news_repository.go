package postgres

import (
	"context"
	"errors"

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

	if err := r.db.WithContext(ctx).Table("news").Where("city_id = ?", cityID).First(news).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewsNotFound
		}
		return nil, err
	}
		// if err := r.db.WithContext(ctx).Table("news").Create(news_en).Error; err != nil {
		// 	log.Fatal("Error writing file to database:", err)
		// 	return
		// }
		// return news, nil

	return news, nil
}
