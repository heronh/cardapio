package initializers

import "github.com/heronh/cardapio/models"

func SyncDatabase() {
	// Sync the database
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Todo{})
	DB.AutoMigrate(&models.Company{})
	DB.AutoMigrate(&models.Dish{})
	DB.AutoMigrate(&models.Image{})
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Dish{}, &models.Image{})
}
