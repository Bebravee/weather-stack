package weathercntrl

import (
	"github.com/chup1x/weather-stack/internal/repository/postgres"
	weatherservice "github.com/chup1x/weather-stack/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterWeatherRoutes(router fiber.Router, db *gorm.DB) {
	weatherRepo := postgres.NewWetherRepository(db)
	weatherCntrl := NewWeatherController(weatherservice.NewWeatherService(weatherRepo))

	weather := router.Group("/weather")
	{
		//weather.Post("/", weatherCntrl.CreateWeatherRecordHandler)
		weather.Get("/weather", weatherCntrl.GetWeatherHandler)
		weather.Get("/weather/clothes", weatherCntrl.GetWeatherClothesHandler)

		weather.Get("/news", weatherCntrl.GetNewsHandler)
	}
}
