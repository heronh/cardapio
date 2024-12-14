package models

type Category struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	Name      string  `gorm:"type:varchar(100);not null"`
	CompanyID uint    `json:"companyid"`
	Company   Company `gorm:"foreignKey:CompanyID"`
	UserID    uint    `json:"userid"`
	User      User    `gorm:"foreignKey:UserID"`
}
