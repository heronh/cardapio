package models

type WeekDay struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Description string `json:"description"`
}
