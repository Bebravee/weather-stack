package domain

import (
	"time"
)

type WeatherEntity struct {
	CITY_ID     string    `json:"city_id" gorm:"city_id"`
	Temperature float64   `json:"temperature" gorm:"temperature"`
	feels_like  float64   `json:"feels" gorm:"feels"`
	Description string    `json:"description" gorm:"description"`
	Humidity    float64   `json:"humidity" gorm:"humidity"`
	Pressure    float64   `json:"pressure" gorm:"pressure"`
	Wind_speed  float64   `json:"wind_speed" gorm:"wind_speed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type WeatherClothesEntity struct {
	PATH 		string	  `json:"path" gorm:"path"`
}

type NewsEntity struct {
	PATH 		string	  `json:"path" gorm:"path"`
}
