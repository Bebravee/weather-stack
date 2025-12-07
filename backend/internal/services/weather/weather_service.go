package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/chup1x/weather-stack/internal/domain"
)

type weatherCreater interface {
	CreateWeatherRequest(context.Context, *domain.WeatherEntity) error
}

type weatherProvider interface {
	GetWeatherByCity(context.Context, string) (*domain.WeatherEntity, error)
}

type weatherStorage interface {
	weatherCreater
	weatherProvider
}

type WeatherService struct {
	repo   weatherStorage
	client *weatherClient
}

func NewWeatherService(repo weatherStorage, client *weatherClient) *WeatherService {
	return &WeatherService{
		repo:   repo,
		client: client,
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

func (s *WeatherService) GetWeatherWithCache(ctx context.Context, city string) (*domain.WeatherEntity, error) {
	weather, err := s.repo.GetWeatherByCity(ctx, city)
	if err == nil {
		return weather, nil
	}

	if !errors.Is(err, domain.ErrWeatherNotFound) || s.client == nil {
		return nil, err
	}

	fetched, fetchErr := s.client.FetchWeather(ctx, city)
	if fetchErr != nil {
		return nil, fetchErr
	}

	if err := s.repo.CreateWeatherRequest(ctx, fetched); err != nil {
		return nil, fmt.Errorf("save fetched weather: %w", err)
	}

	return fetched, nil
}
