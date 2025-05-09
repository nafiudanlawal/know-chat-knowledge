package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"knowledge/api/models"
	"knowledge/api/services"
	"knowledge/api/utils"
	"knowledge/api/validators"
)

// DocumentController handles document-related HTTP requests
type DocumentController struct {
	DB              *gorm.DB
	DocumentService *services.DocumentService
}

// NewDocumentController creates a new instance of DocumentController
func NewDocumentController(db *gorm.DB) *DocumentController {
	return &DocumentController{
		DB:              db,
		DocumentService: services.NewDocumentService(db),
	}
}

// CreateDocument handles the creation of a new document
func (c *DocumentController) CreateDocument(ctx fiber.Ctx) error {
	// Parse multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid form data: "+err.Error())
	}

	// Validate request
	docType := ctx.FormValue("type")
	name := ctx.FormValue("name")
	content := ctx.FormValue("content")
	file := ctx.FormValue("file")
	fmt.Println(docType, name, content, file)

	if err := validators.ValidateCreateDocument(docType, name, content, form.File["file"]); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Process document based on type
	var document *models.Document
	var fileURL string

	switch docType {

	case "file":
		// Process file upload
		files := form.File["documents"]
		// => []*multipart.FileHeader
		log.Println("documents count", len(files))
		fmt.Println("documents count", len(files))
		// Loop through files:
		for _, file := range files {
			log.Fatal(file.Filename, file.Size, file.Header["Content-Type"][0])
		}

		file, err := ctx.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "File upload error: "+err.Error())
		}

		fileURL, err = c.DocumentService.SaveFile(file)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to save file: "+err.Error())
		}

		document = &models.Document{
			Name:    name,
			Type:    "file",
			FileURL: fileURL,
		}

	case "text":
		// Process text content
		document = &models.Document{
			Name:    name,
			Type:    "text",
			Content: content,
		}

	default:
		return fiber.NewError(fiber.StatusBadRequest, "Invalid content type provided")
	}

	// Save document to database
	if err := c.DocumentService.Create(document); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create document: "+err.Error())
	}

	// Return response
	return ctx.Status(fiber.StatusCreated).JSON(models.Response{
		Status:  true,
		Message: "document created",
		Data:    document,
	})
}

// GetDocuments retrieves a paginated list of documents
func (c *DocumentController) GetDocuments(ctx fiber.Ctx) error {
	// Parse pagination parameters
	pagination := utils.NewPagination(ctx)

	// Get documents from service
	documents, err := c.DocumentService.GetAll(pagination)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve documents: "+err.Error())
	}

	// Return response
	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  true,
		Message: "success",
		Data:    documents,
	})
}

// GetDocument retrieves a single document by ID
func (c *DocumentController) GetDocument(ctx fiber.Ctx) error {
	// Parse document ID from URL
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Document ID is required")
	}

	// Parse ID to UUID
	docID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid document ID format")
	}

	// Get document from service
	document, err := c.DocumentService.GetByID(docID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Document not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve document: "+err.Error())
	}

	// Return response
	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  true,
		Message: "success",
		Data:    document,
	})
}

// UpdateDocument updates document metadata (PATCH)
func (c *DocumentController) UpdateDocument(ctx fiber.Ctx) error {
	// Parse document ID from URL
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Document ID is required")
	}

	// Parse ID to UUID
	docID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid document ID format")
	}

	// Parse request body
	var updateData struct {
		Name string `json:"name"`
	}

	if err := ctx.Bind().Body(&updateData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	// Validate name
	if updateData.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Name is required")
	}

	// Update document
	document, err := c.DocumentService.Update(docID, updateData.Name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Document not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update document: "+err.Error())
	}

	// Return response
	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  true,
		Message: "document updated",
		Data:    document,
	})
}

// ReplaceDocument replaces an entire document (PUT)
func (c *DocumentController) ReplaceDocument(ctx fiber.Ctx) error {
	// Parse document ID from URL
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Document ID is required")
	}

	// Parse ID to UUID
	docID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid document ID format")
	}

	// Verify document exists
	_, err = c.DocumentService.GetByID(docID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusNotFound, "Document not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve document: "+err.Error())
	}

	// Parse multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid form data: "+err.Error())
	}

	// Validate request
	docType := ctx.FormValue("type")
	name := ctx.FormValue("name")
	content := ctx.FormValue("content")

	if err := validators.ValidateCreateDocument(docType, name, content, form.File["file"]); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Process document based on type
	var fileURL string

	switch docType {
	case "file":
		// Process file upload
		file, err := ctx.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "File upload error: "+err.Error())
		}

		fileURL, err = c.DocumentService.SaveFile(file)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to save file: "+err.Error())
		}
	}

	// Replace document
	document, err := c.DocumentService.Replace(docID, models.Document{
		Name:    name,
		Type:    docType,
		Content: content,
		FileURL: fileURL,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to replace document: "+err.Error())
	}

	// Return response
	return ctx.Status(fiber.StatusOK).JSON(models.Response{
		Status:  true,
		Message: "document replaced",
		Data:    document,
	})
}
