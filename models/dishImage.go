package models

type DishImage struct {
	DishID  uint `gorm:"primaryKey"`
	ImageID uint `gorm:"primaryKey"`
}
