package models

import (
	"time"
)

type Todo struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Completed   bool      `json:"completed"`
	Created_at  time.Time `json:"createdat"`
	Updated_at  time.Time `json:"updatedat"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user" gorm:"foreignKey:ID"`
	Description string    `json:"description"`
}
