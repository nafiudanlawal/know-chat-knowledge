package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// Pagination represents pagination parameters
type Pagination struct {
	Offset int
	Limit  int
	Page   int
}

// NewPagination creates a new pagination object from request parameters
func NewPagination(c fiber.Ctx) *Pagination {
	// Default values
	offset := 0
	limit := 10
	page := 1

	// Parse offset
	if offsetParam := c.Query("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Parse limit
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Parse page
	if pageParam := c.Query("page"); pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
			// Override offset if page is provided
			offset = (page - 1) * limit
		}
	}

	return &Pagination{
		Offset: offset,
		Limit:  limit,
		Page:   page,
	}
}

// GetOffset returns the calculated offset
func (p *Pagination) GetOffset() int {
	return p.Offset
}

// GetLimit returns the limit
func (p *Pagination) GetLimit() int {
	return p.Limit
}

// GetPage returns the current page
func (p *Pagination) GetPage() int {
	return p.Page
}
