package models

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	Created_at time.Time `json:"createdat"`
	Updated_at time.Time `json:"updatedat"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
}
