package middleware

import (
	"github.com/Prayas-35/Finance-Data-Processing/internal/utils"
	"github.com/gofiber/fiber/v3"
)

func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c fiber.Ctx) error {
		roleAny := c.Locals("role")
		role, ok := roleAny.(string)
		if !ok || role == "" {
			return utils.WriteError(c, fiber.StatusForbidden, "forbidden", "role not found", nil)
		}

		if _, exists := allowed[role]; !exists {
			return utils.WriteError(c, fiber.StatusForbidden, "forbidden", "insufficient permissions", nil)
		}

		return c.Next()
	}
}
