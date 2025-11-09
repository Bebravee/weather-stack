package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type UserID struct {
	uuid.UUID
}

type UserEntity struct {
	ID         UserID `json:"id" gorm:"id"`
	Name       string `json:"name" gorm:"name"`
	Sex        string `json:"sex" gorm:"sex"`
	Age        int    `json:"age" gorm:"age"`
	City_n     string `json:"city_n" gorm:"city_n"`
	City_w     string `json:"city_w" gorm:"city_w"`
	Drop_time  string `json:"drop" gorm:"drop"`
	T_comf	   int	  `json:"t_comf" gorm:"t_comf"`
	T_tol	   int	  `json:"t_tol" gorm:"t_tol"`
	T_puh	   int	  `json:"t_puh" gorm:"t_puh"`
	Temp1	   int	  `json:"temp1" gorm:"temp1"`
	Passw      string `json:"password" gorm:"password"`
	TelegramID int64  `json:"telegram_id" gorm:"telegram_id"`
	CreatedAt  time.Time `json:"created_at" gorm:"created_at"`
}
