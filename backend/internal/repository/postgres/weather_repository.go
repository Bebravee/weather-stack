package postgres

import (
	"context"

	"github.com/chup1x/weather-stack/internal/domain"
	"gorm.io/gorm"
)

type WeatherRepository struct {
	db *gorm.DB
}

func NewWeatherRepository(db *gorm.DB) *WeatherRepository {
	return &WeatherRepository{db: db}
}

func (r *WeatherRepository) CreateWeatherRequest(ctx context.Context, new *domain.WeatherEntity) error {
	return r.db.WithContext(ctx).Table("weather").Create(new).Error
}

func (r *WeatherRepository) GetWeatherByCity(ctx context.Context, city string) ([]*domain.WeatherEntity, error) {
	weather := []*domain.WeatherEntity{}

	if err := r.db.WithContext(ctx).Table("weather").Where("city_id = ?", city).First(weather).Error; err != nil {
		return nil, err
	}

	return weather, nil
}

func (r *WeatherRepository) GetClothesByComb(ctx context.Context, id int) ([]*domain.WeatherClothesEntity, error) {
	clothes := []*domain.WeatherClothesEntity{}

	if err := r.db.WithContext(ctx).Table("clothes").Where("id = ?", id).First(clothes).Error; err != nil {
		return nil, err
	}

	return clothes, nil
}

func (r *WeatherRepository) GetNewsByCity(ctx context.Context, city string) ([]*domain.NewsEntity, error) {
	news := []*domain.NewsEntity{}

	if err := r.db.WithContext(ctx).Table("news").Where("city_id = ?", city).First(news).Error; err != nil {
		return nil, err
	}

	return news, nil
}
