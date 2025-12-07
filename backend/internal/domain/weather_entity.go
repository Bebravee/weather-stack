package domain

import (
	"errors"
	"time"
)

var ErrWeatherNotFound = errors.New("weather not found")
var ErrClothesNotFound = errors.New("clothes not found")

type WeatherEntity struct {
	CityID      string    `json:"city" gorm:"city_id"`
	Temperature float64   `json:"temperature" gorm:"temperature"`
	FeelsLike   float64   `json:"feels" gorm:"feels_like"`
	Description string    `json:"description" gorm:"description"`
	Humidity    float64   `json:"humidity" gorm:"humidity"`
	Pressure    float64   `json:"pressure" gorm:"pressure"`
	WindSpeed   float64   `json:"wind_speed" gorm:"wind_speed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type WeatherClothesEntity struct {
	Code string `json:"code" gorm:"column:city_id"`
	Path string `json:"path" gorm:"column:path"`
}
