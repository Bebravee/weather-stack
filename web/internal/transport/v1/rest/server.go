package rest

import (
	"context"
	"fmt"

	"github.com/chup1x/weather-stack/internal/config"
	"github.com/chup1x/weather-stack/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start(_ context.Context, cfg *config.Config) error {
	s.app = fiber.New()

	api := s.app.Group("/")

	api.Static("/docs", "./docs")

	registerHandlers(api, cfg)

	if err := s.app.Listen(fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port)); err != nil {
		return fmt.Errorf("server start: unable to start web server")
	}

	return nil
}

func registerHandlers(router fiber.Router, cfg *config.Config) {
	handlers := handlers.New(cfg.Server.RemoteHost)

	profile := router.Group("/profile")
	profile.Post("/register", handlers.RegisterUserHandler)
	profile.Get("/:id", handlers.GetUserHandler)
	profile.Get("/by-telegram-id/:telegram_id", handlers.GetUserByTelegramHandler)

	weather := router.Group("/weather")
	weather.Get("/city/:city", handlers.GetWeatherByCityHandler)
	weather.Get("/by-telegram-id/:telegram_id", handlers.GetWeatherByTelegramHandler)
	weather.Get("/clothes/:telegram_id", handlers.GetClothesByTelegramHandler)

	news := router.Group("/news")
	news.Get("/city/:city", handlers.GetNewsByCityHandler)
	news.Get("/by-telegram-id/:telegram_id", handlers.GetNewsByTelegramHandler)
}
