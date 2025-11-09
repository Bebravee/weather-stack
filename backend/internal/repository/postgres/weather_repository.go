package postgres

import (
	"context"

	"github.com/chup1x/weather-stack/internal/domain"
	"gorm.io/gorm"
)

type WetherRepository struct {
	db *gorm.DB
}

func NewWetherRepository(db *gorm.DB) *WetherRepository {
	return &WetherRepository{db: db}
}

func (r *WetherRepository) CreateWeatherRequest(ctx context.Context, new *domain.WeatherEntity) error {
	return r.db.WithContext(ctx).Table("weather").Create(new).Error
}

func (r *WetherRepository) GetWeatherByCity(ctx context.Context, city string) ([]*domain.WeatherEntity, error) {
	weather := []*domain.WeatherEntity{}

	if err := r.db.WithContext(ctx).Table("weather").Where("city_id = ?", city).First(weather).Error; err != nil {
		return nil, err
	}

	return weather, nil
}

func (r *WetherRepository) GetClothesByComb(ctx context.Context, id int) ([]*domain.ClothesEntity, error) {
	clothes := []*domain.ClothesEntity{}

	if err := r.db.WithContext(ctx).Table("clothes").Where("id = ?", id).First(clothes).Error; err != nil {
		return nil, err
	}

	return clothes, nil
}

func (r *WetherRepository) GetNewsByCity(ctx context.Context, city string) ([]*domain.WeatherEntity, error) {
	news := []*domain.NewsEntity{}

	if err := r.db.WithContext(ctx).Table("news").Where("city_id = ?", city).First(news).Error; err != nil {
		return nil, err
	}

	return news, nil
}
