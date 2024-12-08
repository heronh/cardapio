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
	Name      string    `gorm:"type:varchar(255)"` // Original name of the image
	IsSample  bool      `gorm:"type:boolean"`      // If the image is a sample
}
