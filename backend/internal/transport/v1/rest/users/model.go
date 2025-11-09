package usercntrl

import (
	"github.com/chup1x/weather-stack/internal/domain"
)

type GetTelegramProfileRequest struct {
	ID int `param:"telegram_id" validate:"required"`
}

type GetTelegramProfileResponse struct {
	*domain.UserEntity
}

type GetProfileRequest struct {
	ID domain.UserID `json:"user_id" validate:"required"`
}

type GetProfileResponse struct {
	*domain.UserEntity
}

type RegisterProfileRequest struct {
	Name       string `json:"name" validate:"required"`
	Sex        string `json:"sex" validate:"required"`
	Age        int    `json:"age" validate:"required"`
	City_n     string `json:"city_n" validate:"required"`
	City_w     string `json:"city_w" validate:"required"`
	Drop_time  string `json:"drop" validate:"required"`
	t_comf     int    `json:"comf" validate:"required"`
	t_tol      int    `json:"tol" validate:"required"`
	t_puh      int    `json:"puh" validate:"required"`
	temp1      int    `json:"temp1" validate:"required"`
	TelegramID int64  `json:"telegram_id" validate:"required"`

type RegisterProfileResponse struct {
	ID domain.UserID `json:"id"`
}
