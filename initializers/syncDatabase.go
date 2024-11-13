package initializers

func SyncDatabase() {
	// Sync the database
	DB.AutoMigrate(&models.User{})
}
