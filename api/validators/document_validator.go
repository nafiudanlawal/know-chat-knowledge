package validators

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// ValidateCreateDocument validates document creation parameters
func ValidateCreateDocument(docType, name, content string, files []*multipart.FileHeader) error {
	// Validate document type
	if docType == "" {
		return errors.New("document type is required")
	}

	if docType != "file" && docType != "text" {
		return errors.New("document type must be 'file' or 'text'")
	}

	// Validate document name
	if name == "" {
		return errors.New("document name is required")
	}

	// Validate based on document type
	switch docType {
	case "file":
		// Validate file
		if len(files) == 0 {
			return errors.New("file is required for file type documents")
		}

		file := files[0]

		// Validate file extension (only PDF allowed according to API spec)
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".pdf" {
			return errors.New("only PDF files are supported")
		}

		// Validate file size (max 10MB)
		if file.Size > 10*1024*1024 {
			return errors.New("file size exceeds the maximum allowed (10MB)")
		}

	case "text":
		// Validate content
		if content == "" {
			return errors.New("content is required for text type documents")
		}
	}

	return nil
}
