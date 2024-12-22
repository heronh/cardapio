package models

import (
	"time"
)

type Image struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
	UserID    uint      `json:"userid"`
	User      User      `gorm:"foreignKey:UserID"`
	CompanyID uint      `json:"companyid"`
	Company   Company   `gorm:"foreignKey:CompanyID"`
	Name      string    `gorm:"type:varchar(255)"`
	Path      string    `gorm:"type:varchar(255)"`
	Original  string    `gorm:"type:varchar(255)"`
	Dishes    []*Dish   `gorm:"many2many:dish_images;"`
}
