package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CompanyID uint      `json:"company_id"`
	Company   Company   `gorm:"foreignKey:CompanyId"`
}
