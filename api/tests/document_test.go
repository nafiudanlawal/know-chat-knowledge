package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"knowledge/api/config"
	"knowledge/api/database"
	"knowledge/api/migrations"
	"knowledge/api/models"
	"knowledge/api/routes"
	"knowledge/api/services"
)

var (
	app        *fiber.App
	db         *gorm.DB
)

// setupTestApp creates a test fiber application
func setupTestApp() (*fiber.App, *gorm.DB, error) {
	// Load test environment variables if exists
	_ = godotenv.Load("../.env.test")

	// Initialize config
	cfg := config.New()

	// Use in-memory SQLite for testing
	// Note: In a real application, you might use a test PostgreSQL database
	// or mock the database layer
	var err error
	db, err = database.Connect(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := migrations.Run(db); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create test app
	app := fiber.New()

	// Setup routes
	routes.Setup(app, db)

	return app, db, nil
}

// createTestDocument creates a test document in the database
func createTestDocument(t *testing.T) *models.Document {
	docService := services.NewDocumentService(db)

	doc := &models.Document{
		Name:    "Test Document",
		Type:    "text",
		Content: "This is a test document content",
	}

	err := docService.Create(doc)
	assert.NoError(t, err)

	return doc
}

// TestMain is used for setup and teardown of tests
func TestMain(m *testing.M) {
	// Setup test application
	var err error
	app, db, err = setupTestApp()
	if err != nil {
		fmt.Printf("Failed to setup test app: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	// In a real application, you would clean up the test database here

	os.Exit(code)
}

// TestCreateTextDocument tests creating a document with text content
func TestCreateTextDocument(t *testing.T) {
	// Create a new multipart writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add form fields
	_ = writer.WriteField("type", "text")
	_ = writer.WriteField("name", "Test Text Document")
	_ = writer.WriteField("content", "This is a test document content")

	// Close the writer
	writer.Close()

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer test-token")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse response
	var response models.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Status)
	assert.Equal(t, "document created", response.Message)
}

// TestGetDocuments tests retrieving a list of documents
func TestGetDocuments(t *testing.T) {
	// Create test document
	_ = createTestDocument(t)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var response models.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Status)
	assert.Equal(t, "success", response.Message)

	// Assert data is an array
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)
	assert.GreaterOrEqual(t, len(data), 1)
}

// TestGetDocument tests retrieving a single document
func TestGetDocument(t *testing.T) {
	// Create test document
	doc := createTestDocument(t)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/"+doc.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer test-token")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var response models.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Status)
	assert.Equal(t, "success", response.Message)

	// Assert data is a document
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, doc.Name, data["name"])
}

// TestUpdateDocument tests updating a document
func TestUpdateDocument(t *testing.T) {
	// Create test document
	doc := createTestDocument(t)

	// Create update payload
	updateData := map[string]string{
		"name": "Updated Document Name",
	}
	updateJSON, _ := json.Marshal(updateData)

	// Create request
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/documents/"+doc.ID.String(), bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var response models.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Status)

	// Assert document name was updated
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Updated Document Name", data["name"])
}

// TestReplaceDocument tests replacing a document
func TestReplaceDocument(t *testing.T) {
	// Create test document
	doc := createTestDocument(t)

	// Create a new multipart writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add form fields
	_ = writer.WriteField("type", "text")
	_ = writer.WriteField("name", "Replaced Document")
	_ = writer.WriteField("content", "This is the replaced content")

	// Close the writer
	writer.Close()

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/api/v1/documents/"+doc.ID.String(), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer test-token")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse response
	var response models.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Status)

	// Assert document name was updated
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Replaced Document", data["name"])
	assert.Equal(t, "This is the replaced content", data["content"])
}

// TestGetNonExistentDocument tests attempting to retrieve a non-existent document
func TestGetNonExistentDocument(t *testing.T) {
	// Create a random UUID that doesn't exist
	randomUUID := uuid.New().String()

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/"+randomUUID, nil)
	req.Header.Set("Authorization", "Bearer test-token")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["status"])
}

// TestAuthMiddleware tests that the authentication middleware rejects unauthorized requests
func TestAuthMiddleware(t *testing.T) {
	// Create request without auth header
	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents", nil)

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["status"])
	assert.Equal(t, "Access denied", response["message"])
}
