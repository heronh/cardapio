package models

import (
	"time"
)

type Company struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:varchar(255)"`
	Category    string    `gorm:"type:varchar(100)"`
	Address     string    `gorm:"type:varchar(255)"`
	Street      string    `gorm:"type:varchar(100)"`
	Number      string    `gorm:"type:varchar(20)"`
	Complement  string    `gorm:"type:varchar(255)"`
	City        string    `gorm:"type:varchar(100)"`
	State       string    `gorm:"type:varchar(100)"`
	Country     string    `gorm:"type:varchar(100)"`
	PostalCode  string    `gorm:"type:varchar(20)"`
	Phone       string    `gorm:"type:varchar(20)"`
	Website     string    `gorm:"type:varchar(100)"`
}
