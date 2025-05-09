package routes

import (
        "github.com/gofiber/fiber/v3"
        "gorm.io/gorm"

        "knowledge/api/controllers"
        "knowledge/api/middlewares"
)

// Setup configures all the application routes
func Setup(app *fiber.App, db *gorm.DB) {
        // Create controllers
        documentController := controllers.NewDocumentController(db)

        // API v1 group
        api := app.Group("/api/v1")
        
        // Apply authentication middleware to all routes
        api.Use(middlewares.AuthMiddleware())

        // Document routes
        documents := api.Group("/documents")
        
        // GET /api/v1/documents
        documents.Get("/", documentController.GetDocuments)
        
        // POST /api/v1/documents
        documents.Post("/", documentController.CreateDocument)
        
        // GET /api/v1/documents/:id
        documents.Get("/:id", documentController.GetDocument)
        
        // PATCH /api/v1/documents/:id
        documents.Patch("/:id", documentController.UpdateDocument)
        
        // PUT /api/v1/documents/:id
        documents.Put("/:id", documentController.ReplaceDocument)
}
