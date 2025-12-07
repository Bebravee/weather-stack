package weathercntrl

import (
	"github.com/chup1x/weather-stack/internal/config"
	"github.com/chup1x/weather-stack/internal/repository/postgres"
	weatherservice "github.com/chup1x/weather-stack/internal/services/weather"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterWeatherRoutes(router fiber.Router, cfg *config.Config, db *gorm.DB) {
	weatherRepo := postgres.NewWeatherRepository(db)
	userRepo := postgres.NewUserRepository(db)
	weatherClient := weatherservice.NewWeatherClient(cfg.Weather.Host, cfg.Weather.APIKey, cfg.Weather.Lang, cfg.Weather.Units)
	weatherService := weatherservice.NewWeatherService(weatherRepo, weatherClient)

	clothesRepo := postgres.NewClothesRepository(db)
	llmClient := weatherservice.NewLLMClient(cfg.LLM)
	clothesService := weatherservice.NewClothesService(clothesRepo, userRepo, weatherService, llmClient)

	weatherCntrl := NewWeatherController(weatherService, clothesService, userRepo)

	weather := router.Group("/weather")

	weather.Get("/:id", weatherCntrl.GetWeatherHandler)
	weather.Get("/city/:id", weatherCntrl.GetWeatherHandler)
	weather.Get("/by-telegram-id/:telegram_id", weatherCntrl.GetWeatherByTelegramHandler)
	weather.Get("/clothes/:telegram_id", weatherCntrl.GetWeatherClothesHandler)
	weather.Post("/", weatherCntrl.CreateWeatherRecordHandler)
}
