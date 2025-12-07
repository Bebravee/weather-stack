package rest

import (
	"context"
	"fmt"

	"github.com/chup1x/weather-stack/internal/config"
	newscntrl "github.com/chup1x/weather-stack/internal/transport/v1/rest/news"
	usercntrl "github.com/chup1x/weather-stack/internal/transport/v1/rest/users"
	weathercntrl "github.com/chup1x/weather-stack/internal/transport/v1/rest/weather"
	"github.com/chup1x/weather-stack/pkg/database"
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

	api := s.app.Group("/api/v1")

	s.app.Static("/docs", "./docs")

	db, err := database.ConnectPostgres(database.PostgresConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		DBName:   cfg.Postgres.Name,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		return fmt.Errorf("to connect to postgres: %w", err)
	}

	if err := database.PostgresMigrations(db); err != nil {
		return fmt.Errorf("to apply migrations: %w", err)
	}

	usercntrl.RegisterUserRoutes(api, db)
	weathercntrl.RegisterWeatherRoutes(api, cfg, db)
	newscntrl.RegisterNewsRoutes(api, cfg, db)

	if err := s.app.Listen(fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port)); err != nil {
		return fmt.Errorf("server start: unable to start web server")
	}

	return nil
}
