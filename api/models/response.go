package models

// Response represents a standard API response
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status    bool   `json:"status"`
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
}
