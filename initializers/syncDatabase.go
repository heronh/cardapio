package initializers

import "github.com/heronh/cardapio/models"

func SyncDatabase() {
	// Sync the database
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Todo{})
}
