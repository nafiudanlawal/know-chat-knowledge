package middlewares

import (
        "strings"

        "github.com/gofiber/fiber/v3"
)

// AuthMiddleware checks if a valid bearer token is provided in the request
func AuthMiddleware() fiber.Handler {
        return func(c fiber.Ctx) error {
                // Get the Authorization header
                authHeader := c.Get("Authorization")
                
                // Check if the header is empty or doesn't start with "Bearer "
                if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
                        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                                "status":     false,
                                "message":    "Access denied",
                                "error_code": fiber.StatusUnauthorized,
                        })
                }
                
                // Extract the token from the header
                token := strings.TrimPrefix(authHeader, "Bearer ")
                
                // In a real-world application, you would validate the token here
                // This is a simplified implementation that just checks for token presence
                if token == "" {
                        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                                "status":     false,
                                "message":    "Access denied",
                                "error_code": fiber.StatusUnauthorized,
                        })
                }
                
                // Token is valid, continue to the next middleware or handler
                return c.Next()
        }
}
