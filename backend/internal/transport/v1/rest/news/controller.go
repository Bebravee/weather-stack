package weathercntrl

import (
	"context"
	"errors"

	"github.com/chup1x/weather-stack/internal/domain"
	newsservice "github.com/chup1x/weather-stack/internal/services/news"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type newsUserProvider interface {
	SelectByTelegramID(context.Context, int) (*domain.UserEntity, error)
}

type newsController struct {
	s         *newsservice.NewsService
	users     newsUserProvider
	validator *validator.Validate
}

func NewNewsController(service *newsservice.NewsService, users newsUserProvider) *newsController {
	return &newsController{
		validator: validator.New(),
		s:         service,
		users:     users,
	}
}

func (cn *newsController) GetNewsHandler(c *fiber.Ctx) error {
	var req GetNewsRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	news, err := cn.s.GetNews(c.UserContext(), req.CityID)
	if err != nil {
		if errors.Is(err, domain.NewsNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if news["status"] == "error" {
		return c.Status(fiber.StatusInternalServerError).JSON(news)
	}

	return c.JSON(news)
}

func (cn *newsController) GetNewsByTelegramHandler(c *fiber.Ctx) error {
	var req GetNewsByTelegramRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	user, err := cn.users.SelectByTelegramID(c.UserContext(), req.TelegramID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	news, err := cn.s.GetNews(c.UserContext(), user.CityN)
	if err != nil {
		if errors.Is(err, domain.NewsNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if news["status"] == "error" {
		return c.Status(fiber.StatusInternalServerError).JSON(news)
	}

	return c.JSON(news)
}
