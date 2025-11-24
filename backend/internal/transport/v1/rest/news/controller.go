package weathercntrl

import (
	"errors"
	"encoding/json"
	"os"
	"github.com/chup1x/weather-stack/internal/domain"
	newsservice "github.com/chup1x/weather-stack/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type newsController struct {
	s         *newsservice.NewsService
	validator *validator.Validate
}

func NewNewsController(service *newsservice.NewsService) *newsController {
	return &newsController{
		validator: validator.New(),
		s:         service,
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
	data, err := os.ReadFile(news.Path)
	//err = json.Unmarshal(data, &res)
	//res := GetNewsResponse{news}

	return c.JSON(date)
}
