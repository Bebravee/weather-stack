package weathercntrl

import (
	"github.com/chup1x/weather-stack/internal/domain"
)

type GetWeatherHistoryRequest struct {
	ID string `param:"id" validate:"required"`
}

type GetWeatherByTelegramRequest struct {
	TelegramID int `params:"telegram_id" validate:"required"`
}

type GetWeatherHistoryResponse struct {
	*domain.WeatherEntity
}

type GetWeatherClothesRequest struct {
	TelegramID int `params:"telegram_id" validate:"required"`
}

type CreateWeatherRequest struct {
	*domain.WeatherEntity
}

type CreateWeatherResponse struct {
	City string `json:"city"`
}
