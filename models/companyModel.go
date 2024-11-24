package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null"`
	Address     string `gorm:"type:varchar(255)"`
	City        string `gorm:"type:varchar(100)"`
	Complement  string `gorm:"type:varchar(255)"`
	Description string `gorm:"type:varchar(255)"`
	State       string `gorm:"type:varchar(100)"`
	Country     string `gorm:"type:varchar(100)"`
	PostalCode  string `gorm:"type:varchar(20)"`
	PhoneNumber string `gorm:"type:varchar(20)"`
	Website     string `gorm:"type:varchar(100)"`
	Logo        string `gorm:"type:varchar(255)"`
	Category    string `gorm:"type:varchar(20)"`
	UserID      uint
	User        User `gorm:"foreignKey:UserID"`
}
