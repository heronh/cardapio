package models

type DishImage struct {
	ID      uint `gorm:"primaryKey"`
	DishID  uint `gorm:"not null"`
	ImageID uint `gorm:"not null"`
}
