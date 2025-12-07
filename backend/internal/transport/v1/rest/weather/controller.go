package weathercntrl

import (
	"context"
	"errors"

	"github.com/chup1x/weather-stack/internal/domain"
	weatherservice "github.com/chup1x/weather-stack/internal/services/weather"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type userProvider interface {
	SelectByTelegramID(context.Context, int) (*domain.UserEntity, error)
}

type weatherController struct {
	s         *weatherservice.WeatherService
	clothes   *weatherservice.ClothesService
	users     userProvider
	validator *validator.Validate
}

func NewWeatherController(service *weatherservice.WeatherService, clothes *weatherservice.ClothesService, users userProvider) *weatherController {
	return &weatherController{
		validator: validator.New(),
		s:         service,
		clothes:   clothes,
		users:     users,
	}
}

func (cn *weatherController) CreateWeatherRecordHandler(c *fiber.Ctx) error {
	var req CreateWeatherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if err := cn.s.CreateWeatherRecord(c.UserContext(), &domain.WeatherEntity{
		CityID:      req.CityID,
		Temperature: req.Temperature,
		FeelsLike:   req.FeelsLike,
		Description: req.Description,
		Humidity:    req.Humidity,
		Pressure:    req.Pressure,
		WindSpeed:   req.WindSpeed,
	}); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	res := CreateWeatherResponse{req.CityID}
	return c.JSON(res)
}

func (cn *weatherController) GetWeatherHandler(c *fiber.Ctx) error {
	var req GetWeatherHistoryRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	city, err := cn.s.GetWeatherWithCache(c.UserContext(), req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrWeatherNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	res := GetWeatherHistoryResponse{city}

	return c.JSON(res)
}

func (cn *weatherController) GetWeatherByTelegramHandler(c *fiber.Ctx) error {
	var req GetWeatherByTelegramRequest
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
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	city, err := cn.s.GetWeatherWithCache(c.UserContext(), user.CityW)
	if err != nil {
		if errors.Is(err, domain.ErrWeatherNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	res := GetWeatherHistoryResponse{city}
	return c.JSON(res)
}

func (cn *weatherController) GetWeatherClothesHandler(c *fiber.Ctx) error {
	var req GetWeatherClothesRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	clothes, err := cn.clothes.GetClothes(c.UserContext(), req.TelegramID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(clothes)
}
