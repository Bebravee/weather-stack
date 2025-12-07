package handlers

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserRequest struct {
	Name       string `json:"name" validate:"required"`
	Sex        string `json:"sex" validate:"required"`
	Age        int    `json:"age" validate:"required"`
	CityN      string `json:"city_n" validate:"required"`
	CityW      string `json:"city_w" validate:"required"`
	DropTime   string `json:"drop_time" validate:"required"`
	TComfort   int    `json:"t_comfort" validate:"required"`
	TTol       int    `json:"t_tol" validate:"required"`
	TPuh       int    `json:"t_puh" validate:"required"`
	Temp1      int    `json:"temp1" validate:"required"`
	TelegramID int64  `json:"telegram_id" validate:"required"`
}

type RegisterUserResponse struct {
	ID uuid.UUID `json:"id"`
}

type GetUserRequest struct {
	ID uuid.UUID `params:"id"`
}

type GetUserResponse struct {
	ID         uuid.UUID `json:"id" gorm:"id"`
	Name       string    `json:"name" validate:"required"`
	Sex        string    `json:"sex" validate:"required"`
	Age        int       `json:"age" validate:"required"`
	CityN      string    `json:"city_n" validate:"required"`
	CityW      string    `json:"city_w" validate:"required"`
	DropTime   string    `json:"drop_time" validate:"required"`
	TComfort   int       `json:"t_comfort" validate:"required"`
	TTol       int       `json:"t_tol" validate:"required"`
	TPuh       int       `json:"t_puh" validate:"required"`
	Temp1      int       `json:"temp1" validate:"required"`
	TelegramID int64     `json:"telegram_id" validate:"required"`
	CreatedAt  time.Time `json:"created_at" gorm:"created_at"`
}

type GetCityRequest struct {
	City string `params:"city" validate:"required"`
}

type GetTelegramRequest struct {
	ID int `params:"telegram_id" validate:"required"`
}
