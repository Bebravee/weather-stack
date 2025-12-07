package postgres

import (
	"context"
	"errors"
	"time"

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
	return r.db.WithContext(ctx).Table("weather_requests").Create(new).Error
}

func (r *WeatherRepository) GetWeatherByCity(ctx context.Context, city string) (*domain.WeatherEntity, error) {
	weather := &domain.WeatherEntity{}

	startOfDay := time.Now().Truncate(24 * time.Hour)

	if err := r.db.WithContext(ctx).
		Table("weather_requests").
		Where("city_id = ? AND created_at >= ?", city, startOfDay).
		Order("created_at DESC").
		First(weather).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWeatherNotFound
		}
		return nil, err
	}

	return weather, nil
}

/*
	func (r *WeatherRepository) GetClothesByComb(ctx context.Context, id int) ([]*domain.WeatherClothesEntity, error) {
		clothes := []*domain.WeatherClothesEntity{}

		if err := r.db.WithContext(ctx).Table("clothes").Where("id = ?", id).First(clothes).Error; err != nil {
			return nil, err
		}

		return clothes, nil
	}
*/
