package middleware

import (
	"strings"

	"github.com/Prayas-35/Finance-Data-Processing/internal/auth"
	"github.com/Prayas-35/Finance-Data-Processing/internal/utils"
	"github.com/gofiber/fiber/v3"
)

func RequireAuth(jwtManager *auth.JWTManager) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.WriteError(c, fiber.StatusUnauthorized, "unauthorized", "missing authorization header", nil)
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return utils.WriteError(c, fiber.StatusUnauthorized, "unauthorized", "invalid authorization format", nil)
		}

		claims, err := jwtManager.Parse(parts[1])
		if err != nil {
			return utils.WriteError(c, fiber.StatusUnauthorized, "unauthorized", "invalid token", nil)
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}
