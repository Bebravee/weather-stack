package weathercntrl

import (
	"errors"

	"github.com/chup1x/weather-stack/internal/domain"
	weatherservice "github.com/chup1x/weather-stack/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type weatherController struct {
	s         *weatherservice.WeatherService
	validator *validator.Validate
}

func NewWeatherController(service *weatherservice.WeatherService) *weatherController {
	return &weatherController{
		validator: validator.New(),
		s:         service,
	}
}

func (cn *weatherController) CreateWeatherRecordHandler(c *fiber.Ctx) error {
	var req CreateWeatherRecord
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if err := cn.s.CreateWeatherRecord(c.UserContext(), &domain.WeatherEntity{
		Temperature: req.Temperature,
	}); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (cn *weatherController) GetWeatherHandler(c *fiber.Ctx) error {
	
	var req GetWeatherHistoryRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	city, err := cn.s.GetWeather(c.UserContext(), req.City_w)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
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
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	clothes, err := cn.s.GetWeatherClothes(c.UserContext(), req.user)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	res := GetWeatherClothesResponse{clothes}

	return c.JSON(res)
}

func (cn *weatherController) GetNewsHandler(c *fiber.Ctx) error {
	var req GetNewsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	city, err := cn.s.GetNews(c.UserContext(), req.city_n)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	res := GetNewsResponse{city}

	return c.JSON(res)
}
