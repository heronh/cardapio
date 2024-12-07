package models

import (
	"time"

	"gorm.io/gorm"
)

type Weekday struct {
	Monday    bool `json:"monday"`
	Tuesday   bool `json:"tuesday"`
	Wednesday bool `json:"wednesday"`
	Thursday  bool `json:"thursday"`
	Friday    bool `json:"friday"`
	Saturday  bool `json:"saturday"`
	Sunday    bool `json:"sunday"`
}

type Dish struct {
	gorm.Model
	ID            uint      `json:"id" gorm:"primary_key"`
	CreatedAt     time.Time `json:"createdat"`
	UpdatedAt     time.Time `json:"updatedat"`
	Name          string    `gorm:"type:varchar(100);not null"`
	Description   string    `gorm:"type:varchar(255)"`
	Price         float64   `gorm:"type:decimal(10,2)"`
	CompanyID     uint      `json:"companyid"`
	Company       Company   `gorm:"foreignKey:CompanyID"`
	UserID        uint      `json:"userid"`
	User          User      `gorm:"foreignKey:UserID"`
	Enabled       bool      `gorm:"type:boolean"`
	DaysOfWeek    []string  `gorm:"type:text[]"`
	AvailableFrom time.Time `json:"availablefrom"`
	AvailableTo   time.Time `json:"availableto"`
}
