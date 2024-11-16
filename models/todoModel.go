package models

import (
	"time"
)

type Todo struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	Completed   bool      `json:"completed"`
	Created_at  time.Time `json:"createdat"`
	Updated_at  time.Time `json:"updatedat"`
	Description string    `json:"description"`
}
