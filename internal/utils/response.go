package utils

import "github.com/gofiber/fiber/v3"

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func WriteError(c fiber.Ctx, status int, code, message string, details interface{}) error {
	return c.Status(status).JSON(ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}

func WriteOK(c fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{"data": data})
}
