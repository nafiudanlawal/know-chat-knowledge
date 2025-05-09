package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Document represents a document in the system
type Document struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null"`
	Type      string         `json:"type" gorm:"not null"` // "file" or "text"
	Content   string         `json:"content,omitempty" gorm:"type:text"` // For text documents
	FileURL   string         `json:"file_url,omitempty"` // For file documents
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate is called before inserting a new document into the database
func (d *Document) BeforeCreate(tx *gorm.DB) error {
	// Generate a new UUID if not set
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
