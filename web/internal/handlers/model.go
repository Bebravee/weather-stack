package handlers

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserRequest struct {
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
	TelegramID int64  `json:"TelegramID" validate:"required"`
}

type RegisterUserResponse struct {
	ID uuid.UUID `json:"id"`
}

type GetUserRequest struct {
	ID uuid.UUID `param:"id"`
}

type GetUserResponse struct {
	ID         uuid.UUID `json:"id" gorm:"id"`
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
	TelegramID int64  `json:"TelegramID" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
}
