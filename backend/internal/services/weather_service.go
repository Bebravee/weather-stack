package services

import (
	"context"
	"fmt"

	"github.com/chup1x/weather-stack/internal/domain"
)

type weatherCreater interface {
	CreateWeatherRequest(context.Context, *domain.WeatherEntity) error
}

type weatherProvider interface {
	GetWeatherByCity(context.Context, string) (*domain.WeatherEntity, error)
	GetClothesByComb(context.Context, int) (*domain.WeatherClothesEntity, error)
	GetNewsByCity(context.Context, string) (*domain.NewsEntity, error)
}

type weatherStorage interface {
	weatherCreater
	weatherProvider
	userProvider
}

type WeatherService struct {
	repo weatherStorage
}

func NewWeatherService(repo weatherStorage) *WeatherService {
	return &WeatherService{
		repo: repo,
	}
}

func (s *WeatherService) CreateWeatherRecord(ctx context.Context, new *domain.WeatherEntity) error {
	if err := s.repo.CreateWeatherRequest(ctx, new); err != nil {
		return fmt.Errorf("to create a weatcher request: %w", err)
	}
	return nil
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (*domain.WeatherEntity, error) {
	weather, err := s.repo.GetWeatherByCity(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("to select a weather by city: %w", err)
	}
	return weather, nil
}

func (s *WeatherService) GetWeatherClothes(ctx context.Context, id int64) (*domain.WeatherClothesEntity, error) {
	//comb, err := s.repo.SelectByTelegramID(ctx, id)
	clothes, err := s.repo.GetClothesByComb(ctx, 1234)
	if err != nil {
		return nil, fmt.Errorf("to select a weather by city: %w", err)
	}
	return clothes, nil
}

func (s *WeatherService) GetNews(ctx context.Context, city string) (*domain.NewsEntity, error) {
	news, err := s.repo.GetNewsByCity(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("to select a weather by city: %w", err)
	}
	return news, nil
}
