package usercntrl

import (
	"errors"

	"github.com/chup1x/weather-stack/internal/domain"
	"github.com/chup1x/weather-stack/internal/services"
	weatherservice "github.com/chup1x/weather-stack/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type userController struct {
	s         *services.UserService
	validator *validator.Validate
}

func NewUserController(service *weatherservice.UserService) *userController {
	return &userController{
		validator: validator.New(),
		s:         service,
	}
}

func (cn *userController) LoginHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (cn *userController) RegisterHandler(c *fiber.Ctx) error {
	var req RegisterProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	userID, err := cn.s.CreateUser(c.UserContext(), &domain.UserEntity{
		Name:       req.Name,
		Sex:        req.Sex,
		Age:        req.Age,
		City_n:     req.City_n,
		City_w:     req.City_w,
		Drop_time:	req.Drop_time,
		T_comf:		req.t_comf,
		T_tol:		req.t_tol,
		T_puh:		req.t_puh,
		Temp1:		req.temp1,
		TelegramID: req.telegram_id,
	})
	if err != nil {
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	res := RegisterProfileResponse{userID}

	return c.JSON(res)
}

func (cn *userController) GetTelegramProfileHandler(c *fiber.Ctx) error {
	var req GetTelegramProfileRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	user, err := cn.s.GetProfileByTelegramID(c.UserContext(), req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	res := GetTelegramProfileResponse{user}

	return c.JSON(res)
}

func (cn *userController) GetProfileHandler(c *fiber.Ctx) error {
	var req GetProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := cn.validator.Struct(req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	user, err := cn.s.GetProfileByID(c.UserContext(), req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error(err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	res := GetProfileResponse{user}

	return c.JSON(res)
}

func (cn *userController) UpdateProfileHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
