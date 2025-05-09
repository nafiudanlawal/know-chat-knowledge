package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveUploadedFile saves an uploaded file to the specified destination path
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open source file
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	// Copy file contents
	if _, err := io.Copy(out, src); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// GetFileExtension returns the extension of a file
func GetFileExtension(filename string) string {
	return filepath.Ext(filename)
}

// IsPdfFile checks if a file has a .pdf extension
func IsPdfFile(filename string) bool {
	ext := GetFileExtension(filename)
	return ext == ".pdf"
}
