package weathercntrl

import (
	"github.com/chup1x/weather-stack/internal/config"
	"github.com/chup1x/weather-stack/internal/repository/postgres"
	newsservice "github.com/chup1x/weather-stack/internal/services/news"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterNewsRoutes(router fiber.Router, cfg *config.Config, db *gorm.DB) {
	newsRepo := postgres.NewNewsRepository(db)
	newsClient := newsservice.NewNewsClient(cfg.News.Host, cfg.News.APIKey)
	userRepo := postgres.NewUserRepository(db)
	newsCntrl := NewNewsController(newsservice.NewNewsService(newsRepo, newsClient), userRepo)

	news := router.Group("/news")
	news.Get("/:city_id", newsCntrl.GetNewsHandler)
	news.Get("/city/:city_id", newsCntrl.GetNewsHandler)
	news.Get("/by-telegram-id/:telegram_id", newsCntrl.GetNewsByTelegramHandler)
}
