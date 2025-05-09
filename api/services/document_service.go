package services

import (
        "fmt"
        "mime/multipart"
        "os"
        "path/filepath"
        
        "github.com/google/uuid"
        "gorm.io/gorm"

        "knowledge/api/config"
        "knowledge/api/models"
        "knowledge/api/utils"
)

// DocumentService handles business logic for documents
type DocumentService struct {
        DB     *gorm.DB
        Config *config.Config
}

// NewDocumentService creates a new document service instance
func NewDocumentService(db *gorm.DB) *DocumentService {
        return &DocumentService{
                DB:     db,
                Config: config.New(),
        }
}

// Create creates a new document in the database
func (s *DocumentService) Create(document *models.Document) error {
        return s.DB.Create(document).Error
}

// GetAll retrieves all documents with pagination
func (s *DocumentService) GetAll(pagination *utils.Pagination) ([]models.Document, error) {
        var documents []models.Document
        
        query := s.DB.Model(&models.Document{})
        
        // Apply pagination
        query = query.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
        
        if err := query.Find(&documents).Error; err != nil {
                return nil, err
        }
        
        return documents, nil
}

// GetByID retrieves a document by its ID
func (s *DocumentService) GetByID(id uuid.UUID) (*models.Document, error) {
        var document models.Document
        
        if err := s.DB.Where("id = ?", id).First(&document).Error; err != nil {
                return nil, err
        }
        
        return &document, nil
}

// Update updates document metadata
func (s *DocumentService) Update(id uuid.UUID, name string) (*models.Document, error) {
        var document models.Document
        
        // Find document
        if err := s.DB.Where("id = ?", id).First(&document).Error; err != nil {
                return nil, err
        }
        
        // Update document name
        document.Name = name
        
        if err := s.DB.Save(&document).Error; err != nil {
                return nil, err
        }
        
        return &document, nil
}

// Replace replaces an entire document
func (s *DocumentService) Replace(id uuid.UUID, newDoc models.Document) (*models.Document, error) {
        var document models.Document
        
        // Find document
        if err := s.DB.Where("id = ?", id).First(&document).Error; err != nil {
                return nil, err
        }
        
        // Update document fields
        document.Name = newDoc.Name
        document.Type = newDoc.Type
        document.Content = newDoc.Content
        document.FileURL = newDoc.FileURL
        
        if err := s.DB.Save(&document).Error; err != nil {
                return nil, err
        }
        
        return &document, nil
}

// SaveFile saves an uploaded file to disk and returns the file URL
func (s *DocumentService) SaveFile(file *multipart.FileHeader) (string, error) {
        // Generate unique filename
        ext := filepath.Ext(file.Filename)
        filename := uuid.New().String() + ext
        
        // Create full path
        uploadDir := s.Config.FileStorage.Path
        filepath := filepath.Join(uploadDir, filename)
        
        // Save file
        if err := utils.SaveUploadedFile(file, filepath); err != nil {
                return "", err
        }
        
        // Generate file URL
        // In a production environment, this would likely be a URL to a CDN or file storage service
        fileURL := fmt.Sprintf("/uploads/%s", filename)
        
        return fileURL, nil
}

// Delete removes a document from the database
func (s *DocumentService) Delete(id uuid.UUID) error {
        var document models.Document
        
        // Find document
        if err := s.DB.Where("id = ?", id).First(&document).Error; err != nil {
                return err
        }
        
        // If document is a file, delete the physical file
        if document.Type == "file" && document.FileURL != "" {
                // Extract filename from FileURL
                filename := filepath.Base(document.FileURL)
                filePath := filepath.Join(s.Config.FileStorage.Path, filename)
                
                // Delete file
                if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
                        // Log error but continue with database deletion
                        fmt.Printf("Error deleting file %s: %v\n", filePath, err)
                }
        }
        
        // Delete from database
        return s.DB.Delete(&document).Error
}
