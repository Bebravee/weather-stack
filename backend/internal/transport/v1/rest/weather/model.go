package weathercntrl

import (
	"github.com/chup1x/weather-stack/internal/domain"
)

type GetWeatherHistoryRequest struct {
	City_w string `param:"city_w" validate:"required"`
}

type GetWeatherHistoryResponse struct {
	*domain.WeatherEntity
}

type GetWeatherClothesRequest struct {
	user int64 `json:"telegram_id" validate:"required"`
}

type GetWeatherClothesResponse struct {
	*domain.WeatherClothesEntity
}

type GetNewsRequest struct {
	city_n string `json:"city_n" validate:"required"`
}

type GetNewsResponse struct {
	*domain.NewsEntity
}

type CreateWeatherRecord struct {
	TelegramID  int64   `json:"telegram_id"`
	Temperature float64 `json:"temperature"`
	Clothing    string  `json:"clothing"`
}
