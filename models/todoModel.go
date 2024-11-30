package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"primary_key"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	Description string    `json:"description"`
}
