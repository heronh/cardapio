package models

import (
	"time"

	"gorm.io/gorm"
)

type Dish struct {
	gorm.Model
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:varchar(255)"`
	Ingredients string    `gorm:"type:varchar(255)"`
	Price       float64   `gorm:"type:decimal(10,2)"`
	CompanyID   uint      `json:"companyid"`
	Company     Company   `gorm:"foreignKey:CompanyID"`
	CategoryID  uint      `json:"categoryid"`
	Category    Category  `gorm:"foreignKey:CategoryID"`
	UserID      uint      `json:"userid"`
	User        User      `gorm:"foreignKey:UserID"`
	Enabled     bool      `gorm:"type:boolean"`
	Images      []*Image  `gorm:"many2many:dish_images;"`
	//Products  []Product `gorm:"foreignKey:ProductID;references:ID"`
}
