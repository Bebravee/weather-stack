package services

import (
	"context"
	"fmt"

	"github.com/chup1x/weather-stack/internal/domain"
)

type weatherCreater interface {
	CreateWeatherRequest(context.Context, *domain.WeatherRequestEntity) error
}

type weatherProvider interface {
	GetWeatherRequestByUserID(context.Context, domain.UserID) ([]*domain.WeatherRequestEntity, error)
}

type weatherStorage interface {
	weatherCreater
	weatherProvider
}

type WeatherService struct {
	repo weatherStorage
}

func NewWeatherService(repo weatherStorage) *WeatherService {
	return &WeatherService{
		repo: repo,
	}
}

func (s *WeatherService) CreateWeatherRecord(ctx context.Context, new *domain.WeatherRequestEntity) error {
	if err := s.repo.CreateWeatherRequest(ctx, new); err != nil {
		return fmt.Errorf("to create a weatcher request: %w", err)
	}
	return nil
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (domain.WeatherEntity, error) {
	weather, err := s.weather.SelectByCity(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("to select a weather by city: %w", err)
	}
	return weather, nil
}

func (s *WeatherService) GetWeatherClothes(ctx context.Context, id int) (domain.ClothesEntity, error) {
	comb, err := s.users.SelectByID(ctx, id)
	clothes, err := s.clothes.SelectByComb(ctx, comb.temp1)
	if err != nil {
		return nil, fmt.Errorf("to select a weather by city: %w", err)
	}
	return clothes, nil
}

func (s *WeatherService) GetNews(ctx context.Context, city string) (domain.NewsEntity, error) {
	news, err := s.news.SelectByCity(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("to select a weather by city: %w", err)
	}
	return news, nil
}
