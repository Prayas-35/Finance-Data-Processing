package handlers

import (
	"context"
	"time"

	"github.com/Prayas-35/Finance-Data-Processing/internal/services"
	"github.com/Prayas-35/Finance-Data-Processing/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	auth *services.AuthService
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req loginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.WriteError(c, fiber.StatusBadRequest, "invalid_input", "invalid request payload", nil)
	}

	token, err := h.auth.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		return utils.WriteError(c, fiber.StatusUnauthorized, "unauthorized", err.Error(), nil)
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   int(time.Hour.Seconds()),
		},
	})
}
