package migrations

import (
	"fmt"

	"gorm.io/gorm"
	
	"knowledge/api/models"
)

// Run executes all database migrations
func Run(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Auto migrate all models
	if err := db.AutoMigrate(
		&models.Document{},
		
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	fmt.Println("Database migrations completed successfully")
	return nil
}
