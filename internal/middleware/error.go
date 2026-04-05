package middleware

import (
	"log/slog"

	"github.com/Prayas-35/Finance-Data-Processing/internal/utils"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	slog.Error("request failed", "path", c.Path(), "method", c.Method(), "error", err.Error())
	return utils.WriteError(c, fiber.StatusInternalServerError, "internal_error", "internal server error", nil)
}
