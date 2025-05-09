package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"knowledge/api/config"
	"knowledge/api/database"
	"knowledge/api/migrations"
	"knowledge/api/routes"

	_ "github.com/joho/godotenv/autoload"
)

// Load Env FIle .env

func main() {
	// Initialize config
	config := config.New()

	// Initialize database
	db, err := database.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := migrations.Run(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// Default 500 statuscode
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error
				code = e.Code
			}

			// Set Content-Type: application/json
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			// Return statuscode with error message
			return c.Status(code).JSON(fiber.Map{
				"status":     false,
				"message":    err.Error(),
				"error_code": code,
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin, Content-Type, Accept, Authorization"},
	}))

	// Setup routes
	routes.Setup(app, db)

	// Determine port for HTTP service
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Start server
	log.Printf("Starting server on 0.0.0.0:%s\n", port)
	log.Fatal(app.Listen(":" + port, fiber.ListenConfig{
		//EnablePrefork: true,
		DisableStartupMessage: true,
	}))
}
