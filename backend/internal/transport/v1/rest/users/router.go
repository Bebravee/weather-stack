package usercntrl

import (
	"github.com/chup1x/weather-stack/internal/repository/postgres"
	"github.com/chup1x/weather-stack/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router fiber.Router, db *gorm.DB) {
	userRepo := postgres.NewUserRepository(db)

	userCntrl := NewUserController(services.NewUserService(userRepo))
	auth := router.Group("/auth")

	auth.Post("/profile/register", userCntrl.RegisterHandler)
	auth.Post("/login", userCntrl.LoginHandler)

	router.Get("/profile/by-id", userCntrl.GetProfileHandler)
	router.Get("/profile/by-telegram-id/:telegram_id", userCntrl.GetTelegramProfileHandler)
	router.Patch("/profile", userCntrl.UpdateProfileHandler)
}
