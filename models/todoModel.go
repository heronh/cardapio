package models

import (
	"time"
)

type Todo struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	Description string    `json:"description"`
}
